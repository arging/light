// Copyright 2014 li. All rights reserved.

package session

import (
	"github.com/roverli/utils/errors"
	"net/http"
	"net/url"
	"sync/atomic"
	"time"
)

var (
	ErrManagerClosed = errors.New("light/session: Session manager is already closed.")
	ErrGenSessionId  = errors.New("light/session: Gen session id error.")
)

// Provides a way to identify a user across more than one page request or visit
// to a web site and to store information about that user. The web container
// uses this interface to create a session between an HTTP client and an HTTP server.
// The session persists for a specified time period, across more than one connection or
// page request from the user. A session usually corresponds to one user, who may
// visit a site many times. The server can maintain a session in many ways such as
// using cookies or rewriting URLs. This interface allows controllers to view and manipulate
// information about a session, such as the session identifier, creation time, and last accessed time.
// Bind objects to sessions, allowing user information to persist across multiple user connections.
// For session that are invalidated or expire, notifications are sent after
// the session has been invalidated or expired. A controllr should be able to handle cases in which
// the client does not choose to join a session, such as when cookies are
// intentionally turned off. Until the client joins the session, isNew returns true.
// If the client chooses not to join the session, will return a different session
// on each request, and isNew will always return true.
type Session interface {

	// Returns a string containing the unique identifier assigned to this session.
	// The identifier is assigned by the session container and is implementation dependent.
	Id() string

	// Returns the time when this session was created, measured in nanoseconds.
	CreationTime() int64

	// Returns the last time the client sent a request associated with
	// this session, as the number of nanoseconds and marked by the time
	// the container received the request.
	// Actions that your application takes, such as getting or setting
	// a value associated with the session, do not affect the access time.
	LastAccessedTime() int64

	// Returns the object bound with the specified name in this session, or
	// nil if no object is bound under the name.
	GetAttribute(name string) interface{}

	// Returns an enumeration of string objects
	// containing the names of all the objects bound to this session.
	// If no values bound to the session, return empty slice.
	GetAttributeNames() []string

	// Binds an object to this session, using the name specified.
	// If an object of the same name is already bound to the session,
	// the object is replaced. If the value passed in is nil,
	// this has the same effect as calling removeAttribute().
	SetAttribute(name string, value interface{})

	// Removes the object bound with the specified name from
	// this session. If the session does not have an object
	// bound with the specified name, this method does nothing.
	RemoveAttribute(name string)

	// Invalidates this session then unbinds any objects bound to it.
	Invalidate()

	// Returns true if the client does not yet know about the
	// session or if the client chooses not to join the session.For
	// example, if the server used only cookie-based sessions, and
	// the client had disabled the use of cookies, then a session would
	// be new on each request.
	IsNew() bool
}

// Store is the interface for customer session stores.
type Store interface {

	// Return a cached session.
	Get(id string) (Session, errors.Error)

	// Create and return a new session.
	New(id string) (Session, errors.Error)

	// Persist session.
	Save(s Session) errors.Error

	// Initialize the session store.
	Init(config Config) errors.Error

	// Collect the expired sessions.
	Gc()
}

// Manager maintains session objects.
// Responsible for managing opening and closing of sessions.
type Manager struct {
	isClosed int32
	ticker   *time.Ticker
	store    Store
	config   Config
}

// Retrieve a session from context by http request.
func (m *Manager) Get(r *http.Request) (Session, errors.Error) {

	if atomic.LoadInt32(&m.isClosed) == 1 {
		return nil, ErrManagerClosed
	}

	cookie1, err1 := r.Cookie(m.config.CookieName + "1")
	cookie2, err2 := r.Cookie(m.config.CookieName + "2")

	switch {
	case err1 == http.ErrNoCookie && err2 == http.ErrNoCookie:
		return nil, nil

	case err1 != nil:
		return nil, errors.Wrapf(err1, "light/session: Get session cookie error.")

	case err2 != nil:
		return nil, errors.Wrapf(err2, "light/session: Get session cookie error.")

	default:
		return m.store.Get(cookie1.Value + m.config.Sed + cookie2.Value)
	}
}

// Create a plain Session for the current application context.
// Will usually be a new ClientSession for the current request.
func (m *Manager) Create(r *http.Request, w http.ResponseWriter) (Session, errors.Error) {
	if atomic.LoadInt32(&m.isClosed) == 1 {
		return nil, ErrManagerClosed
	}

	id1, err1 := Md5Id(r.RemoteAddr + m.config.Sed)
	id2, err2 := Md5Id(r.RemoteAddr + m.config.Sed)
	if err1 != nil || err2 != nil {
		return nil, ErrGenSessionId
	}

	cookie1 := m.newCookie(id1)
	cookie2 := m.newCookie(id2)

	if m.config.EnableCookie {
		http.SetCookie(w, cookie1)
		http.SetCookie(w, cookie2)
	}
	r.AddCookie(cookie1)
	r.AddCookie(cookie2)

	return m.store.New(id1 + m.config.Sed + id2)
}

func (m *Manager) newCookie(id string) *http.Cookie {
	return &http.Cookie{
		Name:     m.config.CookieName,
		Value:    url.QueryEscape(id),
		Path:     "/",
		HttpOnly: m.config.HttpOnly,
		Secure:   m.config.Secure,
		MaxAge:   m.config.MaxAge,
		Domain:   m.config.Domain}
}

// Persist session to the underlying store.
func (m *Manager) Save(s Session) errors.Error {
	if atomic.LoadInt32(&m.isClosed) == 1 {
		return ErrManagerClosed
	}
	return m.store.Save(s)
}

// Close this session manager, stop gc timer.
// shutting down all internal resources.
func (m *Manager) Close() {
	atomic.StoreInt32(&m.isClosed, 1)
	m.ticker.Stop()
}

// Start the session manager
func (m *Manager) start() errors.Error {
	m.ticker = time.NewTicker(time.Duration(m.config.GcInterval) * time.Second)
	go func() {
		for _ = range m.ticker.C {
			m.store.Gc()
		}
	}()
	return nil
}
