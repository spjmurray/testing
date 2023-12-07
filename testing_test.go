/*
Copyright 2023-2024 Simon Murray.

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

package testing_test

import (
	"os"
	"testing"
	"time"

	smtest "github.com/spjmurray/testing"
)

const (
	ResourceCPU = "cpu"
	ResourceRAM = "memory"
)

func TestMain(m *testing.M) {
	resources := smtest.ResourceSet{
		ResourceCPU: 16,
		ResourceRAM: 64,
	}

	smtest.Start(resources)

	os.Exit(m.Run())
}

func TestSuccess1(t *testing.T) {
	resources := smtest.ResourceSet{
		ResourceCPU: 8,
		ResourceRAM: 32,
	}

	defer smtest.Parallel(t, resources)()

	time.Sleep(time.Second)
}

func TestSuccess2(t *testing.T) {
	resources := smtest.ResourceSet{
		ResourceCPU: 8,
		ResourceRAM: 32,
	}

	defer smtest.Parallel(t, resources)()

	time.Sleep(time.Second)
}

func TestSuccess3(t *testing.T) {
	resources := smtest.ResourceSet{
		ResourceCPU: 8,
		ResourceRAM: 32,
	}

	defer smtest.Parallel(t, resources)()

	time.Sleep(time.Second)
}

func TestSuccess4(t *testing.T) {
	resources := smtest.ResourceSet{
		ResourceCPU: 16,
		ResourceRAM: 64,
	}

	defer smtest.Parallel(t, resources)()

	time.Sleep(time.Second)
}

func TestSkip1(t *testing.T) {
	resources := smtest.ResourceSet{
		ResourceCPU: 32,
	}

	defer smtest.Parallel(t, resources)()
}

func TestSkip2(t *testing.T) {
	resources := smtest.ResourceSet{
		"gpu": 2,
	}

	defer smtest.Parallel(t, resources)()
}
