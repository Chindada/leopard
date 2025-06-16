// Package log package log
package log

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/chindada/leopard/pkg/log/hook"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

type Log struct {
	*logrus.Logger
	*config
}

type config struct {
	Level          string `env:"LOG_LEVEL"`
	NeedCaller     bool   `env:"LOG_NEED_CALLER"`
	DisableConsole bool   `env:"LOG_DISABLE_CONSOLE"`
	DisableFile    bool   `env:"LOG_DISABLE_FILE"`
}

var (
	singleton *Log
	once      sync.Once
)

func Init() {
	var appName string
	_, caller, _, ok := runtime.Caller(1)
	if ok {
		appName = strings.ToUpper(filepath.Base(filepath.Dir(caller)))
		appName = appName[:3]
	} else {
		appName = "UNKNOWN"
	}
	once.Do(func() {
		l := &Log{
			Logger: logrus.New(),
		}
		l.readConfig(appName)
		l.setFileHook(appName)
		singleton = l
	})
}

func Get() *Log {
	if singleton == nil {
		panic("log not initialized")
	}
	return singleton
}

func (l *Log) readConfig(appName string) {
	cfg := config{}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic(err)
	}
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	l.Logger.SetFormatter(&consoleFormatter{
		appName:          appName,
		callerPrettyfier: l.callerPrettyfier(),
	})
	l.Logger.SetLevel(level)
	l.Logger.SetReportCaller(cfg.NeedCaller)
	if cfg.DisableConsole {
		l.Logger.SetOutput(io.Discard)
	}
	l.config = &cfg
}

func (l *Log) setFileHook(appName string) {
	if l.config.DisableFile {
		return
	}
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	fileHook := hook.Get(
		l.Logger.Level,
		filepath.Join(filepath.Clean(filepath.Dir(ex)), "..", "logs"),
		appName,
	)
	fileHook.SetReportCaller(l.config.NeedCaller, l.callerPrettyfier())
	l.Logger.AddHook(fileHook)
}

func (l *Log) callerPrettyfier() func(*runtime.Frame) (string, string) {
	return func(frame *runtime.Frame) (string, string) {
		split := strings.Split(frame.File, "/")
		if len(split) < 2 {
			return "", ""
		}
		return fmt.Sprintf("%s:%d", split[len(split)-1], frame.Line), ""
	}
}

type consoleFormatter struct {
	appName          string
	callerPrettyfier func(*runtime.Frame) (string, string)
}

const (
	red    = 31
	yellow = 33
	blue   = 36
	gray   = 37
)

func (f *consoleFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	levelText := strings.ToUpper(entry.Level.String())[0:4]
	var levelColor int
	switch entry.Level {
	case logrus.DebugLevel, logrus.TraceLevel:
		levelColor = gray
	case logrus.WarnLevel:
		levelColor = yellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		levelColor = red
	case logrus.InfoLevel:
		levelColor = blue
	default:
		levelColor = blue
	}
	header := fmt.Sprintf(
		"\x1b[%dm%s\x1b[33m[%s]\x1b[0m[%s]",
		levelColor,
		levelText,
		f.appName,
		entry.Time.Format(time.DateTime),
	)
	if entry.Caller != nil && f.callerPrettyfier != nil {
		caller, _ := f.callerPrettyfier(entry.Caller)
		header = fmt.Sprintf(header+"[%s]", caller)
	}

	msg := fmt.Sprintf("%s %-44s\n", header, entry.Message)
	_, e := b.WriteString(msg)
	if e != nil {
		return nil, e
	}
	return b.Bytes(), nil
}
