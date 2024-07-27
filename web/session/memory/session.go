package memory

import (
	"awesomeProject/web/session"
	"context"
	"errors"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

var (
	errorKeyNotFound = errors.New("session: key not found")
)

type Store struct {
	mu         sync.RWMutex
	expiration time.Duration
	sessions   *cache.Cache
}

func NewStore(expiration time.Duration) *Store {
	return &Store{
		sessions: cache.New(expiration, time.Second),
	}
}

func (s *Store) Generate(ctx context.Context, id string) (session.Session, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess := &Session{
		id:     id,
		values: sync.Map{},
	}

	s.sessions.Set(id, sess, s.expiration)
	return sess, nil
}

func (s *Store) Refresh(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.sessions.Get(id)
	if !ok {
		return errorKeyNotFound
	}
	s.sessions.Set(id, sess, s.expiration)

	return nil
}

func (s *Store) Remove(ctx context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions.Delete(id)
	return nil
}

func (s *Store) Get(ctx context.Context, id string) (session.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sess, ok := s.sessions.Get(id)
	if !ok {
		return nil, errorKeyNotFound
	}

	return sess.(session.Session), nil
}

type Session struct {
	id     string
	values sync.Map
}

func (s *Session) Get(ctx context.Context, key string) (any, error) {
	val, ok := s.values.Load(key)
	if !ok {
		return "", errorKeyNotFound
	}

	return val.(string), nil
}

func (s *Session) Set(ctx context.Context, key string, value any) error {
	s.values.Store(key, value)
	return nil
}

func (s *Session) ID() string {
	return s.id
}
