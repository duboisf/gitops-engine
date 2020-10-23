package io

import (
	"os"

	"k8s.io/klog/v2/klogr"
)

var (
	TempDir   string
	NopCloser = NewCloser(func() error {
		return nil
	})
	Logger = klogr.New()
)

func init() {
	fileInfo, err := os.Stat("/dev/shm")
	if err == nil && fileInfo.IsDir() {
		TempDir = "/dev/shm"
	}
}

// DeleteFile is best effort deletion of a file
func DeleteFile(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	}
	_ = os.Remove(path)
}

type Closer interface {
	Close() error
}

type inlineCloser struct {
	close func() error
}

func (c *inlineCloser) Close() error {
	return c.close()
}

func NewCloser(close func() error) Closer {
	return &inlineCloser{close: close}
}

// Close is a convenience function to close a object that has a Close() method, ignoring any errors
// Used to satisfy errcheck lint
func Close(c Closer) {
	if err := c.Close(); err != nil {
		Logger.Error(err, "Failed to close %v", c)
	}
}
