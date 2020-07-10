// Logging defines a logger interface to log to standard out and a local file.
package dicebae

import (
	"fmt"
	"log"
	"os"
	"path"
)

func (db *diceBae) initLogger(logDir string) error {
	logFile := "bae.log"
	logPath := path.Join(logDir, logFile)
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open the logfile %q: %v", logPath, err)
	}
	db.logFile = f
	db.logger = log.New(f, "", log.LstdFlags)
	db.LogInfo("Logging to %q", logPath)
	return nil
}

func (db *diceBae) LogInfo(format string, args ...interface{}) {
	db.logger.SetPrefix("info:")
	db.logger.Printf(format, args...)
	fmt.Printf(format+"\n", args...)
}

func (db *diceBae) LogError(format string, args ...interface{}) {
	db.logger.SetPrefix("awww:")
	db.logger.Printf(format, args...)
	fmt.Errorf(format+"\n", args...)
}
