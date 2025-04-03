package config

import (
	"os"
	"service-base-go/pkg/logger"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

type Config struct {
	LogLevel    string
	DatabaseUrl string
	LogFileUrl  string
	Port        string
	AppName     string

	// Timeoutlar, limitler,
	IdleTimeout  int
	ReadTimeout  int
	WriteTimeout int
}

func getEnvAsInt(name string, defaultValue int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func LoadConfig() Config {
	if err := godotenv.Load(".env"); err != nil {
		panic(".env dosyası okunamadı")
	}
	return Config{
		IdleTimeout:  getEnvAsInt("IDLE_TIMEOUT", 1),
		ReadTimeout:  getEnvAsInt("READ_TIMEOUT", 1),
		WriteTimeout: getEnvAsInt("WRITE_TIMEOUT", 1),
		LogLevel:     os.Getenv("LOG_LEVEL"),
		Port:         os.Getenv("PORT"),
		AppName:      os.Getenv("APP_NAME"),
		DatabaseUrl:  os.Getenv("DATABASE_URL"),
		LogFileUrl:   os.Getenv("LOG_FILE_URL"),
	}
}

func SetupLoggingToFile(logFilePath string) (*os.File, error) {
	// Log dosyasını açma
	logData := logger.GetGlobalLogData()
	logData["class"] = "Config"
	logData["logFilePath"] = logFilePath
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logData["error"] = err
		logger.PushLog(logger.GetLogger(), zerolog.ErrorLevel, "Log dosyası açılmadı", logData)
		return nil, err
	}
	logger := logger.GetLogger()
	// Zerolog'u dosyaya yazacak şekilde yapılandırma
	if logFile != nil {
		logger.Logger = logger.Logger.Output(zerolog.ConsoleWriter{Out: os.Stdout}) // Geliştirme ortamından çıkınca out parametresini logfile yap
	} else {
		return nil, err
	}

	return logFile, nil
}
