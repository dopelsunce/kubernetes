/*
Copyright 2016 The Kubernetes Authors All rights reserved.

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

package framework

import (
	"testing"
	"time"

)

// TestPopReleaseLock tests that when processor listener blocks on chan,
// it should release the lock for pendingNotifications.
func TestPopReleaseLock(t *testing.T) {
	pl := newProcessListener(nil)
	stopCh := make(chan struct{})
	defer close(stopCh)
	// make pop() block on nextCh: waiting for receiver to get notification.
	pl.add(1)
	go pl.pop(stopCh)

	resultCh := make(chan struct{})
	go func() {
		pl.lock.Lock()
		close(resultCh)
	}()

	select {
	case <-resultCh:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Timeout after %v", 100 * time.Millisecond)
	}
	pl.lock.Unlock()
}

func TestDDNeo(t *testing.T) {
	pl := newProcessListener(nil)
	stopCh := make(chan struct{})
	// make pop() block on nextCh: waiting for receiver to get notification.
	pl.add(1)

	resultCh := make(chan struct{})
	go func() {
		pl.pop(stopCh)
		resultCh <- struct{}{}
	}()
	pl.lock.Lock()
	pl.lock.Unlock()
	close(stopCh)

	select {
	case <-resultCh:
	case <-time.After(100 * time.Millisecond):
		t.Errorf("Timeout! DD triggered")
	}
}