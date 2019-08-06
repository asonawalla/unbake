/*
Copyright © 2019 Azim Sonawalla <azim.sonawalla@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"context"
	"fmt"

	"github.com/docker/buildx/bake"
)

func unbake(file string) ([]string, error) {
	var targets, err = bake.ReadTargets(context.Background(), []string{file}, []string{"default"}, []string{})
	if err != nil {
		return nil, err
	}
	return targetsToCommands(targets), nil
}

func targetsToCommands(targets map[string]bake.Target) []string {
	var result []string
	for _, t := range targets {
		for _, tag := range t.Tags {
			var cmd = fmt.Sprintf("DOCKER_BUILDKIT=1 docker build -t %s -f %s ", tag, *t.Dockerfile)
			for k, v := range t.Args {
				cmd += fmt.Sprintf("--build-arg %s=%s ", k, v)
			}
			cmd += "."
			result = append(result, cmd)
		}
	}
	return result
}
