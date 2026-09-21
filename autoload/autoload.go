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

// Envs contains auto loaded environmental variables.
var Envs map[string]string = nil

// FileNotFound optionally handles the case of file not found.
var FileNotFound func(dir, name string) = nil

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
		if notFound := FileNotFound; notFound != nil {
			notFound(dir, name)
		}
		return
	}
	envs, err := env.Load(path)
	if err != nil {
		panic(err)
	}
	Envs = envs
}
