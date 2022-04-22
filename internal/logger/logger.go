package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var Log zerolog.Logger

// InitLog sets up the desired level of logging, and
func InitLog(level int) error {

	switch level {
	case 5:
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case 4:
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case 3:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case 2:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case 1:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case 0:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		return fmt.Errorf("could not set log level for level %d", level)
	}

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s |", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	Log = zerolog.New(output).With().Timestamp().Logger()

	return nil
}
