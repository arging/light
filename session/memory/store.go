// Copyright 2014 li. All rights reserved.

package memory

import (
	"github.com/roverli/light/session"
	"github.com/roverli/utils/errors"
	"sync"
	"time"
)

var _ session.Store = &MemoryStore{}

// MemoryStore stores the session in the memory.
type MemoryStore struct {
	mutex    sync.RWMutex              // lock for sessions
	sessions map[string]*MemorySession // session storage
	config   session.Config
}

func (s *MemoryStore) Get(id string) (session.Session, errors.Error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	session := s.sessions[id]
	if session != nil {
		session.incAccess()
		return session, nil
	}
	return nil, nil
}

func (s *MemoryStore) New(id string) (session.Session, errors.Error) {

	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, dup := s.sessions[id]; dup {
		return nil, errors.New("session/memory: duplicate session id[" + id + "].")
	}

	current := time.Now().UnixNano()
	session := &MemorySession{
		id:               id,
		creationTime:     current,
		lastAccessedTime: current,
		attributes:       make(map[string]interface{}),
		hits:             1}
	s.sessions[id] = session

	return session, nil
}

func (s *MemoryStore) Gc() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now().UnixNano()
	interv := int64(time.Second / time.Nanosecond * time.Duration(s.config.MaxActiveInterval))
	for id, session := range s.sessions {
		t := session.lastAccessedTime
		if t+interv < now {
			delete(s.sessions, id)
		}
	}
}

func (s *MemoryStore) Save(session session.Session) errors.Error {
	return nil
}

func (s *MemoryStore) Init(config session.Config) errors.Error {
	s.sessions = make(map[string]*MemorySession)
	s.config = config
	return nil
}

func (s *MemoryStore) Close() {
	s.sessions = nil
}
