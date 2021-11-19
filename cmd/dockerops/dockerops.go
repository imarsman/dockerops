package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/alexflint/go-arg"
	"gopkg.in/yaml.v2"
)

// Parameters a set of parameters
type Parameters struct {
	Volumes []Volume `yaml:"volumes"`
	Env     []Env    `yaml:"environment"`
}

// Volume a set of volume arguments
type Volume struct {
	Local     string `yaml:"local"`
	Container string `yaml:"container"`
	Rw        bool   `yaml:"rw"`
}

// Env a set of environment arguments
type Env struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// var configPath string
// var arg string

func exists(name string) (exists bool, err error) {
	_, err = os.Stat(name)
	if err == nil {
		return true, err
	}
	if errors.Is(err, os.ErrNotExist) {
		return exists, nil
	}
	return exists, err
}

func main() {

	var args struct {
		ConfigPath string `arg:"-C" help:"config path - defaults to [dockeropps dir]/dockerops.yml"`
		Call       string `arg:"positional" help:"call to ops to make"`
	}

	arg.MustParse(&args)

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePath := filepath.Dir(ex)

	if len(args.ConfigPath) == 0 {
		args.ConfigPath = filepath.Join(exePath, "dockerops.yml")
		fmt.Println("config path", args.ConfigPath)
	}

	args.ConfigPath = path.Clean(args.ConfigPath)

	exists, err := exists(args.ConfigPath)
	if !exists || err != nil {
		fmt.Println("Could not find config file at", args.ConfigPath, "exiting")
		os.Exit(1)
	}

	params := Parameters{}
	cfgBytes, err := os.ReadFile(args.ConfigPath)
	// fmt.Println("bytes", string(cfgBytes))
	err = yaml.Unmarshal(cfgBytes, &params)
	if err != nil {
		fmt.Println("Error reading config file", args.ConfigPath, "exiting")
		os.Exit(1)
	}

	dockerArgs := make([]string, 0, len(params.Env)+len(params.Volumes)+5)
	dockerArgs = append(dockerArgs, "--platform linux/amd64")

	for _, v := range params.Env {
		key, val := strings.TrimSpace(v.Key), strings.TrimSpace(v.Value)

		if key == "" || val == "" {
			continue
		}
		dockerArgs = append(dockerArgs, fmt.Sprintf("--env \"%s=%s\"", key, val))
	}
	for _, v := range params.Volumes {
		local, container := strings.TrimSpace(v.Local), strings.TrimSpace(v.Container)

		rw := "ro"
		if v.Rw {
			rw = "rw"
		}
		dockerArgs = append(dockerArgs, fmt.Sprintf("--volume \"%s:%s:%s\"", local, container, rw))
	}
	dockerArgs = append(dockerArgs, "--name nanos")
	dockerArgs = append(dockerArgs, "nanos:latest")
	dockerArgs = append(dockerArgs, "/app/ops")
	dockerArgs = append(dockerArgs, args.Call)

	cmd := exec.Command("docker run", dockerArgs...)
	fmt.Printf("command \"%s\"\n", cmd.String())

	// fmt.Print(string(stdout))
}
