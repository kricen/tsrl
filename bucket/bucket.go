// Copyright 2018 oliver kricen
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package bucket

import (
	"errors"
	"sync/atomic"
	"time"
)

const (
	defaultMaxSize = 10000
	defaultTimeout = 3 * time.Second
)

var (
	ErrTimeout = errors.New("access token timeout")
)

// Bucket :
type Bucket struct {
	maxSize           int64
	vacantSize        int64
	undistributedSize int64
	RWChan            chan interface{}
	timeoutDuration   time.Duration
}

//New : init a traffic-shaping Bucket with two param : maxSize  and timeoutDuration
//      maxSize : max size that this bucket can hold
//      timeoutDuration : max time that borrow a token from bucket
func New(maxSize int64, timeoutDuration time.Duration) *Bucket {
	if maxSize <= 0 {
		maxSize = defaultMaxSize
	}
	if timeoutDuration <= 0 {
		timeoutDuration = defaultTimeout
	}
	channel := make(chan interface{}, maxSize)
	bucket := &Bucket{maxSize: maxSize, undistributedSize: maxSize,
		RWChan: channel, vacantSize: maxSize, timeoutDuration: timeoutDuration}
	return bucket
}

// BorrowToken : access token from bucket.
// Specification : if timeoutDruation is little than zero , program will set timeoutDruation equals Buckte's
//                  timeoutDruation set in init , while  if not set when function init, timeoutDruation is
//                  a default value : 3 seconds
func (b *Bucket) BorrowToken(timeoutDruation time.Duration) (token interface{}, err error) {
	// check undistributedSize is great than zero ,if gt 0 , produce a token

	if uds := atomic.LoadInt64(&b.undistributedSize); uds > 0 {
		val := atomic.AddInt64(&b.undistributedSize, -1)
		if val < 0 {
			atomic.AddInt64(&b.undistributedSize, 1)
		} else {
			go func() {
				b.RWChan <- uds
			}()
		}

	}
	if timeoutDruation <= 0 {
		timeoutDruation = b.timeoutDuration
	}
	// define a received-only timeout channel
	timeout := time.After(timeoutDruation)

	select {
	case <-timeout:
		return nil, ErrTimeout
	case token = <-b.RWChan:
		atomic.AddInt64(&b.vacantSize, -1)
	}
	return
}

// ReleaseToken : return token borrow from bucket
func (b *Bucket) ReleaseToken(token interface{}) {
	atomic.AddInt64(&b.vacantSize, 1)
	go func() {
		b.RWChan <- token
	}()
}
