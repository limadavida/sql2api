package logger

import (
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

var (
	once     sync.Once
	instance *log.Logger
)

func GetLogger() *log.Logger {
	once.Do(func() {
		instance = log.NewWithOptions(os.Stderr, log.Options{
			ReportCaller:    true,
			ReportTimestamp: true,
			TimeFormat:      time.Kitchen,
			Prefix:          "[SQL2API🛠️]",
		})
	})
	return instance
}
