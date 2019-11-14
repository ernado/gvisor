// Copyright 2019 The gVisor Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stack

import "math/rand"

// RNG is a random number generator. It is not expected to be cryptographically
// secure, and should not be used for such applications.
type RNG interface {
	// Int63n returns, as an int64, a non-negative pseudo-random number in
	// [0,n). It may panic if n <= 0.
	Int63n(int64) int64
}

// defaultRNG is an implementation of RNG that uses the global random number
// generator from the math/rand package.
type defaultRNG struct{}

// Int63n implements RNG.Int63n.
//
// It is equivalent to calling Int63n(n) from the math/rand package.
func (*defaultRNG) Int63n(n int64) int64 {
	return rand.Int63n(n)
}
