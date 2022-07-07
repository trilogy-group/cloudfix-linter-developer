package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

/* Instruction to use logger:-
   Import ("github.com/trilogy-group/cloudfix-linter/logger") in your go file which is to be logged.

   Initialise logger in your main function of go file using "appLogger := logger.New()"
   After this a logs folder would be created if not present

   Writing to the log file using logger:
   To log info to the logger:- "appLogger.Info().Println("message")"
   To log Warning to the logger:- "appLogger.Warning().Println("message")"
   To log Error to the logger:- "appLogger.Error().Println("message")"

   results would get stored in logs/day-month-year.log files
*/
const (
	LogsDirpath = "logs"
)

type LogDir struct {
	LogDirectory string
}

func New() *LogDir {
	err := os.Mkdir(LogsDirpath, 0666)
	if err != nil {
		return nil
	}
	return &LogDir{
		LogDirectory: LogsDirpath,
	}
}

func SetLogFile() *os.File {
	year, month, day := time.Now().Date()
	fileName := fmt.Sprintf("%v-%v-%v.log", day, month.String(), year)
	filePath, _ := os.OpenFile(LogsDirpath+"/"+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	return filePath
}
func (l *LogDir) Info() *log.Logger {
	getFilePath := SetLogFile()
	return log.New(getFilePath, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (l *LogDir) Warning() *log.Logger {
	getFilePath := SetLogFile()
	return log.New(getFilePath, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func (l *LogDir) Error() *log.Logger {
	getFilePath := SetLogFile()
	return log.New(getFilePath, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
