package main

import (
	"io/fs"
	"time"
)

type zxPluginsInjectInfo struct{}

func (zxPluginsInjectInfo) Name() string {
	return "zxPluginsInject.dylib"
}

func (zxPluginsInjectInfo) Size() int64 {
	return 70368
}

func (zxPluginsInjectInfo) Mode() fs.FileMode {
	return 0755
}

func (zxPluginsInjectInfo) ModTime() time.Time {
	return time.Now()
}

func (zxPluginsInjectInfo) IsDir() bool {
	return false
}

func (zxPluginsInjectInfo) Sys() any {
	return nil
}
