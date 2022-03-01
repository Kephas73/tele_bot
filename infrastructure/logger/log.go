package logger

import (
	"errors"
	"fmt"
	"io/ioutil"

	"botTele/constant"
	"github.com/getsentry/sentry-go"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var mLog *log.Logger

func NewLogger(logPath string, logPrefix string) *log.Logger {
	if mLog != nil {
		return mLog
	}

	logPathMap := lfshook.PathMap{
		log.InfoLevel:  logPath + "/" + logPrefix + "_success.log",
		log.TraceLevel: logPath + "/" + logPrefix + "_success.log",
		log.WarnLevel:  logPath + "/" + logPrefix + "_success.log",

		log.DebugLevel: logPath + "/" + logPrefix + "_debug.log",

		log.ErrorLevel: logPath + "/" + logPrefix + "_error.log",
		log.FatalLevel: logPath + "/" + logPrefix + "_error.log",
		log.PanicLevel: logPath + "/" + logPrefix + "_error.log",
	}

	logFormatter := new(log.TextFormatter)
	logFormatter.TimestampFormat = "02-01-2006 15:04:05"
	logFormatter.FullTimestamp = true

	mLog = log.New()
	mLog.Hooks.Add(lfshook.NewHook(
		logPathMap,
		logFormatter,
	))

	if !viper.GetBool(`Debug`) {
		mLog.Out = ioutil.Discard
	}

	return mLog
}

/**
 * Logs trace
 */
func Trace(format string, v ...interface{}) {
	mLog.Tracef(format, v)
}

/**
 * Logs info
 */
func Info(format string, v ...interface{}) {
	mLog.Infof(constant.LogInfoPrefix+format, v...)
}

/**
 * Logs warning
 */
func Warn(format string, v ...interface{}) {
	mLog.Warnf(format, v...)
}

/**
 * Logs debug
 */
func Debug(format string, v ...interface{}) {
	mLog.Debugf(format, v...)
}

/**
 * Logs error
 */
func Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	sentry.CaptureMessage(msg)
	err := errors.New(msg)
	sentry.CaptureException(err)
	mLog.Errorf(constant.LogErrorPrefix+format, v...)
}

/**
 * Logs fatal
 */
func Fatal(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	sentry.CaptureMessage(msg)
	err := errors.New(msg)
	sentry.CaptureException(err)
	mLog.Fatalf(format, v...)
}

func Panic(format string, v ...interface{}) {
	sentry.CaptureMessage(format)
	mLog.Panicf(format, v...)
}