# Docker Ops

This project is an attempt to enable Nanos unikernels to be managed by Ops on
non-intel architectures such as the Mac M1 ARM64.

Unless there is something I have missed (as of 20 November, 2021) Ops does not
run properly on the M1 ARM64 architecture. This is because the implementation of
nanos used by Ops currently assumes an Intel64 environment. You can run Ops for
things that don't involve building, running, or deploying nanos images. Putting
Ops (and therefore Nanos) in an Intel64 container allows Nanos to operate in its
intended environment.

[Ops](https://ops.city) is a build and deployment tool for the
[Nanos](https://nanovms.com/) unikernel. A unikernel is a minimal operating
environment which is used to create a compatible image to run a single
applicaiton in cloud environments and for use locally using `qemu`. I use it to
test out unikernels running Go applications. Docker is similar in that it is a
scaled down Linux or other Operating system container intended to run at a
single entry point. A Docker container can have most things installed that would
normally be installed on any operating system, including additional user
accounts, logging daemons, etc. Nanos only runs one thing and its purpose is to
handle calls for an Intel64 linux architecture so that the single application will
be able to run. Nanos also provides useful things such as network port and
filesystem access.

The goal of this code is to provide a workflow that is as friction-free as
possible in terms of building and deploying nanos unikernels. Using the
dockerops application you can call Ops running in a Docker intel64 image as if
it were running on its own.

This project should work fine on a non-M1 mac but that would be redundant, as
Ops runs well on Intel64 macs.

## To run

1) install the [TaskFile runner](https://taskfile.dev/). A few conventions are
   used to ensure that the proper container is called and the tasks consistently
   manage that. I know how to write make files but I find Taskfile to be easy to
   read and use. On a mac with Homebrew you can install Taskfile with 
   ```
   brew install go-task/tap/go-task
   ```
2) Make sure that you have a running and recent Docker installation that
   supports multiple architectures (for the purposes of this application,
   Intel64). Any Mac release from 18 April, 2021 or later should.
3) Build the container using the task `task build` in the `build` directory.
4) Compile dockerops in `cmd/dockerops` using `go build .` .
5) Make sure you have a valid config file (see the sample in the `config`
   directory. This file needs either to be in the same directory as the binary
   or have its location indicated using the `-c` flag when invoking dockerops.
   1) Note that to do useful things you will need to expose a directory
      containing things like GCP authentication files. See the config file for
      this. 
6) Run dockerops. See the [Ops site](https://ops.city/) for information on how
   to run Ops and use it to make containers and deploy them to the cloud.

### Sample usage

Here is the usage output for the dockerops app

```
% ~/bin/dockerops -h
Usage: dockerops [--configpath CONFIGPATH] [--env ENV] [--verbose] [CALL [CALL ...]]

Positional arguments:
CALL                   call to ops - surround with quotes

Options:
--configpath CONFIGPATH, -c CONFIGPATH
                        config path - defaults to [dockeropps dir]/dockerops.yml
--env ENV, -e ENV      Set environment variable as key=val
--verbose, -v          print out what is being handled and done
--help, -h             display this help and exit
```

Here is a sample invocation

```
% ./dockerops
Usage:
ops [command]

Available Commands:
build       Build an image from ELF
deploy      Build an image from ELF and deploy an instance
env         Cross-build environment commands
help        Help about any command
image       manage nanos images
instance    manage nanos instances
pkg         Package related commands
profile     Profile
run         Run ELF binary as unikernel
update      check for updates
version     Version
volume      manage nanos volumes

Flags:
-h, --help            help for ops
    --show-debug      display debug messages
    --show-errors     display error messages
    --show-warnings   display warning messages

Use "ops [command] --help" for more information about a command.
```

Here is an invocation to list existing images

```
% ~/bin/dockerops ops image list
+---------------------+---------------------------------------+---------+--------------+
|        NAME         |                 PATH                  |  SIZE   |  CREATEDAT   |
+---------------------+---------------------------------------+---------+--------------+
| nanoapplinux.img    | /root/.ops/images/nanoapplinux.img    | 41.8 MB | 1 week ago   |
+---------------------+---------------------------------------+---------+--------------+
| nats-test-image.img | /root/.ops/images/nats-test-image.img | 44.8 MB | 2 months ago |
+---------------------+---------------------------------------+---------+--------------+
| natslinux.img       | /root/.ops/images/natslinux.img       | 44.7 MB | 2 months ago |
+---------------------+---------------------------------------+---------+--------------+
```

In the background the script run.sh is invoked in the container. This script
takes all passed in args and uses them to create an invocation of Ops, which is
in the container at `/app/ops`. If you put `/app/ops` or `ops` in your call it
will be cleaned up and the call made will be proper for Ops.

## Things to do
- Ensure that things like building and running Ops work
  - This so far has not been tested. Possible issues include stdout and stderr
    interaction when doing things like running an image in the container.
- Use this for development and make any improvements that arise from that
