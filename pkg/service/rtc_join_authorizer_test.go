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

package service

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/livekit/protocol/auth"
	"github.com/livekit/protocol/livekit"
)

func TestRTCJoinAuthorizer(t *testing.T) {
	svc, manager := &RTCService{}, &RoomManager{}
	denied := errors.New("denied")
	check := func(_ context.Context, _ *auth.ClaimGrants, roomID livekit.RoomID) error {
		if roomID == "old" {
			return denied
		}
		return nil
	}
	svc.SetJoinAuthorizer(check)
	manager.SetRTCJoinAuthorizer(check)
	if err := svc.authorizeJoin(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	if err := manager.authorizeRTCJoin(context.Background(), nil, "old"); !errors.Is(err, denied) {
		t.Fatal("room-generation rejection was lost")
	}
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			svc.SetJoinAuthorizer(nil)
			manager.SetRTCJoinAuthorizer(nil)
			svc.SetJoinAuthorizer(check)
			manager.SetRTCJoinAuthorizer(check)
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < 1000; i++ {
			_ = svc.authorizeJoin(context.Background(), nil)
			_ = manager.authorizeRTCJoin(context.Background(), nil, "old")
		}
	}()
	wg.Wait()
}
