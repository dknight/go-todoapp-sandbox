package lib

import (
	"fmt"
	"log"
	"os"
	"time"
)

// TODO error, warn, notice
func NewLogger(logDir string) *log.Logger {
	logfile := fmt.Sprintf(
		"%s/%s.log", logDir, time.Now().Format("2006-01-02"),
	)
	fp, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	return log.New(fp, "", log.LstdFlags)
}
