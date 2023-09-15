package logger

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"tugboat/fs"
)

type LogManager struct {
	fsMan *fs.FSManager
	logLevel string
	logLocation string
}

func CreateLogManager(fsm *fs.FSManager) *LogManager {
	lm := &LogManager{
		fsMan: fsm,
		logLevel: "info",
		logLocation: "/logs",
	}
	return lm
}

func (logman *LogManager) SetLevel(level string){
	logman.logLevel = level
}
func (logman *LogManager) SetLogLocation(location string){
	re := regexp.MustCompile(`memfs:/(/.*)`)
	matches := re.FindStringSubmatch(location)
	if len(matches) == 2 {
		logman.logLocation = matches[1]
	}else{
		log.Printf("Logger: invalid format of location: %v\n", matches)
	}
}
func (logman *LogManager) GetLogLocation() string{
	return logman.logLocation
}

func (logman *LogManager) Log(facility string, severity string, message string) bool{

	logfs := logman.fsMan.GetLogFS()


	_, err := logman.fsMan.CheckIfMemDirExists(logfs, logman.logLocation)
	if err != nil {
		log.Printf("Logger: unable to verify log directory: %v\n", err)
		log.Printf("Logger: attempting to create log directory\n")

		err = logman.fsMan.CreateDir(logfs, logman.logLocation)
		if err != nil {
			log.Printf("Logger: unable to create top level logging %v\n", err)
			return false
		}
	}


	timeNowFmt := time.Now().Format("2006-01-02 15:04:05 MST")
	logFileLocation := strings.Join([]string{logman.logLocation,facility}, "/")
	logMessage := fmt.Sprintf("%-22s | %-8s| %-5s| %s\n", timeNowFmt, facility, severity, message)

	//log.Printf("Logger: Logging to : %s with %s\n", logFileLocation, logMessage)
	_, err = logman.fsMan.WriteFile("log", logFileLocation, []byte(logMessage))
	if err != nil {
		log.Printf("Logger: unable to create logging for: %s : ", facility, err)
		return false
	}
	return true
}