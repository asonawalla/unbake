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
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestBasicLoad(t *testing.T) {
	var testFile, err = createTestConfigFile()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Remove(testFile); err != nil {
			t.Fatalf("removing test file: %s", err)
		}
	}()

	result, err := unbake(testFile)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatal("wrong number of command results")
	}

	var expect = map[string]struct{}{
		"docker build -t foo -f docker/go_container.dockerfile --build-arg cmd=foo .": {},
		"docker build -t bar -f docker/go_container.dockerfile --build-arg cmd=bar .": {},
	}

	var actual = make(map[string]struct{})
	for _, val := range result {
		actual[val] = struct{}{}
	}

	if !reflect.DeepEqual(expect, actual) {
		t.Fatalf("command mismatch, expected:\n %v\n got:\n %v", expect, actual)
	}
}

func TestLoadWithBuildKit(t *testing.T) {
	buildKit = true
	var testFile, err = createTestConfigFile()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Remove(testFile); err != nil {
			t.Fatalf("removing test file: %s", err)
		}
	}()

	result, err := unbake(testFile)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatal("wrong number of command results")
	}

	var expect = map[string]struct{}{
		"DOCKER_BUILDKIT=1 docker build -t foo -f docker/go_container.dockerfile --build-arg cmd=foo .": {},
		"DOCKER_BUILDKIT=1 docker build -t bar -f docker/go_container.dockerfile --build-arg cmd=bar .": {},
	}

	var actual = make(map[string]struct{})
	for _, val := range result {
		actual[val] = struct{}{}
	}

	if !reflect.DeepEqual(expect, actual) {
		t.Fatalf("command mismatch, expected:\n %v\n got:\n %v", expect, actual)
	}
}

func TestLoadWithCustomConfig(t *testing.T) {
	buildKit = true
	dockerCfg = "/path/to/config"
	var testFile, err = createTestConfigFile()
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Remove(testFile); err != nil {
			t.Fatalf("removing test file: %s", err)
		}
	}()

	result, err := unbake(testFile)
	if err != nil {
		t.Fatal(err)
	}

	if len(result) != 2 {
		t.Fatal("wrong number of command results")
	}

	var expect = map[string]struct{}{
		"DOCKER_BUILDKIT=1 docker --config=/path/to/config build -t foo -f docker/go_container.dockerfile --build-arg cmd=foo .": {},
		"DOCKER_BUILDKIT=1 docker --config=/path/to/config build -t bar -f docker/go_container.dockerfile --build-arg cmd=bar .": {},
	}

	var actual = make(map[string]struct{})
	for _, val := range result {
		actual[val] = struct{}{}
	}

	if !reflect.DeepEqual(expect, actual) {
		t.Fatalf("command mismatch, expected:\n %v\n got:\n %v", expect, actual)
	}
}

var testConfig = []byte(`
group "default" {
  targets = ["foo", "bar"]
}

target "go_container" {
  dockerfile = "docker/go_container.dockerfile"
}

target "foo" {
  inherits = ["go_container"]
  tags = ["foo"]
  args = {
    cmd = "foo"
  }
}

target "bar" {
  inherits = ["go_container"]
  tags = ["bar"]
  args = {
    cmd = "bar"
  }
}
`)

func createTestConfigFile() (string, error) {
	tmpfile, err := ioutil.TempFile("", "unbake_test")
	if err != nil {
		return "", err
	}
	if _, err := tmpfile.Write(testConfig); err != nil {
		return "", err
	}
	if err := tmpfile.Close(); err != nil {
		return "", err
	}
	return tmpfile.Name(), nil
}
