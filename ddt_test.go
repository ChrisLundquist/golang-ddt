package ddt

import (
	"sync"
	"testing"
	"time"
)

func TestDDT100ms(t *testing.T) {
	tracker, _ := New(1)
	scale := 100 * time.Millisecond

	for i := 0; i < 10; i++ {
		tracker.Tick("Foo")
		time.Sleep(scale)
	}

	info, _ := tracker.Get("Foo")

	if info.Average > scale*120/100 || info.Average < scale*90/100 {
		t.Fatalf("bad average: %v. Expected: %v", info.Average, scale)
	}
}

func TestDDT10ms(t *testing.T) {
	tracker, _ := New(1)
	scale := 10 * time.Millisecond

	for i := 0; i < 100; i++ {
		tracker.Tick("Foo")
		time.Sleep(scale)
	}

	info, _ := tracker.Get("Foo")

	if info.Average > scale*120/100 || info.Average < scale*90/100 {
		t.Fatalf("bad average: %v. Expected: %v", info.Average, scale)
	}
}

func TestDDT1ms(t *testing.T) {
	tracker, _ := New(1)
	scale := 1 * time.Millisecond

	for i := 0; i < 1000; i++ {
		tracker.Tick("Foo")
		time.Sleep(scale)
	}

	info, _ := tracker.Get("Foo")

	if info.Average > scale*120/100 || info.Average < scale*90/100 {
		t.Fatalf("bad average: %v. Expected: %v", info.Average, scale)
	}
}

func TestConcurrentTick(t *testing.T) {
	var wg sync.WaitGroup
	tracker, _ := New(1)
	scale := 100 * time.Millisecond

	for i := 0; i < 10; i++ {
		go func() {
			wg.Add(1)
			defer wg.Done()
			for i := 0; i < 10; i++ {
				tracker.Tick("Foo")
				time.Sleep(scale)
			}
		}()
	}

	wg.Wait()

	info, _ := tracker.Get("Foo")

	if info.Average > scale*110/100 || info.Average < scale*90/100 {
		t.Fatalf("bad average: %v. Expected: %v", info.Average, scale)
	}
}
