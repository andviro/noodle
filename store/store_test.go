package store

import (
	"golang.org/x/net/context"
	"gopkg.in/tylerb/is.v1"
	"testing"
	"time"
)

func TestGetSet(t *testing.T) {
	is := is.New(t)
	s := New()
	s.Set("existingKey", 10)
	val, ok := s.Get("existingKey")
	is.True(ok)
	is.Equal(val.(int), 10)
	_, ok = s.Get("not existingKey")
	is.False(ok)
}

func TestMust(t *testing.T) {
	is := is.New(t)
	s := New()
	s.Set("existingKey", 10)
	is.Equal(s.Must("existingKey").(int), 10)

}

func TestMustThrowsError(t *testing.T) {
	is := is.New(t)
	var err error
	func(s *Store) {
		defer func() {
			err = recover().(error)
		}()
		_ = s.Must("not existingKey")
	}(New())
	is.Err(err)
	is.Equal(err, NotFoundError)
}

func TestReadWaitsForUpdate(t *testing.T) {
	is := is.New(t)
	s := New()
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
	err := s.Read(func(t map[string]interface{}) error {
		x, ok := t["key"].(int)
		is.True(ok)
		is.Equal(x, 1)
		return nil
	})
	is.NotErr(err)
}

func TestGlobalLocal(t *testing.T) {
	is := is.New(t)
	root := context.Background()
	ctx := WithLocal(WithGlobal(root))
	_, ok := Global(ctx)
	is.True(ok)
	_, ok = Local(ctx)
	is.True(ok)
}
