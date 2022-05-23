package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"sync"
)

type Logger struct {
	*zap.Logger
}

var once sync.Once

func InitLogger(mod bool, level string) *Logger {
	var l Logger
	once.Do(func() {
		writer := getLogWriter(mod)
		encoder := getEncoder(mod)

		level, err := zapcore.ParseLevel(level)
		if err != nil {
			log.Fatalln(err)
		}

		core := zapcore.NewCore(encoder, writer, level)

		logger := zap.New(core)

		l = Logger{logger}
	})

	return &l
}

func getLogWriter(isDev bool) zapcore.WriteSyncer {
	if isDev == false {
		file, _ := os.Create("./logs/log.log")
		return zapcore.AddSync(file)
	}

	return zapcore.AddSync(os.Stdout)
}

func getEncoder(isDev bool) zapcore.Encoder {
	if isDev == false {
		return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}

	return zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
}
