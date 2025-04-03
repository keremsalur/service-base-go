package logger

import (
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Logger struct {
	zerolog.Logger
}

var globalLogger *Logger
var globalLogData map[string]interface{}

func newLogger(loggerName string) *Logger {
	output := zerolog.ConsoleWriter{
		TimeFormat: time.RFC3339,
		Out:        os.Stdout, // Çıktının nereye yazılacağı
		NoColor:    false,     // Renkli çıktı oluşturmasın
	}

	logger := zerolog.New(output).With().Timestamp().Str("logger", loggerName).Logger()

	level := os.Getenv("LOG_LEVEL")

	if level != "" {
		if parsedLevel, err := zerolog.ParseLevel(level); err == nil {
			zerolog.SetGlobalLevel(parsedLevel)
		}
	}

	return &Logger{logger}
}

func InitLogger(loggerName string) {
	globalLogger = newLogger(loggerName)
}

func GetLogger() *Logger {
	if globalLogger == nil {
		panic("Global logger oluşturulmamış! Önce InitGlobalLogger fonksiyonunu çağırın")
	}
	return globalLogger
}

func GetGlobalLogData() map[string]interface{} {
	globalLogData = make(map[string]interface{})
	return globalLogData
}

func PushLog(logger *Logger, level zerolog.Level, message string, fields map[string]interface{}) {

	event := logger.WithLevel(level).
		Str("detected_level", level.String()).
		Str("method", GetMethodName(2)) // 2, PushLog fonksiyonunu çağıranı temsil eder

	// Dinamik alanları ekle
	for key, value := range fields {
		switch v := value.(type) {
		case string:
			event = event.Str(key, v)
		case int:
			event = event.Int(key, v)
		case bool:
			event = event.Bool(key, v)
		case float64:
			event = event.Float64(key, v)
		case nil:
			event = event.Interface(key, nil)
		case time.Time:
			event = event.Time(key, v) // Zaman tipi için
		case time.Duration:
			event = event.Dur(key, v) // Süre tipi için
		default:
			event = event.Interface(key, v) // Diğer türler için
		}
		// Mesajı logla

	}
	event.Msg(message)
}

func GetMethodName(callerParent int) string {
	pc, _, _, _ := runtime.Caller(callerParent) // 1, çağıran fonksiyonu temsil eder
	funcName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(funcName, ".")
	return parts[len(parts)-2] + "." + parts[len(parts)-1]
}
