package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
)

type volume struct {
	local     string `yaml:"local"`
	container string `yaml:"container"`
	rw        bool   `yaml:"rw"`
}

type env struct {
	key   string `yaml:"key"`
	value string `yaml:"key"`
}

type params struct {
	volumes []volume `yaml:"volumes"`
	env     []env    `yaml:"env"`
}

var configPath string
var arg string

func exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func main() {

	arg = path.Clean(arg)

	if len(configPath) == 0 {
		fmt.Println("No config path specified. exiting")
		os.Exit(1)
	}

	exists, err := exists(configPath)
	if !exists || err != nil {
		fmt.Println("Could not find config file at", configPath, "exiting")
		os.Exit(1)
	}

	arg1 := "there"
	arg2 := "are three"
	arg3 := "falcons"
	cmd := exec.Command("", arg1, arg2, arg3)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Print(string(stdout))
}
