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

package tsrl

import (
	"sync"
	"time"

	"github.com/kricen/tsrl/cmap"
	"github.com/kricen/tsrl/model"
)

//BucketPool : add , remove ,ReleaseToken ,BorrowToken
type BucketPool struct {
	bmap cmap.ConcurrentMap
	rw   sync.RWMutex
}

//New :
func New() (pool *BucketPool) {
	bmap := cmap.New()
	pool = &BucketPool{bmap: bmap}
	return
}

// GetBucket :
func (b *BucketPool) GetBucket(url string) (bk *model.Bucket) {
	tmp, ok := b.bmap.Get(url)
	if !ok {
		bk = b.AddBucket(url, 0, 0, model.BUCKET_TYPE_TRAFFIC_SHAPING)
		return
	}
	bk, _ = tmp.(*model.Bucket)
	return
}

// AddBucket : Add Traffic Shaping,default bucket type is traffic shaping
func (b *BucketPool) AddBucket(url string, maxSize int64, timeoutDuration time.Duration, bucketType string) (bk *model.Bucket) {
	b.rw.Lock()
	defer b.rw.Unlock()
	bk = model.New(maxSize, timeoutDuration, bucketType)
	b.bmap.Set(url, bk)
	return
}
