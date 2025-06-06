package main

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/STARRY-S/zip"
	"howett.net/plist"
)

var (
	ErrNoPlist   = errors.New("No Info.plist found in ipa")
	ErrNoPlugins = errors.New("No plugins found")
)

func injectAll(args Args, tmpdir string) (map[string]string, error) {
	z, err := zip.OpenReader(args.Input)
	if err != nil {
		return nil, err
	}
	defer z.Close()

	plists, err := findPlists(z.File, args.PluginsOnly)
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

		zipPath := path.Join(path.Dir(p), execName)
		fsPath, err := extractToPath(z, tmpdir, zipPath)
		if err != nil {
			return nil, fmt.Errorf("Extract %s: %w", execName, err)
		}

		logger.Infof("Injecting: %s", execName)
		if err = injectLC(fsPath, lcName, tmpdir); err != nil {
			return nil, fmt.Errorf("Inject %s: %w", execName, err)
		}

		paths[fsPath] = zipPath
	}

	return paths, nil
}

func findPlists(files []*zip.File, pluginsOnly bool) (plists []string, err error) {
	plists = make([]string, 0, 10)

	for _, f := range files {
		if strings.Contains(f.Name, ".app/Watch") || strings.Contains(f.Name, ".app/WatchKit") || strings.Contains(f.Name, ".app/com.apple.WatchPlaceholder") {
			logger.Warnf("Watch app found: %s", path.Dir(f.Name))
			continue
		}
		if strings.HasSuffix(f.Name, ".appex/Info.plist") {
			plists = append(plists, f.Name)
			continue
		}
		if !pluginsOnly && strings.HasSuffix(f.Name, ".app/Info.plist") {
			plists = append(plists, f.Name)
			continue
		}
	}

	if len(plists) == 0 {
		if pluginsOnly {
			return nil, ErrNoPlugins
		}
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
	ff, err := os.OpenFile(output, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0755)
	if err != nil {
		return "", err
	}
	defer ff.Close()

	_, err = io.Copy(ff, f)
	return output, err
}

func appendFileToUpdater(ud *zip.Updater, path, zippedPath string) error {
	o, err := os.Open(path)
	if err != nil {
		return err
	}
	defer o.Close()

	fi, err := o.Stat()
	if err != nil {
		return err
	}

	return appendToUpdater(ud, zippedPath, fi, o)
}

func appendToUpdater(ud *zip.Updater, zippedPath string, fi fs.FileInfo, r io.Reader) error {
	hdr, err := zip.FileInfoHeader(fi)
	if err != nil {
		return err
	}

	hdr.Name = zippedPath
	hdr.Method = zip.Deflate

	w, err := ud.AppendHeader(hdr, zip.APPEND_MODE_OVERWRITE)
	if err != nil {
		return err
	}

	_, err = io.Copy(w, r)
	return err
}
