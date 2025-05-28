package main

import (
	"fmt"
	"strings"

	"go.uber.org/zap/zapcore"
)

const helpText = `Usage: ipapatch [-h/--help] --input <path> [--output <path] [--dylib <path>] [--inplace] [--noconfirm] [--plugins-only] [--version]

Flags:
  --input path      The path to the ipa file to patch
  --output path     The path to the patched ipa file to create
  --dylib path      The path to the dylib to use instead of the embedded zxPluginsInject
  --inplace         Takes priority over --output, use this to overwrite the input file
  --noconfirm       Skip interactive confirmation when not using --inplace, overwriting a file that already exists, etc
  --plugins-only    Only inject into plugin binaries (not the main executable)

Info:
  -h, --help        Show usage and exit
  --version         Show version and exit`

type Args struct {
	Input       string `arg:"--input,required"`
	Output      string `arg:"--output"`
	Dylib       string `arg:"--dylib"`
	InPlace     bool   `arg:"--inplace"`
	NoConfirm   bool   `arg:"--noconfirm"`
	PluginsOnly bool   `arg:"--plugins-only"`
}

func (Args) Version() string {
	return "ipapatch v1.1.2"
}

func AskInteractively(question string) bool {
	var reply string
	logger.Infof("%s [Y/n]", question)
	if _, err := fmt.Scanln(&reply); err != nil && err.Error() != "unexpected newline" {
		logger.Warnw("Input scan failed", "err", err)
		return false
	}
	reply = strings.TrimSpace(reply)
	return reply == "" || reply == "y" || reply == "Y"
}
