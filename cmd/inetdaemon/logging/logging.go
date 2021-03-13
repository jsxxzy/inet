// Author: d1y<chenhonzhou@gmail.com>
// 日志

package logging

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jsxxzy/inet/cmd/inetdaemon/xfs"
)

// Logging 日志
type Logging struct {
	Dir string // 相对路径目录, 用于创建日志的文件夹
}

// CleanType 清除类型
type CleanType int

const (
	ClearDay CleanType = 0 // 清除当天
	ClearAll CleanType = 1 // 清除所有
)

// LogType 打印日志类型
type LogType int

const (
	ErrorMsg LogType = 2 // 错误日志
	InfoMsg  LogType = 3 // 正常日志
)

const (
	ErrorTag string = "Error"
	InfoTag  string = "Info"
)

func getDateFormat() string {
	t := time.Now()
	formatted := fmt.Sprintf("%d-%02d-%02d", t.Year(), t.Month(), t.Day())
	return formatted + ".log"
}

// New 创建一个日志, 传入一个日志文件夹
func New(dirName string) *Logging {
	var pulicPrem os.FileMode = 0777
	xfs.CreateIfNotExists(dirName, pulicPrem)
	return &Logging{
		Dir: dirName,
	}
}

// GetLog 获取日志
func (l *Logging) GetLog() string {
	var file = l.initDayFile()
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "日志访问错误"
	}
	var tmp = string(data)
	if len(tmp) <= 0 {
		return "日志为空"
	}
	return tmp
}

// 初始化日志文件
func (l *Logging) initDayFile() string {
	var filename = getDateFormat()
	filename = filepath.Join(l.Dir, filename)
	// fmt.Println("filename: ", filename)
	if !xfs.Exists(filename) {
		var file, _ = os.Create(filename)
		defer file.Close()
	}
	return filename
}

func easyGetLogTag(action LogType) (r string) {
	switch action {
	case ErrorMsg:
		r = ErrorTag
	case InfoMsg:
		r = InfoTag
	default:
		r = "unknown"
	}
	r = "[" + r + "]"
	return r
}

func (l *Logging) Println(msg string, action LogType) {
	var targetFile = l.initDayFile()
	f, err := os.OpenFile(targetFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	var t = time.Now()
	var msga = t.Format("2006-01-02 15:04:05") + " " + easyGetLogTag(action) + " " + msg + "\n"
	fmt.Println(msga)
	if _, err := f.WriteString(msga); err != nil {
		log.Println(err)
	}
	defer f.Close()
}

func (l *Logging) Info(msg string) {
	l.Println(msg, InfoMsg)
}

func (l *Logging) Error(msg string) {
	l.Println(msg, ErrorMsg)
}

// Clean 清理日志文件
func (l *Logging) Clean(actionType CleanType) bool {
	switch actionType {
	case ClearDay:
		// TODO
	case ClearAll:
		// TODO
	}
	return false
}
