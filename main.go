package main

import (
    "botTele/global"
    "botTele/infrastructure/logger"
    "botTele/module/bot"
    "fmt"
    logger2 "github.com/Kephas73/go-lib/logger"
    "github.com/labstack/echo"
    "github.com/labstack/echo/middleware"
    "github.com/spf13/viper"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func init() {
    viper.SetConfigFile(`config.json`)

    err := viper.ReadInConfig()

    if err != nil {
        panic(err)
    }

    if viper.GetBool(`Debug`) {
        fmt.Println("Service RUN on DEBUG mode")
    } else {
        fmt.Println("Service RUN on PRODUCTION mode")
    }
}

func main() {
    logPath := viper.GetString("Log.Path")
    logPrefix := viper.GetString("Log.Prefix")
    logger2.NewLogger(logPath, logPrefix)
    logger.NewLogger(logPath, logPrefix)

    //s3_client.InstallS3Client()
    //sql_client.InstallSQLClientManager()
    //lock_etcd.InstanceEtcdManger()
    //redis_client.InstallRedisClientManager()

    timeout := time.Duration(viper.GetInt("Context.Timeout")) * time.Second

    e := echo.New()
    e.Server.SetKeepAlivesEnabled(false)
    e.Server.ReadTimeout = time.Minute * 60
    e.Server.WriteTimeout = time.Minute * 60

    e.Use(middleware.CORS())

    signChan := make(chan os.Signal, 1)
    
    // Init global biến class luôn
    _, err := global.InitClassGlobal()
    if err != nil {
        panic(e)
    }
    
    //healthcheck.Initialize(e, timeout)
    bot.Initialize(e, timeout)

    go e.Start(viper.GetString("Server.Address"))
    signal.Notify(signChan, os.Interrupt, syscall.SIGTERM)
    <-signChan
    logger.Info("Shutdown.....")
    //bot.BotServiceGlobal.SendChatShutdown()
}
