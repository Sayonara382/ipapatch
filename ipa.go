package main

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"howett.net/plist"
)

var ErrNoPlist = errors.New("no Info.plist found in ipa")

// key - path to file in provided tmpdir, now patched
// val - path inside ipa
func injectAll(args Args, tmpdir string) (map[string]string, error) {
	z, err := zip.OpenReader(args.Input)
	if err != nil {
		return nil, err
	}
	defer z.Close()

	plists, err := findPlists(z.File)
	if err != nil {
		return nil, err
	}
	paths := make(map[string]string, len(plists))

	lcName := "@rpath/"
	if args.Dylib == "" {
		lcName += "zxPluginsInject.dylib"
	} else {
		lcName += filepath.Base(args.Dylib)
	}

	for _, p := range plists {
		execName, err := getExecutableName(z, p)
		if err != nil {
			return nil, err
		}

		path := filepath.Join(filepath.Dir(p), execName)
		fsPath, err := extractToPath(z, tmpdir, path)
		if err != nil {
			return nil, fmt.Errorf("error extracting %s: %w", execName, err)
		}

		if err = injectLC(fsPath, lcName, tmpdir); err != nil {
			return nil, fmt.Errorf("couldnt inject into %s: %w", execName, err)
		}

		paths[fsPath] = path
	}

	return paths, nil
}

func findPlists(files []*zip.File) (plists []string, err error) {
	plists = make([]string, 0, 10)

	for _, f := range files {
		if strings.HasSuffix(f.Name, ".app/Info.plist") || strings.HasSuffix(f.Name, ".appex/Info.plist") {
			plists = append(plists, f.Name)
		}
	}

	if len(plists) == 0 {
		return nil, ErrNoPlist
	}
	return plists, nil
}

func getExecutableName(z *zip.ReadCloser, plistName string) (string, error) {
	f, err := z.Open(plistName)
	if err != nil {
		return "", err
	}
	defer f.Close()

	contents, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	var pl struct {
		Executable string `plist:"CFBundleExecutable"`
	}
	if _, err = plist.Unmarshal(contents, &pl); err != nil {
		return "", err
	}

	return pl.Executable, nil
}

func extractToPath(z *zip.ReadCloser, dir, name string) (string, error) {
	f, err := z.Open(name)
	if err != nil {
		return "", err
	}
	defer f.Close()

	output := filepath.Join(dir, filepath.Base(name))
	ff, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0700)
	if err != nil {
		return "", err
	}
	defer ff.Close()

	_, err = io.Copy(ff, f)
	return output, err
}
