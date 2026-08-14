package autoload

import "github.com/aileron-projects/go-env"

// FilePath is the environmental variable file path.
var FilePath = ".env"

func init() {
	if _, err := env.Load(FilePath); err != nil {
		panic(err)
	}
}
