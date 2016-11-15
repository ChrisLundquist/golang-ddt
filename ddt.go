package ddt

import (
	//	"fmt"
	"github.com/hashicorp/golang-lru"
	"time"
)

type DDT struct {
	lru *lru.Cache
}

type TimerInfo struct {
	Last    time.Time
	Average time.Duration
}

// The number of values you'd like to track
func New(size int) (*DDT, error) {
	cache, err := lru.New(size)
	if err != nil {
		return nil, err
	}

	return &DDT{lru: cache}, nil
}

func (d *DDT) Tick(key interface{}) {
	value, found := d.lru.Peek(key)

	now := time.Now()
	var info TimerInfo
	if found == false { /* New Key */
		info = TimerInfo{Last: now, Average: 0}
	} else { /* Existing Key */
		info = *value.(*TimerInfo)
		last := info.Last
		average := info.Average

		duration := now.Sub(last)
		if average == 0 { // this is the second sample
			average = duration
		} else {
			average = (average + duration) / 2
		}
		//fmt.Printf("avg: %v\n", average)
		info = TimerInfo{Last: now, Average: average}
	}
	d.lru.Add(key, &info)
}

func (d *DDT) Get(key interface{}) (*TimerInfo, bool) {
	info, found := d.lru.Get(key)

	return info.(*TimerInfo), found
}
