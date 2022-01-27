// package drain collects and stores messages sent over a channel
package drain

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Store struct {
	LastRead Index
	Msg      Message
}

type Message struct {
	sync.RWMutex
	Array []interface{}
}

type Index struct {
	sync.RWMutex
	Idx int
}

func New(ch chan interface{}, ctx context.Context) *Store {

	s := Store{}
	s.Flush()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-ch:
				s.Msg.Lock()
				s.Msg.Array = append(s.Msg.Array, msg)
				s.Msg.Unlock()
			}
		}

	}()

	return &s
}

func (s *Store) Flush() {
	s.Msg.Lock()
	s.Msg.Array = make([]interface{}, 0)
	s.Msg.Unlock()
	s.LastRead.Lock()
	s.LastRead.Idx = -1
	s.LastRead.Unlock()
}

func (s *Store) NextNoWait() (interface{}, error) {

	if s.IsEmpty() {
		return nil, errors.New("empty store")
	}

	s.LastRead.Lock()
	s.Msg.RLock()

	var lastMsg interface{}

	err := errors.New("No new messages")

	if len(s.Msg.Array) > (s.LastRead.Idx + 1) {

		s.LastRead.Idx++
		lastMsg = s.Msg.Array[s.LastRead.Idx]
		err = nil

	}

	s.LastRead.Unlock()
	s.Msg.RUnlock()

	return lastMsg, err
}

// blocking version of get next message
func (s *Store) Next(ch chan interface{}, ctx context.Context) <-chan interface{} {

	c := make(chan interface{}, 1)

	go func() {

		for {

			msg, err := s.NextNoWait()

			if err == nil {
				c <- msg
				return
			}

			<-time.After(time.Millisecond) // enter blocking state to let the message get sent
		}

	}()

	return c

}

func (s *Store) IsEmpty() bool {
	return s.Count() == 0
}

func (s *Store) All() []interface{} {
	s.Msg.RLock()
	all := s.Msg.Array
	s.Msg.RUnlock()
	return all
}

// does not move read counter
func (s *Store) PeekLatest() (interface{}, error) {
	if s.IsEmpty() {
		return nil, errors.New("empty store")
	}
	s.Msg.RLock()
	latest := s.Msg.Array[len(s.Msg.Array)-1]
	s.Msg.RUnlock()
	return latest, nil
}

// does not move read counter
func (s *Store) PeekLastRead() (interface{}, error) {
	if s.IsEmpty() {
		return nil, errors.New("empty store")
	}
	s.Msg.RLock()
	s.LastRead.RLock()
	last := s.Msg.Array[s.LastRead.Idx]
	s.Msg.RUnlock()
	s.LastRead.RUnlock()
	return last, nil
}

func (s *Store) LastReadIndex() (int, error) {
	if s.IsEmpty() {
		return 0, errors.New("empty store")
	}
	s.LastRead.RLock()
	idx := s.LastRead.Idx
	s.LastRead.RUnlock()
	if idx < 0 {
		return idx, errors.New("no reads yet")
	}

	return idx, nil
}

func (s *Store) Count() int {
	s.Msg.RLock()
	count := len(s.Msg.Array)
	s.Msg.RUnlock()
	return count
}
