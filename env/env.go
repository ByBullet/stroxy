package env

import (
	"flag"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Env struct {
	Mode    string
	WorkDir string
}

var env = Env{
	Mode: "release",
}

func LoadEnv() {
	flag.StringVar(&env.Mode, "mode", "release", "-mode=[mode]")
	flag.Parse()

	if env.Mode == "debug" {
		if env.WorkDir == "" {
			_, env.WorkDir, _, _ = runtime.Caller(0)
			env.WorkDir = env.WorkDir[:strings.LastIndex(env.WorkDir, "env/env.go")-1]
		}
	} else {
		if ex, err := os.Executable(); err == nil {
			env.WorkDir = filepath.Dir(ex)
			return
		}
		env.WorkDir = "./"
	}
}

func GetEnv() Env {
	return env
}
