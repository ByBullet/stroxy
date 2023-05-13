package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
)

var prodLogger *zap.Logger

func consoleEncoder(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
	encoderConfig := zap.NewDevelopmentEncoderConfig() // 同样使用 Development 进行默认设置
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("06/01/02 03:04pm")

	return zapcore.NewConsoleEncoder(encoderConfig), nil
}

func Init() {

	err := zap.RegisterEncoder("my-console", consoleEncoder)

	if err != nil {
		log.Fatal(err)
	}
	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      false,
		Encoding:         "my-console",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "debug":
			cfg.Development = true
		default:
			log.Fatalf("未知参数 %s,参数只能是release/debug\n", os.Args[1])
		}
	}

	prodLogger, err = cfg.Build()

	if err != nil {
		log.Fatal(err)
	}
	defer prodLogger.Sync()
}

func PROD() *zap.Logger {
	return prodLogger
}
