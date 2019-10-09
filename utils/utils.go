package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hhkbp2/go-logging"
	"os"
	"reflect"
	"strings"
	"time"
)

var startTime time.Time

func SplitStringParameter(parameterValue string, separator string) []string {
	return strings.Split(parameterValue, separator)
}

func InitLogger(name string) logging.Logger {
	logger := logging.GetLogger(name)
	handler := logging.NewStdoutHandler()

	format := "%(asctime)s %(levelname)s (%(filename)s:%(lineno)d) " +
		"%(name)s %(message)s"
	// the format for the time part
	dateFormat := "%Y-%m-%d %H:%M:%S.%3n"
	// create a formatter(which controls how log messages are formatted)
	formatter := logging.NewStandardFormatter(format, dateFormat)

	handler.SetFormatter(formatter)
	logger.SetLevel(logging.LevelInfo)
	logger.AddHandler(handler)
	return logger
}

func CheckError(logger logging.Logger, message string, err error) {
	if err != nil {
		logger.Fatalf(message, err)
	}
}

func contains(a interface{}, e interface{}) bool {
	v := reflect.ValueOf(a)

	for i := 0; i < v.Len(); i++ {
		if v.Index(i).Interface() == e {
			return true
		}
	}
	return false
}

func Intersect(a interface{}, b interface{}) []interface{} {
	set := make([]interface{}, 0)
	av := reflect.ValueOf(a)

	for i := 0; i < av.Len(); i++ {
		el := av.Index(i).Interface()
		if contains(b, el) {
			set = append(set, el)
		}
	}

	return set
}

func InterfaceToString(t []interface{}) []string {
	s := make([]string, len(t))
	for i, v := range t {
		s[i] = fmt.Sprint(v)
	}
	return s
}

func FileLastModification(filename string) string {
	file, err := os.Stat(filename)

	if err != nil {
		fmt.Println(err)
	}

	return file.ModTime().Format("02/01/2006 15:04")
}

func FileSize(filename string) int64 {
	file, err := os.Stat(filename)

	if err != nil {
		fmt.Println(err)
	}

	return file.Size()
}

func FileInfos(filename string) gin.H {
	return gin.H{
		"name":        filename,
		"last_update": FileLastModification(filename),
		"size":        FileSize(filename),
	}
}

func Uptime() time.Duration {
	return time.Since(startTime)
}

func Init() {
	startTime = time.Now()
}
