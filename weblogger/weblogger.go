package weblogger

import (
	"log"
	"os"
	"runtime"
	"time"
)

// Weblogger is logger with attributes
type Weblogger struct {
	Logger    *log.Logger
	ClientIP  string
	RPCmethod string
}

// severity tag
const (
	tagError = " ERROR"
	tagFatal = " FATAL"
)

func isLogFileExist(filePath string) (*os.File, error) {
	logfile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	return logfile, err
}

// InitWebLogger is to init a web logger
func (w *Weblogger) InitWebLogger(filePath string) {
	logfile, err := isLogFileExist(filePath)
	if err != nil {
		w.ServerFatalPrintln("File open error", err)
	}
	w.Logger = log.New(logfile, time.Now().Format("2006-01-02T15:04:05.99-07:00")+" ", 0)
}

// AccessPrintln print to the accessLog with access message
func (w *Weblogger) AccessPrintln(rpcMethod string, para string) {
	w.RPCmethod = rpcMethod
	w.Logger.Println(w.ClientIP, w.RPCmethod, para)
}

// ErrorPrintln print to the errorLog with ERROR message
func (w *Weblogger) ErrorPrintln(rpcMethod string, s string) {
	w.RPCmethod = rpcMethod
	_, fileName, line, _ := runtime.Caller(1)
	w.Logger.Println(tagError, w.ClientIP, w.RPCmethod, fileName, line, s)
}

// FatalPrintln print to the errorLog with FATAL message
func (w *Weblogger) FatalPrintln(rpcMethod string, s string, err error) {
	w.RPCmethod = rpcMethod
	_, fileName, line, _ := runtime.Caller(1)
	w.Logger.Println(tagFatal, w.ClientIP, w.RPCmethod, fileName, line, s, err)
}

// ServerFatalPrintln print to the errorLog with FATAL message
func (w *Weblogger) ServerFatalPrintln(s string, err error) {
	_, fileName, line, _ := runtime.Caller(1)
	w.Logger.Println(tagFatal, fileName, line, s, err)
}
