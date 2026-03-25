[![Go Reference][goreference_badge]][goreference_link]
[![Go Report Card][goreportcard_badge]][goreportcard_link]
 
# Overview

Package stack defines goroutine\-safe methods for manipulating a generic [stack](<https://en.wikipedia.org/wiki/Stack_(abstract_data_type)>)

# Migration and Installation

## Breaking Changes in v2
 
- import path is now suffixed by v2
- an additional import is required to choose a Stack implementation:
  - `github.com/bruceesmith/stack/v2/mutex` for the mutex implementation
  - `github.com/bruceesmith/stack/v2/rendezvous` for the rendezvous implementation
- Stack variables continue to have type `*stack.Stack[T]`
- creation of a new Stack requires either `mutex.New` or `rendezvous.New` instead of `stack.New` 
- `rendezvous.New` takes a `context.Context` argument whereas `mutex.New` does not

## Upgrading to v2

Version 2.0.0 is a major version bump to add an alternate stack implementation. While the API for every stack function
remains identical to v1, the signature of `rendezvous.New()` differs from that of `mutex.New()`, and you must update your
import paths to include the /v2 suffix, and add a second import path to choose the desired implementation.

1. Update Import Paths

    Update all source files to use the new module path:
    ```go

    // Old v1 import
    import "github.com/bruceesmith/stack/v2"

    // New v2 import (using the mutex implementation)
    import "github.com/bruceesmith/stack/v2"
    import "github.com/bruceesmith/stack/v2/mutex"


    // New v2 import (using the rendezvous implementation)
    import "github.com/bruceesmith/stack/v2"
    import "github.com/bruceesmith/stack/v2/rendezvous"

    ```

2. Update your go.mod

    Run the following command in your project root to fetch the new version:
    ```bash
    go get github.com/bruceesmith/stack/v2
    ```

3. Summary of Changes

    * Breaking Changes: Import paths and New() function signatures, however the API signatures for stack operations are 
    identical
    * New Features: Two alternate implementations - mutex-based and rendezvous-based
    * Requirement: Go 1.26
