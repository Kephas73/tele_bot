package global

import (
    "botTele/constant"
    "botTele/model"
    "encoding/csv"
    "errors"
    "github.com/jszwec/csvutil"
    "github.com/spf13/viper"
    "io"
    "os"
    "unicode/utf8"
)

var class []*model.Class

func InitClassGlobal() ([]*model.Class, error) {
    // path file csv
    path := viper.GetString("PathClassFile")
    file, err := os.Open(path)
    defer file.Close()
    if err != nil {
        return nil, err
    }

    rs, err := ReadFileClass(file)
    if err != nil {
        return nil, err
    }

    if len(rs) == 0 {
        return nil, errors.New("class is empty")
    }

    class = rs
    return class, nil
}

func GetClassGlobal() ([]*model.Class, error) {
    // Kiểm tra nó xem còn tồn tại biến này hay ko
    // Ko có thì Init lại
    if len(class) == 0 {
        rs, err := InitClassGlobal()
        if err != nil {
            return nil, err
        }

        if len(rs) == 0 {
            return nil, errors.New("class is empty")
        }

        class = rs
    }
  
    return class, nil
}

func ReadFileClass(r io.Reader) (class []*model.Class, err error) {
    csvReader := csv.NewReader(r)
    csvReader.TrimLeadingSpace = constant.ModeTrimLeadingSpace
    csvReader.LazyQuotes = constant.ModeLazyQuotes
    csvReader.FieldsPerRecord = constant.ModeFieldsPerRecord
    comma, _ := utf8.DecodeRuneInString(constant.ModeComma2)
    csvReader.Comma = comma

    // in real application this should be done once in init function.
    // Ko cần header thì dùng đoạn này
    /*userHeader, err := csvutil.Header(model.Class{}, "csv")
    if err != nil {
       return nil, err
    }

    dec, err := csvutil.NewDecoder(csvReader)
    if err != nil {
        return nil, err
    }*/
    
    // Có header thì dùng đoạn này.
    dec, err := csvutil.NewDecoder(csvReader)
    if err != nil {
        return nil, err
    }

    for {
        var u model.Class
        if err := dec.Decode(&u); err == io.EOF {
            break
        } else if err != nil {
            return nil, err
        }

        class = append(class, &u)
    }
    return class, nil
}
