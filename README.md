# unbake
[![CircleCI](https://circleci.com/gh/asonawalla/unbake.svg?style=svg)](https://circleci.com/gh/asonawalla/unbake)
[![GoDoc](https://godoc.org/github.com/asonawalla/unbake?status.svg)](https://godoc.org/github.com/asonawalla/unbake)
[![Go Report Card](https://goreportcard.com/badge/github.com/asonawalla/unbake)](https://goreportcard.com/report/github.com/asonawalla/unbake)

**Archive Notice** This repository is no longer maintained. If you want to use it for your own purposes, please fork it and update all upstream dependencies first.

Reverse docker buildx bake files into sequential docker commands

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

## Usage

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

## Compatibility
Though upstream docker/buildx's high level build constructs are a WIP and buildx itself is listed as a "tech preview",
this repository uses a forked and pinned version of buildx at asonawalla/buildx. Unbake will update with releases
as the fork tracks upstream. That means that if you use asonawalla/buildx, it will reliably be compatible with
asonawalla/unbake.
