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

package model

import (
	"errors"
	"sync/atomic"
	"time"
)

const (
	defaultMaxSize              = 500
	defaultTimeout              = 3 * time.Second
	BUCKET_TYPE_TRAFFIC_SHAPING = "traffic_shaping"
	BUCKET_TYPE_RATE_LIMIT      = "rate_limit"
)

var (
	ErrTimeout = errors.New("access token timeout")
)

// Bucket :
type Bucket struct {
	maxSize           int64
	qps               int64
	vacantSize        int64
	undistributedSize int64
	RWChan            chan interface{}
	timeoutDuration   time.Duration
	BucketType        string
	consume           int64
	produceToken      int64
}

//New : init a traffic-shaping Bucket with two param : maxSize  and timeoutDuration
//      maxSize : max size that this bucket can hold
//      timeoutDuration : max time that borrow a token from bucket
func New(maxSizeOrQPS int64, timeoutDuration time.Duration, bucketType string) (bk *Bucket) {
	if maxSizeOrQPS <= 0 {
		maxSizeOrQPS = defaultMaxSize
	}
	if timeoutDuration <= 0 {
		timeoutDuration = defaultTimeout
	}

	if bucketType == BUCKET_TYPE_TRAFFIC_SHAPING {
		channel := make(chan interface{}, maxSizeOrQPS)
		bk = &Bucket{maxSize: maxSizeOrQPS, undistributedSize: maxSizeOrQPS,
			RWChan: channel, vacantSize: maxSizeOrQPS, timeoutDuration: timeoutDuration, BucketType: BUCKET_TYPE_TRAFFIC_SHAPING}
	} else {

		channel := make(chan interface{}, maxSizeOrQPS)
		bk = &Bucket{qps: maxSizeOrQPS, RWChan: channel, vacantSize: 0, timeoutDuration: timeoutDuration, maxSize: maxSizeOrQPS, BucketType: BUCKET_TYPE_RATE_LIMIT}

		go bk.produce()
	}

	return bk
}

// BorrowToken : access token from bucket.
// Specification : if timeoutDruation is little than zero , program will set timeoutDruation equals Buckte's
//                  timeoutDruation set in init , while  if not set when function init, timeoutDruation is
//                  a default value : 3 seconds
func (b *Bucket) BorrowToken(timeoutDruation time.Duration) (token interface{}, err error) {

	if timeoutDruation <= 0 {
		timeoutDruation = b.timeoutDuration
	}

	// define a received-only timeout channel
	// check undistributedSize is great than zero ,if gt 0 , produce a token
	if b.BucketType == BUCKET_TYPE_TRAFFIC_SHAPING {
		if uds := atomic.LoadInt64(&b.undistributedSize); uds > 0 {
			val := atomic.AddInt64(&b.undistributedSize, -1)
			if val < 0 {
				atomic.AddInt64(&b.undistributedSize, 1)
			} else {
				go func() {
					atomic.AddInt64(&b.produceToken, 1)
					b.RWChan <- uds
				}()
			}

		}
	}
	timeout := time.After(timeoutDruation)
	select {
	case <-timeout:
		return nil, ErrTimeout
	case token = <-b.RWChan:
		atomic.AddInt64(&b.vacantSize, -1)
		atomic.AddInt64(&b.consume, 1)
	}

	return
}

// ReleaseToken : return token borrow from bucket
func (b *Bucket) ReleaseToken(token interface{}) {
	if b.BucketType == BUCKET_TYPE_TRAFFIC_SHAPING {

		atomic.AddInt64(&b.vacantSize, 1)
		go func() {
			b.RWChan <- token
		}()
	}

}

// produce : a function to produce the token ,and throw token into bucket
func (b *Bucket) produce() {

	// calcute the produce speed
	speed := int(float64(1) / float64(b.qps) * 1 * 1000 * 1000)
	produceSpeed := time.Duration(speed) * time.Microsecond

	tick := time.Tick(produceSpeed)

	for {
		select {
		case <-tick:
			if vacantSize := atomic.LoadInt64(&b.vacantSize); vacantSize < b.maxSize {
				atomic.AddInt64(&b.vacantSize, 1)
				go func() {
					b.RWChan <- 1
					atomic.AddInt64(&b.produceToken, 1)
				}()
			}
		}

	}
}
