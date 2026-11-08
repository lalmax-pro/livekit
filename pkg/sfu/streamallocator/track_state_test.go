// Copyright 2026 LiveKit, Inc.
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

package streamallocator

import (
	"sync"
	"testing"
)

func TestStreamStateConcurrentSubscribeMutable(t *testing.T) {
	track := &Track{}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			track.SetStreamState(StreamStatePaused)
			track.SetStreamState(StreamStateActive)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			_ = track.IsSubscribeMutable()
		}
	}()
	wg.Wait()
	if !track.IsSubscribeMutable() || track.SetStreamState(StreamStateActive) {
		t.Fatal("active stream must remain mutable and report no state change")
	}
}
