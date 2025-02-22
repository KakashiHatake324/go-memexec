//go:build !linux
// +build !linux

package memexec

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func open(b []byte, name string) (*os.File, error) {
	pattern := name
	if runtime.GOOS == "windows" {
		pattern = fmt.Sprintf("%s.exe", name)
	}

	if file, err := os.Stat(filepath.Join(tempDir(), pattern)); err == nil {
		os.Remove(file.Name())
	}
	f, err := os.Create(filepath.Join(tempDir(), pattern))
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = clean(f)
		}
	}()
	if err = os.Chmod(f.Name(), 0o500); err != nil {
		return nil, err
	}
	if _, err = f.Write(b); err != nil {
		return nil, err
	}
	if err = f.Close(); err != nil {
		return nil, err
	}
	return f, nil
}

func clean(f *os.File) error {
	return os.Remove(f.Name())
}

func tempDir() string {
	os.TempDir()
	dir := os.TempDir()
	if dir == "" {
		if runtime.GOOS == "android" {
			dir = "/data/local/tmp"
		} else {
			dir = "/tmp"
		}
	}
	return dir
}
