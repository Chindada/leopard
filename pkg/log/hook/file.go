// Package hook package hook
package hook

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	_defaultTimeFormat            = time.RFC3339
	_defaultLogFileNameTimeFormat = time.DateOnly
	_defaultFileMaxAge            = 7 * 24 * time.Hour
)

var (
	singlton *File
	once     sync.Once
)

type File struct {
	logFolder string
	appName   string

	reportCaller     bool
	callerPrettyfier func(*runtime.Frame) (function string, file string)

	levels   []logrus.Level
	f        *os.File
	lock     sync.Mutex
	lastTime time.Time
}

func Get(level logrus.Level, path, appName string) *File {
	if singlton == nil {
		once.Do(func() {
			hook := &File{
				logFolder: path,
				appName:   appName,
			}
			for _, l := range logrus.AllLevels {
				if l <= level {
					hook.levels = append(hook.levels, l)
				}
			}
			singlton = hook
		})
		return Get(level, path, appName)
	}
	return singlton
}

func (h *File) SetReportCaller(reportCaller bool, f func(*runtime.Frame) (function string, file string)) {
	h.reportCaller = reportCaller
	h.callerPrettyfier = f
}

func (h *File) Levels() []logrus.Level {
	return h.levels
}

func (h *File) Fire(entry *logrus.Entry) error {
	defer h.lock.Unlock()
	h.lock.Lock()

	msg, err := h.Format(entry)
	if err != nil {
		return err
	}
	if h.lastTime.Day() != entry.Time.Day() || h.lastTime.IsZero() {
		h.getFile()
	}
	_, err = h.f.Write(msg)
	if err != nil {
		return err
	}
	h.lastTime = entry.Time
	return nil
}

func (h *File) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	levelText := strings.ToUpper(entry.Level.String())[0:4]
	header := fmt.Sprintf(
		"%s[%s][%s]",
		levelText,
		h.appName,
		entry.Time.Format(time.DateTime),
	)
	if entry.Caller != nil && h.callerPrettyfier != nil {
		caller, _ := h.callerPrettyfier(entry.Caller)
		header = fmt.Sprintf(header+"[%s]", caller)
	}
	msg := fmt.Sprintf("%s %-44s\n", header, entry.Message)
	_, e := b.WriteString(msg)
	if e != nil {
		return nil, e
	}
	return b.Bytes(), nil
}

func (h *File) getFile() {
	if h.f != nil {
		_ = h.f.Close()
	}
	if _, err := os.Stat(h.logFolder); os.IsNotExist(err) {
		if err = os.MkdirAll(h.logFolder, os.ModePerm); err != nil {
			panic(err)
		}
	}
	path := fmt.Sprintf("%s/%s-%s.log", h.logFolder, h.appName, time.Now().Format(_defaultLogFileNameTimeFormat))
	if _, err := os.Stat(path); err == nil {
		f, openErr := os.OpenFile(filepath.Clean(path), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if openErr != nil {
			panic(openErr)
		}
		h.f = f
	} else {
		f, createErr := os.Create(filepath.Clean(path))
		if createErr != nil {
			panic(createErr)
		}
		h.f = f
	}
	h.deleteExpireFile()
}

func (h *File) deleteExpireFile() {
	files, err := os.ReadDir(h.logFolder)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		re := regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`)
		match := re.FindString(file.Name())
		if match == "" {
			continue
		}
		t, pErr := time.ParseInLocation(_defaultLogFileNameTimeFormat, match, time.Local)
		if pErr != nil {
			continue
		}
		if time.Since(t) > _defaultFileMaxAge {
			if removeErr := os.Remove(filepath.Join(h.logFolder, file.Name())); removeErr != nil {
				panic(removeErr)
			}
		}
	}
}
