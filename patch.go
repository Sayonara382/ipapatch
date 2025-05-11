package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mholt/archives"
)

// Patch patches the executable and all plugins.
func Patch(args Args) error {
	tmpdir, err := os.MkdirTemp(".", ".ipapatch-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpdir)

	//

	logger.Info("extracting and injecting..")
	paths, err := injectAll(args, tmpdir)
	if err != nil {
		return fmt.Errorf("error injecting: %w", err)
	}

	if args.Output != args.Input {
		logger.Info("copying input to output..")
		if err = copyfile(args.Input, args.Output); err != nil {
			return fmt.Errorf("failed to copy input to output: %w", err)
		}
	}

	zipArgs := make([]string, 0, len(paths)+2)
	zipArgs = append(zipArgs, "-d", args.Output)
	for _, val := range paths {
		zipArgs = append(zipArgs, val)
	}
	appName := strings.Split(zipArgs[2], "/")[1]

	err = exec.Command("zip", zipArgs...).Run()
	if err != nil {
		return fmt.Errorf("error deleting from zipfile: %w", err)
	}

	//

	logger.Info("adding files back to ipa..")

	o, err := os.OpenFile(args.Output, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer o.Close()

	ctx := context.Background()
	ffd, err := archives.FilesFromDisk(ctx, nil, paths)
	if err != nil {
		return fmt.Errorf("FilesFromDisk err: %w", err)
	}

	if args.Dylib != "" {
		ffdDylib, err := archives.FilesFromDisk(
			ctx, nil,
			map[string]string{
				args.Dylib: fmt.Sprintf("Payload/%s/Frameworks/%s", appName, filepath.Base(args.Dylib)),
			},
		)
		if err != nil {
			return fmt.Errorf("FilesFromDisk err (--dylib): %w", err)
		}

		ffd = append(ffd, ffdDylib...)
	} else {
		ffd = append(ffd, archives.FileInfo{
			FileInfo:      zxPluginsInjectInfo{},
			NameInArchive: fmt.Sprintf("Payload/%s/Frameworks/zxPluginsInject.dylib", appName),
			Open:          zxPluginsInjectOpen,
		})
	}

	if err = (archives.Zip{}).Insert(ctx, o, ffd); err != nil {
		return fmt.Errorf("error inserting files back into ipa: %w", err)
	}

	return nil
}

func copyfile(from, to string) error {
	f1, err := os.Open(from)
	if err != nil {
		return err
	}
	defer f1.Close()

	f2, err := os.OpenFile(to, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f2.Close()

	_, err = io.Copy(f2, f1)
	return err
}
