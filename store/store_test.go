package store_test

import (
	"github.com/andviro/noodle/store"
	"gopkg.in/tylerb/is.v1"
	"testing"
	"time"
)

func TestGetSet(t *testing.T) {
	is := is.New(t)
	s := store.New()
	s.Set("existingKey", 10)
	val, ok := s.Get("existingKey")
	is.True(ok)
	is.Equal(val.(int), 10)
	_, ok = s.Get("not existingKey")
	is.False(ok)
}

func TestDelete(t *testing.T) {
	is := is.New(t)
	s := store.New()
	s.Set("key", "value")
	is.Equal(s.MustGet("key").(string), "value")
	s.Delete("key")
	_, ok := s.Get("key")
	is.False(ok)
}

func TestMustGet(t *testing.T) {
	is := is.New(t)
	s := store.New()
	s.Set("existingKey", 10)
	is.Equal(s.MustGet("existingKey").(int), 10)
}

func TestMustGetPanics(t *testing.T) {
	is := is.New(t)
	var err error
	func(s *store.Store) {
		defer func() {
			err = recover().(error)
		}()
		_ = s.MustGet("not existingKey")
	}(store.New())
	is.Err(err)
	_, ok := err.(store.KeyError)
	is.True(ok)
}

func TestViewWaitsForUpdate(t *testing.T) {
	is := is.New(t)
	s := store.New()
	begin := make(chan struct{})
	s.Set("key", 0)
	go func() {
		s.Update(func(t map[string]interface{}) error {
			close(begin) // Continue main thread here
			time.Sleep(100 * time.Millisecond)
			x := t["key"].(int)
			x++
			t["key"] = x
			return nil
		})
	}()
	<-begin // Gate to begin read transaction
	err := s.View(func(t map[string]interface{}) error {
		x, ok := t["key"].(int)
		is.True(ok)
		is.Equal(x, 1)
		return nil
	})
	is.NotErr(err)
}
