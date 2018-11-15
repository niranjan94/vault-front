package multilock

import (
	"runtime"
	"testing"
	"time"
)

func HammerMutex(m *MultiLock, id interface{}, loops int, cdone chan bool) {
	for i := 0; i < loops; i++ {
		m.Lock(id)
		time.Sleep(1 * time.Millisecond)
		m.Unlock(id)
	}
	cdone <- true
}

func HammerMutexReadLock(m *MultiLock, id interface{}, loops int, cdone chan bool) {
	for i := 0; i < loops; i++ {
		m.RLock(id)
		time.Sleep(1 * time.Millisecond)
		m.RUnlock(id)
	}
	cdone <- true
}

func TestMultiLock(t *testing.T) {
	multiLock := NewMultiLock()
	if n := runtime.SetMutexProfileFraction(1); n != 0 {
		t.Logf("got mutexrate %d expected 0", n)
	}
	defer runtime.SetMutexProfileFraction(0)
	c := make(chan bool)
	for i := 0; i < 10; i++ {
		go HammerMutex(multiLock, "this-id", 500, c)
	}
	for i := 0; i < 10; i++ {
		go HammerMutex(multiLock, "this-another-id", 500, c)
	}
	for i := 0; i < 20; i++ {
		<-c
	}
}

func TestReadMultiLock(t *testing.T) {
	GetMultiLock()
	multiLock := NewMultiLock()
	if n := runtime.SetMutexProfileFraction(1); n != 0 {
		t.Logf("got mutexrate %d expected 0", n)
	}
	defer runtime.SetMutexProfileFraction(0)
	c := make(chan bool)
	for i := 0; i < 10; i++ {
		go HammerMutexReadLock(multiLock, "this-id", 500, c)
	}
	for i := 0; i < 10; i++ {
		go HammerMutexReadLock(multiLock, "this-another-id", 500, c)
	}
	for i := 0; i < 20; i++ {
		<-c
	}
}