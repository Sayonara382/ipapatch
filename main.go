package main

import (
	"embed"
	"errors"
	"fmt"
	"os"

	"github.com/alexflint/go-arg"
)

var zxPluginsInject embed.FS

func main() {
	var args Args
	if err := arg.Parse(&args); err != nil {
		if errors.Is(err, arg.ErrHelp) {
			fmt.Println(helpText)
			return
		} else if errors.Is(err, arg.ErrVersion) {
			fmt.Println(args.Version())
			return
		} else if args.Input == "" {
			fmt.Println(helpText)
			logger.Error("Input file is required")
			return
		}
		logger.Fatalf("Argument parsing error: %v (see --help for usage)", err)
	}

	if args.InPlace {
		logger.Info("In-place mode: will overwrite input file")
		args.Output = args.Input
	}
	if args.Output == "" {
		if args.NoConfirm {
			logger.Fatal("Neither --output nor --inplace specified")
		}
		if !AskInteractively("--inplace not specified, overwrite the input?") {
			logger.Warn("User declined to overwrite input file")
			return
		}
		args.Output = args.Input
	} else {
		_, err := os.Stat(args.Output)
		if err == nil && !args.InPlace {
			if args.NoConfirm {
				logger.Info("--output already exists, overwriting")
			} else if !AskInteractively("--output already exists, overwrite?") {
				logger.Warn("User declined to overwrite output file")
				return
			}
		}
	}

	if args.Dylib != "" {
		_, err := os.Stat(args.Dylib)
		if os.IsNotExist(err) {
			logger.Fatalw("Dylib path does not exist", "path", args.Dylib)
		}
	}

	if err := Patch(args); err != nil {
		logger.Errorw("Patch failed", "error", err)
		os.Exit(1)
	}
}