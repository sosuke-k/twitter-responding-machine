package trm

import (
	"fmt"
	"log"
	"os"
	"path"
)

type TRMLogger struct {
	logPath string
}

var instance *TRMLogger

// GetLogger return singlton instance
func GetLogger() *TRMLogger {
	if instance == nil {
		pwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stdout, "error opening file: %v", err)
			return nil
		}
		logPath := path.Join(pwd, "trm.log")
		instance = &TRMLogger{logPath: logPath}
	}
	return instance
}

func (mylogger *TRMLogger) Println(s string) {
	logf, err := os.OpenFile(mylogger.logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stdout, "error opening file: %v", err)
	}
	defer logf.Close()

	log.SetOutput(logf)
	log.Println(s)
}

func (mylogger *TRMLogger) Fatalln(e error) {
	logf, err := os.OpenFile(mylogger.logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stdout, "error opening file: %v", err)
	}
	defer logf.Close()

	log.SetOutput(logf)
	log.Fatalln(e)
}

func (mylogger *TRMLogger) Printf(format string, e interface{}) {
	logf, err := os.OpenFile(mylogger.logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stdout, "error opening file: %v", err)
	}
	defer logf.Close()

	log.SetOutput(logf)
	log.Printf(format, e)
}
