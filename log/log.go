// Copyright 2014 li. All rights reserved.

package log

import (
	"github.com/agtorre/gocolorize"
	"github.com/roverli/light/conf"
	"io"
	"log"
	"os"
	"runtime"
)

type colorLogs struct {
	c gocolorize.Colorize
	w io.Writer
}

func (r *colorLogs) Write(p []byte) (n int, err error) {
	return r.w.Write([]byte(r.c.Paint(string(p))))
}

var (
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
)

var logConf, _ = conf.Load("log.conf")

func init() {
	if logConf == nil {
		return
	}
	if runtime.GOOS == "windows" {
		gocolorize.SetPlain(true)
	}

	debugLogger = newLogger("debug")
	infoLogger = newLogger("info")
	warnLogger = newLogger("warn")
	errorLogger = newLogger("error")
}

func newLogger(name string) *log.Logger {

	color := gocolorize.NewColor(logConf.String(name+".color", "black"))

	var w io.Writer
	switch out := logConf[name+".output"]; out {
	case "", "off":
		return nil

	case "stdout":
		w = &colorLogs{c: color, w: os.Stdout}

	case "stderr":
		w = &colorLogs{c: color, w: os.Stderr}

	default:
		f, err := os.OpenFile(out, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("light/log: Open file %s error, %v.", out, err)
			return nil
		}
		w = f
	}
	return log.New(w, logConf[name+".prefix"],
		logConf.Int(name+".flags", log.Ldate|log.Ltime|log.Lshortfile))
}

func Debugf(format string, v ...interface{}) {
	if conf.IsDebug {
		printf(debugLogger, format, v...)
	}
}

func Debug(v ...interface{}) {
	if conf.IsDebug {
		println(debugLogger, v...)
	}
}

func Infof(format string, v ...interface{}) {
	printf(infoLogger, format, v...)
}

func Info(v ...interface{}) {
	println(infoLogger, v...)
}

func Warnf(format string, v ...interface{}) {
	printf(warnLogger, format, v...)
}

func Warn(v ...interface{}) {
	println(warnLogger, v...)
}

func Errorf(format string, v ...interface{}) {
	printf(errorLogger, format, v...)
}

func Error(v ...interface{}) {
	println(errorLogger, v...)
}

func println(logger *log.Logger, v ...interface{}) {
	if logger != nil {
		logger.Println(v...)
	}
}
func printf(logger *log.Logger, format string, v ...interface{}) {
	if logger != nil {
		logger.Printf(format, v...)
	}
}
