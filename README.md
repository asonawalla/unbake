# unbake
Reverse docker buildx bake files into sequential docker commands

[![CircleCI](https://circleci.com/gh/asonawalla/unbake.svg?style=svg)](https://circleci.com/gh/asonawalla/unbake)

docker/buildx is a plugin to the most recent docker distribution
that allows users to define high level build information (bake files)
and transmit that build info to the buildkit-powered docker daemon
in a single invocation.

Since this is a relatively new piece of technolgoy, many environments (including
primarily CI and CD systems) don't yet have the most recent version of docker or
the buildx plugin, and it introduces significant overhead to install these.

unbake solves for this discrepancy by taking a bake file as an input and generating
plain docker commands to build those targets. You'd want to run the unbake container
in your CI against a bake file and pipe the output of it to shell to continue with a
multi-invocation docker build process.

## usage

### Installing using the go tool
Clone this repository and run `go install` in the module:
```
git clone https://github.com/asonawalla/unbake.git
cd unbake && go install
```

Then invoke with a bake file and pipe the output to `sh`:
```
unbake -f bake.hcl | sh
```

### Using with Docker
Alternately, you can run unbake using the docker container uploaded to docker hub (docker is the only
requirement, making this a promising option for CI environments):
```
docker run --rm -v$(pwd):/bake asonawalla/unbake /bin/unbake -f /bake/bake.hcl | sh
```
