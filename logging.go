package logging

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// log fields
type Logging struct {
	LogDir, LogName, LogPrefix string
	LogSize, LogBackupNum      int64
	LogFlags                   int
}

var Logger Logging

// init logfile config, or redefine by Log afterwards
func init() {
	dir, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	Logger = Logging{
		LogDir:    dir,
		LogName:   "Unnamed.log",
		LogPrefix: "",

		LogSize:      1024,
		LogBackupNum: 1,
		LogFlags:     1 | 2,
	}
}

// log level
func (logger Logging) EMEG(v interface{}) {
	handle("EMEG", logger, v)
}

func (logger Logging) ALERT(v interface{}) {
	handle("ALERT", logger, v)
}
func (logger Logging) CRIT(v interface{}) {
	handle("CRIT", logger, v)
}
func (logger Logging) ERR(v interface{}) {
	handle("ERR", logger, v)
}
func (logger Logging) WARNING(v interface{}) {
	handle("WARNING", logger, v)
}
func (logger Logging) NOTICE(v interface{}) {
	handle("NOTICE", logger, v)
}
func (logger Logging) INFO(v interface{}) {
	handle("INFO", logger, v)
}
func (logger Logging) DEBUG(v interface{}) {
	handle("DEBUG", logger, v)
}

func handle(level string, logger Logging, v interface{}) {
	logFilePath := fmt.Sprintf("%v/%v", logger.LogDir, logger.LogName)
	flushlogFile(logger, logFilePath)
	f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	newlogger := log.New(f, logger.LogPrefix, logger.LogFlags)
	newlogger.Printf("[%v] %v", level, v)
}

func flushlogFile(logger Logging, logFilePath string) {
	fileInfo, err := os.Stat(logFilePath)
	dif := fileInfo.Size() - logger.LogSize
	// logfile not exist
	if err != nil {
		return
	}
	// logfile exist
	// expected logfile size
	if dif < 0 {
		return
	}
	// logfile exist
	// not expected logfile size
	// logfile backup number = 1, will keep logfile size and flush content
	if logger.LogBackupNum == 1 {
		data, err := ioutil.ReadFile(logFilePath)
		if err != nil {
			log.Panic(err)
		}
		dataList := strings.Split(string(data), "\n")
		if len(dataList) < 2 {
			return
		}
		for i, _ := range dataList {
			if int64(len(strings.Join(dataList[:i+1], "\n"))) > dif {
				ioutil.WriteFile(logFilePath, []byte(strings.Join(dataList[i+1:], "\n")), 0755)
				break
			}
		}
		return
	}
	// logfile exist
	// not expected logfile size
	// logfile backup number > 1, will turn backup logfile
	for i := logger.LogBackupNum - 2; i > -1; i-- {
		if i == 0 {
			if _, err := os.Stat(logFilePath); err == nil {
				os.Rename(logFilePath, fmt.Sprintf("%v.%v", logFilePath, i+1))
				os.Truncate(logFilePath, 0)
				continue
			}
		}
		if _, err := os.Stat(fmt.Sprintf("%v.%v", logFilePath, i)); err == nil {
			os.Rename(fmt.Sprintf("%v.%v", logFilePath, i), fmt.Sprintf("%v.%v", logFilePath, i+1))
		}
	}
}
