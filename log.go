package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger

func init() {
	config := zap.NewProductionConfig()
	config.Encoding = "console"
	config.EncoderConfig.TimeKey = "" // omit time from logs
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.ConsoleSeparator = " "
	config.DisableCaller = true
	config.DisableStacktrace = true
	l, err := config.Build()
	if err != nil {
		panic(err)
	}
	logger = l.Sugar()
}
