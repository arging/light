// Copyright 2014 li. All rights reserved.

package memory

import (
	"github.com/roverli/light/session"
	"sync"
	"sync/atomic"
	"time"
)

var _ session.Session = &MemorySession{}

type MemorySession struct {
	id               string
	creationTime     int64
	lastAccessedTime int64
	attributes       map[string]interface{}
	lock             sync.RWMutex // lock for attributes
	hits             uint64
	isClosed         bool
}

func (s *MemorySession) Id() string {
	return s.id
}

func (s *MemorySession) CreationTime() int64 {
	return s.creationTime
}

func (s *MemorySession) LastAccessedTime() int64 {
	return atomic.LoadInt64(&s.lastAccessedTime)
}

func (s *MemorySession) GetAttribute(name string) interface{} {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.attributes[name]
}

func (s *MemorySession) GetAttributeNames() []string {
	s.lock.RLock()
	defer s.lock.RUnlock()

	i := 0
	names := make([]string, len(s.attributes))
	for k, _ := range s.attributes {
		names[i] = k
		i++
	}
	return names
}

func (s *MemorySession) SetAttribute(name string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.attributes[name] = value
}

func (s *MemorySession) RemoveAttribute(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.attributes, name)
}

func (s *MemorySession) Invalidate() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.attributes = make(map[string]interface{})
	s.isClosed = true
}

func (s *MemorySession) IsNew() bool {
	return atomic.LoadUint64(&s.hits) == 1
}

func (s *MemorySession) incAccess() {
	atomic.AddUint64(&s.hits, 1)
	atomic.StoreInt64(&s.lastAccessedTime, time.Now().UnixNano())
}
