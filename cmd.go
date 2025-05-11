package main

import (
	"fmt"
	"strings"

	"go.uber.org/zap/zapcore"
)

const helpText = `usage: ipapatch [-h/--help] --input <path> [--output <path] [--dylib <path>] [--inplace] [--noconfirm] [--version]

flags:
  --input path      the path to the ipa file to patch
  --output path     the path to the patched ipa file to create
  --dylib path      the path to the dylib to use instead of the embedded zxPluginsInject
  --inplace         takes priority over --output, use this to overwrite the input file
  --noconfirm       skip interactive confirmation when not using --inplace, overwriting a file that already exists, etc

info:
  -h, --help        show usage and exit
  --version         show version and exit`

const version = "ipapatch v1.0.0"

func AskInteractively(question string) bool {
	var reply string
	logger.Infof("%s [Y/n]", question)
	if _, err := fmt.Scanln(&reply); err != nil && err.Error() != "unexpected newline" {
		logger.Logw(zapcore.ErrorLevel, "couldnt scan reply", "err", err)
		return false
	}
	reply = strings.TrimSpace(reply)
	return reply == "" || reply == "y" || reply == "Y"
}
