package logging

import (
	"log"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

var once sync.Once

func InitLogger(InProd bool, level string) *Logger {
	var l Logger
	once.Do(
		func() {
			writer := getLogWriter(InProd)
			encoder := getEncoder(InProd)

			lvl, err := zapcore.ParseLevel(level)
			if err != nil {
				log.Fatalln(err)
			}

			core := zapcore.NewCore(encoder, writer, lvl)

			logger := zap.New(core)

			l = Logger{logger}
		},
	)

	return &l
}

func getLogWriter(InProd bool) zapcore.WriteSyncer {
	if InProd == true {
		file, _ := os.Create("./logs/log.log")
		return zapcore.AddSync(file)
	}

	return zapcore.AddSync(os.Stdout)
}

func getEncoder(InProd bool) zapcore.Encoder {
	if InProd == true {
		return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	}

	return zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
}
