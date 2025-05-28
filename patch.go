package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/STARRY-S/zip"
)

func Patch(args Args) error {
	tmpdir, err := os.MkdirTemp(".", ".ipapatch-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpdir)

	logger.Info("Extracting and injecting")
	paths, err := injectAll(args, tmpdir)
	if err != nil {
		return fmt.Errorf("Injecting: %w", err)
	}

	if args.Output != args.Input {
		logger.Info("Copying input to output")
		if err = copyfile(args.Input, args.Output); err != nil {
			return fmt.Errorf("Copy input to output: %w", err)
		}
	}

	zipArgs := make([]string, 0, len(paths)+2)
	zipArgs = append(zipArgs, "-d", args.Output)
	for _, val := range paths {
		zipArgs = append(zipArgs, val)
	}
	appName := strings.Split(zipArgs[2], "/")[1]

	if len(paths) > 0 {
		logger.Info("Removing old files from IPA")
		if err := removeFilesFromZip(args.Output, paths); err != nil {
			return fmt.Errorf("Delete from zipfile: %w", err)
		}
	}

	logger.Info("Adding files back to IPA")
	o, err := os.OpenFile(args.Output, os.O_RDWR, 0)
	if err != nil {
		return err
	}
	defer o.Close()

	ud, err := zip.NewUpdater(o)
	if err != nil {
		return err
	}
	defer ud.Close()

	for sysPath, zippedPath := range paths {
		if err = appendFileToUpdater(ud, sysPath, zippedPath); err != nil {
			return err
		}
	}

	if args.Dylib != "" {
		return appendFileToUpdater(ud, args.Dylib, fmt.Sprintf("Payload/%s/Frameworks/%s", appName, filepath.Base(args.Dylib)))
	}

	zxpi, err := zxPluginsInject.Open("resources/zxPluginsInject.dylib")
	if err != nil {
		return err
	}
	defer zxpi.Close()

	return appendToUpdater(
		ud,
		fmt.Sprintf("Payload/%s/Frameworks/zxPluginsInject.dylib", appName),
		zxPluginsInjectInfo{},
		zxpi,
	)
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

func removeFilesFromZip(zipPath string, removePaths map[string]string) error {
	orig, err := os.Open(zipPath)
	if err != nil {
		return err
	}
	defer orig.Close()

	stat, err := orig.Stat()
	if err != nil {
		return err
	}

	r, err := zip.NewReader(orig, stat.Size())
	if err != nil {
		return err
	}

	tmpZip, err := os.CreateTemp(".", ".ipapatch-tmpzip-*")
	if err != nil {
		return err
	}
	defer func() {
		tmpZip.Close()
		os.Remove(tmpZip.Name())
	}()

	w := zip.NewWriter(tmpZip)
	defer w.Close()

	removeSet := make(map[string]struct{}, len(removePaths))
	for _, p := range removePaths {
		removeSet[p] = struct{}{}
	}

	for _, f := range r.File {
		if _, shouldRemove := removeSet[f.Name]; shouldRemove {
			continue
		}
		fr, err := f.Open()
		if err != nil {
			return err
		}
		hdr := f.FileHeader
		fw, err := w.CreateHeader(&hdr)
		if err != nil {
			fr.Close()
			return err
		}
		_, err = io.Copy(fw, fr)
		fr.Close()
		if err != nil {
			return err
		}
	}

	w.Close()
	tmpZip.Close()

	orig.Close()
	if err := os.Rename(tmpZip.Name(), zipPath); err != nil {
		return err
	}
	return nil
}