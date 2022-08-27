package service

import (
    "botTele/constant"
    "botTele/infrastructure/bot"
    "botTele/infrastructure/logger"
    "botTele/model"
    "botTele/module/bot/repository"
    "bytes"
    "context"
    "encoding/csv"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/Kephas73/go-lib/lock_etcd"
    "github.com/Kephas73/go-lib/redis_client"
    "github.com/Kephas73/go-lib/s3_client"
    "github.com/Kephas73/go-lib/util"
    "github.com/go-redis/redis/v8"
    "github.com/jszwec/csvutil"
    "io"
    "io/ioutil"
    "mime/multipart"
    "strings"
    "sync"
    "time"
    "unicode/utf8"
)

type IBotService interface {
    SendChat(data model.RawData) error
    SendChatShutdown() error
    AutoReply()
    UploadFiles(ctx context.Context, files []*multipart.FileHeader) (*model.File, error)
    WorkerUploadFile()
    RandomIP() (string, error)
    InitIP() ([]string, error)
}

type BotService struct {
    Bot            *bot.TelegramBot
    FileRepository repository.IFileRepository
    Timeout        time.Duration
    LockEtcd       *lock_etcd.GEtcd
    sync.Mutex
    CacheRepository *redis_client.RedisPool
}

var botServiceInstance *BotService

func NewBotService(timeout time.Duration) IBotService {

    if botServiceInstance != nil {
        return botServiceInstance
    }

    //sqlxInstance := sql_client.GetSQLClient(constant.DB_FILE_UP_LOAD).DB

    return &BotService{
        Bot:     nil,
        Timeout: timeout,
        //FileRepository: repository.NewFileRepository(sqlxInstance),
        CacheRepository: redis_client.GetRedisClient(constant.TeleBotCache),
        LockEtcd:        lock_etcd.GetEtcdDiscoveryInstance(),
    }
}

func (bot *BotService) SendChat(data model.RawData) error {

    go func() {
        if len(data.Text) != constant.ValueEmpty {
            err := bot.Bot.SendMessageForApiTele(data.Text)
            if err != nil {
                logger.Error("BotService::SendChat: - Send chat error: %v", err)
            }
        }
    }()

    go func() {
        if data.Object != nil {
            b, _ := json.Marshal(data.Object)
            err := bot.Bot.SendMessageForApiTele(bytes.NewBuffer(b))
            if err != nil {
                logger.Error("BotService::SendChat: - Send chat error: %v", err)
            }
        }
    }()

    return nil
}

func (bot *BotService) AutoReply() {
    bot.Bot.AutoReply()
}

func (bot *BotService) SendChatShutdown() error {
    err := bot.Bot.SendChat(fmt.Sprintf(constant.ShutdownBot, bot.Bot.Bot.Self.FirstName+" "+bot.Bot.Bot.Self.LastName))
    if err != nil {
        logger.Error("BotService::SendChat: - Send chat error: %v", err)
        return err
    }
    return nil
}

func (bot *BotService) UploadFiles(ctx context.Context, files []*multipart.FileHeader) (*model.File, error) {

    ctx, cancel := context.WithTimeout(ctx, bot.Timeout)
    defer cancel()

    // Upload s3
    cdn, err := UploadS3(files)
    if err != nil {
        return nil, err
    }

    // Save db
    rs, err := bot.FileRepository.Insert(ctx, &model.File{
        FilePath:    cdn,
        State:       constant.ProcessingState,
        Description: "",
        CreatedTime: int32(time.Now().Unix()),
        UpdatedTime: int32(time.Now().Unix()),
    })

    if err != nil {
        logger.Error("BotService::UploadFiles: - Insert file error: %v", err)
        return nil, err
    }

    return rs, nil
}

func (bot *BotService) WorkerUploadFile() {

    ctx, cancel := context.WithTimeout(context.Background(), bot.Timeout)
    defer cancel()

    // Save db
    rs, err := bot.FileRepository.SelectFileByState(ctx, constant.ProcessingState)
    if err != nil {
        logger.Error("WorkerUploadFile: - Get file processing error: %v", err)
        return
    }

    if len(rs) != 0 {
        wg := sync.WaitGroup{}
        wg.Add(len(rs))
        for idx, v := range rs {
            go func(i int, file *model.File) {
                defer wg.Done()
                // Get file S3
                members := make([]*model.Member, 0)
                content, err := s3_client.GetS3ClientInstance().GetFileS3(file.FilePath)
                if err != nil {
                    logger.Error("WorkerUploadFile: - Get file path: %s S3 error: %v", file.FilePath, err)
                    rs[i].State = constant.FailedState
                    rs[i].Description = err.Error()
                    rs[i].UpdatedTime = int32(time.Now().Unix())
                    goto SAVE_STATE
                }

                // Decode content
                members, err = DecodeMember(strings.NewReader(content))
                if err != nil {
                    logger.Error("WorkerUploadFile: - Decode file path: %s error: %v", file.FilePath, err)
                    rs[i].State = constant.FailedState
                    rs[i].Description = err.Error()
                    rs[i].UpdatedTime = int32(time.Now().Unix())
                    goto SAVE_STATE
                }

                // Validate
                err = ValidateMember(members)
                if err != nil {
                    logger.Error("WorkerUploadFile: - Validate file path: %s error: %v", file.FilePath, err)
                    rs[i].State = constant.FailedState
                    rs[i].Description = err.Error()
                    rs[i].UpdatedTime = int32(time.Now().Unix())
                    goto SAVE_STATE
                }

                // Insert DB
                // Trường hợp insert lỗi thì cũng gán lỗi vòa trạng thái file đó.
                // Chổ này có thể sử dụng transaction để insert data và update trạng thái
                // Vì có trường hợp update trạng thái lỗi nhưng đã insert rồi, thì lần sau job chạy có thể bị dub data
                for _, m := range members {
                    fmt.Println(fmt.Sprintf("Name: %s - Address: %s - Age: %d", m.Name, m.Address, m.Age))
                }

                rs[i].State = constant.SuccessState
                rs[i].UpdatedTime = int32(time.Now().Unix())

                // Save state DB
            SAVE_STATE:
                _, err = bot.FileRepository.Update(context.Background(), rs[i])
                if err != nil {
                    logger.Error("WorkerUploadFile: - Update state file path error: %v", err)
                    return
                }

                // Remove file S3
                // Xóa cho sạch, trường hợp không xóa lỗi thì vẫn ko có vấn đề gì
                go s3_client.GetS3ClientInstance().RemoveFileS3(file.FilePath)

            }(idx, v)
        }
        wg.Wait()
    }

    return
}

func UploadS3(files []*multipart.FileHeader) (string, error) {
    if len(files) == 0 {
        logger.Error("file is empty!")
        return "", fmt.Errorf("file is empty")
    }

    fileObject := files[0]
    fileMimeType, err := util.GetExtFile(fileObject.Header.Get("Content-Type"))
    if err != nil {
        logger.Error("get ext file, err: %v", err)
        return "", err
    }

    file, err := fileObject.Open()
    if err != nil {
        logger.Error("can't open file, err: %v", err)
        return "", err
    }
    defer file.Close()

    imageBytes, err := ioutil.ReadAll(file)
    if err != nil {
        logger.Error("can't read file, err: %v", err)
        return "", err
    }

    fileObject.Filename = fmt.Sprintf("%v_%d.%v", strings.ReplaceAll(strings.Split(fileObject.Filename, ".")[0], " ", ""), time.Now().Unix(), fileMimeType)
    s3FilePath := fmt.Sprintf("%v/%v", "go-lib", fileObject.Filename)

    //4. Upload to S3
    cdn, err := s3_client.GetS3ClientInstance().UploadFile(imageBytes, s3FilePath)
    if err != nil {
        logger.Error("can't upload file s3, err: %v", err)
        return "", err
    }

    logger.Info("Cdn S3: %v", cdn)

    return s3FilePath, nil
}

func DecodeMember(r io.Reader) (members []*model.Member, err error) {
    csvReader := csv.NewReader(r)
    csvReader.TrimLeadingSpace = constant.ModeTrimLeadingSpace
    csvReader.LazyQuotes = constant.ModeLazyQuotes
    csvReader.FieldsPerRecord = constant.ModeFieldsPerRecord
    comma, _ := utf8.DecodeRuneInString(constant.ModeComma)
    csvReader.Comma = comma

    // in real application this should be done once in init function.
    userHeader, err := csvutil.Header(model.Member{}, "csv")
    if err != nil {
        return nil, err
    }
    dec, err := csvutil.NewDecoder(csvReader, userHeader...)
    if err != nil {
        return nil, err
    }

    for {
        var u model.Member
        if err := dec.Decode(&u); err == io.EOF {
            break
        } else if err != nil {
            return nil, err
        }

        members = append(members, &u)
    }
    return members, nil
}

func ValidateMember([]*model.Member) error {
    return nil
}

func (bot *BotService) RandomIP() (string, error) {

    //// Request vào chờ để lấy khóa
    //mux := bot.LockEtcd.Locker(constant.LockRandomIP)
    ///*    mux.Lock()
    //     mux.Unlock()*/
    //if err := mux.Lock().LockerActionError; err != nil {
    //   logger.Error("BotService::RandomIP: Lock Etcd error: %v", err)
    //   return "", err
    //}
    //defer mux.Unlock()

    bot.Lock()
    defer bot.Unlock()

    timeNow := time.Now()
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
    defer cancel()

    // Lấy phần tử có trọng số thấp nhấp 
    // 0: lấy bắt đầu từ phần tử thứ 0
    // -1: lấy tất cả.
    listIP, err := bot.CacheRepository.Get().ZRevRange(ctx, constant.KeyListIP, 0, -1).Result()
    if err != nil {
        // Chạy lại quá trình init
        if listIP, err = bot.initIP(); err != nil {
            logger.Error("BotService::RandomIP: Get list or init ip error: %v", err)
            return "", err
        }
    }

    if len(listIP) == 0 {
        logger.Error("BotService::RandomIP: List ip empty!")
        return "", errors.New("empty")
    }

    // Sau khi lấy thì thực hiện tăng trong số lên 1 cho phần tử
    err = bot.CacheRepository.Get().ZIncrBy(context.Background(), constant.KeyListIP, 1, listIP[len(listIP)-1]).Err()
    if err != nil {
        logger.Error("BotService::RandomIP: ZIncrBy error: %v", err)
        return "", err
    }
    logger.Info("BotService::RandomIP: Done: %d", time.Since(timeNow).Milliseconds())

    // Lấy thằng vị trí cuối cùng, vì list đã được sắp xếp rồi.
    return listIP[len(listIP)-1], nil
}

// InitIP func:
// Chỉ chạy khi lấy redis bị lỗi
// Hoặc khi thêm mới các ip mới
// Hoặc admin muốn reset lại từ đầu
// Hoặc khi khởi động lại lại service
func (bot *BotService) InitIP() ([]string, error) {

    //// Request vào chờ để lấy khóa
    //mux := bot.LockEtcd.Locker(constant.LockRandomIP)
    //if err := mux.Lock().LockerActionError; err != nil {
    //    logger.Error("BotService::RandomIP: Lock Etcd error: %v", err)
    //}
    //defer mux.Unlock()
    bot.Lock()
    defer bot.Unlock()

    return bot.initIP()
}

func (bot *BotService) initIP() ([]string, error) {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
    defer cancel()

    // Đoạn này có thể lưu IP ở DB để lấy ra và lưu
    pipe := bot.CacheRepository.Get().Pipeline()
    for _, ip := range constant.ListIP {
        pipe.ZAdd(ctx, constant.KeyListIP, &redis.Z{
            Member: ip,
            Score:  0, // Default init
        })
    }

    _, err := pipe.Exec(ctx)
    if err != nil {
        logger.Error("BotService::initIP: Error: %v", err)
        return nil, err
    }

    return constant.ListIP, nil
}
