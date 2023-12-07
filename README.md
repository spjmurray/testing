# Go Testing on Steriods

Ever had a bunch of tests that battled over finite resources?
Ever been plagued by non-deterministic errors due to race conditions?
Fear not!

This testing module is a simple wrapper around Go's native testing framework, but the difference is that it's aware of resources.
All you need to do is:

* Poll your infrastructure for physical resources or quotas
* Register these in `TestMain()`
* Call the custom `Parallel()` function in each test, with the required number of resources to run the test
* Profit like a Ferengi

## Documentation

Have a look at the test code.
