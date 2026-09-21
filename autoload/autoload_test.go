package autoload

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/aileron-projects/go-tester"
)

func TestLoadIfExist(t *testing.T) {
	t.Parallel()
	t.Run("file exist", func(t *testing.T) {
		tmp := t.TempDir()
		name := "application.env"
		err := os.WriteFile(filepath.Join(tmp, name), []byte("FOO=bar"), os.ModePerm)
		tester.AssertEqual(t, nil, err)
		loadIfExist(tmp, name)
		tester.AssertEqual(t, "bar", os.Getenv("FOO"))
		os.Unsetenv("FOO")
	})
	t.Run("file not exist", func(t *testing.T) {
		tmp := t.TempDir()
		name := "application.env"
		loadIfExist(tmp, name)
	})
	t.Run("file not found handler", func(t *testing.T) {
		tmp := t.TempDir()
		name := "application.properties"
		nfDir, nfName := "", ""
		FileNotFound = func(d, n string) {
			nfDir, nfName = d, n
		}
		defer func() { FileNotFound = nil }()
		loadIfExist(tmp, name)
		tester.AssertEqual(t, tmp, nfDir)
		tester.AssertEqual(t, name, nfName)
	})
	t.Run("panic", func(t *testing.T) {
		tmp := t.TempDir()
		name := "application.env"
		err := os.WriteFile(filepath.Join(tmp, name), []byte("=test"), os.ModePerm)
		tester.AssertEqual(t, nil, err)
		tester.AssertPanic(t, func() {
			loadIfExist(tmp, name)
		})
	})
}
