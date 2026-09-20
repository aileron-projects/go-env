package autoload

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/aileron-projects/go-env"
)

var (
	EnvDir         = "./"
	EnvPath        = "application.env"
	ProfilePattern = "application-${profile}.env"
	Profile        = os.Getenv("ENV_ACTIVE_PROFILE")
)

func init() {
	loadIfExist(EnvDir, EnvPath)
	if Profile != "" {
		path := strings.ReplaceAll(ProfilePattern, "${profile}", Profile)
		loadIfExist(EnvDir, path)
	}
}

func loadIfExist(dir, name string) {
	path := filepath.Join(dir, name)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return
	}
	if _, err := env.Load(path); err != nil {
		panic(err)
	}
}
