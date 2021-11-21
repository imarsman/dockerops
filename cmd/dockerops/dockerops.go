package main

import (
	"errors"
	"fmt"
	"log"
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
}

// Env a set of environment arguments
type Env struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// Check if a file exists
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

	// Define args for go-args
	var args struct {
		ConfigPath string `arg:"-c" help:"config path - defaults to [dockeropps dir]/dockerops.yml"`
		// The rest of the call args as a slice
		// NOTE:
		// If there are any flagged parameters to be passed to ops the whole ops
		// call will need to be surrounded in quotes to ensure that the arg
		// parser treats them as part of the Call arg.
		Env     []string `arg:"-e,separate" help:"Set environment variable as key=val"`
		Verbose bool     `arg:"-v" help:"print out what is being handled and done"`
		Call    []string `arg:"positional" help:"call to ops - surround with quotes"`
	}

	arg.MustParse(&args)

	// Get path for app for when no config arg supplied and a search is to be
	// made for the config file.
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exePath := filepath.Dir(ex)

	if len(args.ConfigPath) == 0 {
		args.ConfigPath = filepath.Join(exePath, "dockerops.yml")
	}

	args.ConfigPath = path.Clean(args.ConfigPath)

	exists, err := exists(args.ConfigPath)
	if !exists || err != nil {
		log.Fatalf("could not find config file at %s - exiting", args.ConfigPath)
	}

	// Read in config file parameters
	params := Parameters{}
	cfgBytes, err := os.ReadFile(args.ConfigPath)
	if err != nil {
		log.Fatalf("error reading config file %s - exiting", args.ConfigPath)
	}

	err = yaml.Unmarshal(cfgBytes, &params)
	if err != nil {
		log.Fatalf("error reading config file %s - exiting", args.ConfigPath)
	}

	dockerArgs := make([]string, 0, len(params.Env)+len(params.Volumes)+10)

	// Start building call to docker
	dockerArgs = append(dockerArgs, "run")
	// When double-hyphen arguments are handled along with their values it
	// messes things up.
	dockerArgs = append(dockerArgs, "--platform")
	dockerArgs = append(dockerArgs, "linux/amd64")

	// If there are any config environment parameters handle them
	for _, v := range params.Env {
		key, val := strings.TrimSpace(v.Key), strings.TrimSpace(v.Value)

		if key == "" || val == "" {
			continue
		}
		if args.Verbose {
			log.Printf("setting environment key %s to val %s\n", key, val)
		}
		os.Setenv(key, val)
	}

	// Handle arg environment pairs
	if len(args.Env) > 0 {
		for _, arg := range args.Env {
			envSet := strings.TrimSpace(arg)
			parts := strings.Split(envSet, ",")
			if len(parts) >= 1 {
				for _, p := range parts {
					kv := strings.Split(p, "=")
					key, val := strings.TrimSpace(kv[0]), strings.TrimSpace(kv[1])

					if args.Verbose {
						log.Printf("setting arg environment key %s to val %s\n", key, val)
					}
					os.Setenv(key, val)
				}
			}
		}
	}

	// Handle setting volume mounts to container as defined in config file
	for _, v := range params.Volumes {
		local, container := strings.TrimSpace(v.Local), strings.TrimSpace(v.Container)

		if args.Verbose {
			log.Printf("setting host path %s to container path %s\n", local, container)
		}

		dockerArgs = append(dockerArgs, "--mount")
		dockerArgs = append(dockerArgs, fmt.Sprintf("type=bind,source=%s,target=%s", local, container))
	}

	// Add args with naming tied to Taskfile value
	dockerArgs = append(dockerArgs, "nanos:latest")

	// Add in call string to be handled by run.sh to invoke ops
	dockerArgs = append(dockerArgs, strings.Join(args.Call, " "))

	// Define the call to be made
	cmd := exec.Command("docker", dockerArgs...)
	// Allow stdout and stderr to be sent to shell
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if args.Verbose {
		log.Printf("running %s\n\n", cmd.String())
	}

	// Run the command and report any errors
	err = cmd.Run()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}
