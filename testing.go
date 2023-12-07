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

package testing

import (
	"fmt"
	"testing"
	"time"
)

// queueItem constains all the bits to hold a test up until enough
// resources are free.
type queueItem struct {
	// wait is closed to release the test.
	wait chan interface{}

	// required is the set of resources that are required for the
	// test to successfully execute.
	required ResourceSet
}

// transaction is used to enqueue an item.
type transaction struct {
	// name is the test name e.g. unique.
	name string

	// item is the item to add to the queue.
	item *queueItem
}

var (
	// available are the global set of resources that are available.
	// Sadly the standard testing package doesn't allow a context etc.
	// to be passed from TestMain to individual tests, so we're stack
	// with "bad practice".  We can use this value to check if a test
	// can actually be run.
	available ResourceSet

	// unallocated is the set of resources that are not in use.
	unallocated = ResourceSet{}

	// queue is the set of tests waiting to run.
	queue = map[string]*queueItem{}

	// enqueue adds a test to our scheduler.
	enqueue chan *transaction

	// release is called on test exit to release resources.
	release chan ResourceSet
)

// Start is called from TestMain to set things up for example:
//
//	import (
//	  "testing"
//	  "os"
//
//	  smtest "github.com/spjmurray/testing"
//	)
//
//	const (
//	  ResourceTypeCPU    = "cpu"
//	  ResourceTypeMemory = "memory"
//	)
//
//	func TestMain(m *testing.M) {
//	   // Do something here like interrogate the infrastructure for
//	   // resource or quota limits...
//	   resources := smtest.RedourceSet{
//	     ResourceTypeCPU:    16,
//	     ResourceTypeMemory: 64,
//	   }
//
//	   smtest.Start(resources)
//
//	   os.Exit(m.Run())
//	}
func Start(resources ResourceSet) {
	available = resources

	for k, v := range available {
		unallocated[k] = v
	}

	enqueue = make(chan *transaction)
	release = make(chan ResourceSet)

	go func() {
		for {
			// Process new tests, and finishing tests in a concurrency
			// safe way.  New tests go on the queue, finished tests will
			// release their resource allocations.
			select {
			case transaction := <-enqueue:
				queue[transaction.name] = transaction.item
			case allocated := <-release:
				for k, v := range allocated {
					unallocated[k] += v
				}
			}

			// For every item on the queue...
			for name, item := range queue {
				ok := true

				// If all of its required resources can be satisfied...
				for k, v := range item.required {
					if unallocated[k] < v {
						ok = false
						break
					}
				}

				if ok {
					// Remove them from the unallocated pool, remove the
					// enqueued item and release the test.
					for k, v := range item.required {
						unallocated[k] -= v
					}

					delete(queue, name)
					close(item.wait)
				}
			}
		}
	}()
}

// Parallel is called from individual tests, it delegates concurrency to the native
// testing library, but crucially only releases a test for execution once resource
// is available.  If a test requires too many resources, or none are available at all
// then the test is skipped.
func Parallel(t *testing.T, required ResourceSet) func() {
	for k, v := range required {
		availableResource, ok := available[k]
		if !ok || v > availableResource {
			t.Skipf("test requires %d %s, %d available", v, k, availableResource)
		}
	}

	// This call pops the test onto the queue, and will respect go's standard
	// concurrency guarantees...
	t.Parallel()

	wait := make(chan interface{})

	// Enqueue the test with the scheduler...
	transaction := &transaction{
		name: t.Name(),
		item: &queueItem{
			wait:     wait,
			required: required,
		},
	}

	enqueue <- transaction

	fmt.Printf("+++ ALLOC %s\n", t.Name())

	// Wait for resource to become available...
	<-wait

	fmt.Printf("+++ SCHED %s\n", t.Name())

	start := time.Now()

	return func() {
		fmt.Printf("+++ END   %s (%.2fs)\n", t.Name(), time.Since(start).Seconds())

		release <- required
	}
}
