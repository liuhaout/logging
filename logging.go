package logging

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Log fields
type Logging struct {
	LogDir, LogName, LogPrefix string
	LogSize, LogBackupNum      int64
	LogFlags                   int
}

var Log *Logging

// init logfile config, or redefine by Log afterwards
func init() {
	dir, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	Log = &Logging{
		LogDir:    dir,
		LogName:   "UnnamedLog.log",
		LogPrefix: "",

		LogSize:      1024,
		LogBackupNum: 1,
		LogFlags:     1 | 2,
	}
}

// Log level
func (l *Logging) EMEG(v interface{}) {
	handle("EMEG", v)
}
func (l *Logging) ALERT(v interface{}) {
	handle("ALERT", v)
}
func (l *Logging) CRIT(v interface{}) {
	handle("CRIT", v)
}
func (l *Logging) ERR(v interface{}) {
	handle("ERR", v)
}
func (l *Logging) WARNING(v interface{}) {
	handle("WARNING", v)
}
func (l *Logging) NOTICE(v interface{}) {
	handle("NOTICE", v)
}
func (l *Logging) INFO(v interface{}) {
	handle("INFO", v)
}
func (l *Logging) DEBUG(v interface{}) {
	handle("DEBUG", v)
}

func handle(level string, v interface{}) {
	flushlogFile()
	f, err := os.OpenFile(fmt.Sprintf("%v/%v", Log.LogDir, Log.LogName), os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	logger := log.New(f, Log.LogPrefix, Log.LogFlags)
	logger.Printf("[%v] %v", level, v)
}

func flushlogFile() {
	logFilePath := fmt.Sprintf("%v/%v", Log.LogDir, Log.LogName)
	fileInfo, err := os.Stat(logFilePath)
	// logfile not exist
	if err != nil {
		return
	}
	// logfile exist
	// expected logfile size
	if fileInfo.Size() <= Log.LogSize {
		return
	}
	// logfile exist
	// not expected logfile size
	// logfile backup number = 1, will keep logfile size and flush content
	if Log.LogBackupNum == 1 {
		data, err := ioutil.ReadFile(logFilePath)
		if err != nil {
			log.Panic(err)
		}
		tmp := strings.Join(strings.Split(string(data), "\n")[1:], "\n")
		ioutil.WriteFile(logFilePath, []byte(tmp), 0755)
		return
	}
	// logfile exist
	// not expected logfile size
	// logfile backup number > 1, will turn backup logfile
	for i := Log.LogBackupNum - 2; i > -1; i-- {
		if i == 0 {
			if _, err := os.Stat(logFilePath); err == nil {
				os.Rename(logFilePath, fmt.Sprintf("%v/%v.%v", Log.LogDir, Log.LogName, i+1))
				os.Truncate(logFilePath, 0)
				continue
			}
		}
		if _, err := os.Stat(fmt.Sprintf("%v/%v.%v", Log.LogDir, Log.LogName, i)); err == nil {
			os.Rename(fmt.Sprintf("%v/%v.%v", Log.LogDir, Log.LogName, i), fmt.Sprintf("%v/%v.%v", Log.LogDir, Log.LogName, i+1))
		}
	}
}
