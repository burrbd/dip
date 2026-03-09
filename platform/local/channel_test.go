package local

import (
	"sync"
	"testing"

	"github.com/cheekybits/is"
)

func TestChannel_PostAndHistory(t *testing.T) {
	is := is.New(t)
	ch := NewChannel()

	err := ch.Post("chan1", "hello")
	is.NoErr(err)
	err = ch.Post("chan1", "world")
	is.NoErr(err)

	msgs, err := ch.History("chan1")
	is.NoErr(err)
	is.Equal(2, len(msgs))
	is.Equal("hello", msgs[0])
	is.Equal("world", msgs[1])
}

func TestChannel_History_EmptyChannel(t *testing.T) {
	is := is.New(t)
	ch := NewChannel()

	msgs, err := ch.History("nonexistent")
	is.NoErr(err)
	is.Equal(0, len(msgs))
}

func TestChannel_SendDMAndDMHistory(t *testing.T) {
	is := is.New(t)
	ch := NewChannel()

	err := ch.SendDM("user1", "dm1")
	is.NoErr(err)
	err = ch.SendDM("user1", "dm2")
	is.NoErr(err)

	msgs, err := ch.DMHistory("user1")
	is.NoErr(err)
	is.Equal(2, len(msgs))
	is.Equal("dm1", msgs[0])
	is.Equal("dm2", msgs[1])
}

func TestChannel_DMHistory_EmptyUser(t *testing.T) {
	is := is.New(t)
	ch := NewChannel()

	msgs, err := ch.DMHistory("nobody")
	is.NoErr(err)
	is.Equal(0, len(msgs))
}

func TestChannel_PostImage(t *testing.T) {
	is := is.New(t)
	ch := NewChannel()

	data := []byte{0x89, 0x50, 0x4E, 0x47}
	err := ch.PostImage("chan1", data)
	is.NoErr(err)

	imgs := ch.ImagesSince("chan1", 0)
	is.Equal(1, len(imgs))
	is.Equal(data, imgs[0])
}

func TestChannel_MessagesSince(t *testing.T) {
	is := is.New(t)
	ch := NewChannel()

	ch.Post("chan1", "a")
	ch.Post("chan1", "b")
	cursor := ch.MessageCount("chan1")
	ch.Post("chan1", "c")
	ch.Post("chan1", "d")

	result := ch.MessagesSince("chan1", cursor)
	is.Equal(2, len(result))
	is.Equal("c", result[0])
	is.Equal("d", result[1])
}

func TestChannel_MessagesSince_CursorAtEnd(t *testing.T) {
	is := is.New(t)
	ch := NewChannel()

	ch.Post("chan1", "a")
	cursor := ch.MessageCount("chan1")

	result := ch.MessagesSince("chan1", cursor)
	is.Equal(0, len(result))
}

func TestChannel_DMsSince(t *testing.T) {
	is := is.New(t)
	ch := NewChannel()

	ch.SendDM("user1", "x")
	cursor := ch.DMCount("user1")
	ch.SendDM("user1", "y")

	result := ch.DMsSince("user1", cursor)
	is.Equal(1, len(result))
	is.Equal("y", result[0])
}

func TestChannel_DMsSince_CursorAtEnd(t *testing.T) {
	is := is.New(t)
	ch := NewChannel()
	ch.SendDM("user1", "x")
	cursor := ch.DMCount("user1")
	result := ch.DMsSince("user1", cursor)
	is.Equal(0, len(result))
}

func TestChannel_ImagesSince(t *testing.T) {
	is := is.New(t)
	ch := NewChannel()

	ch.PostImage("chan1", []byte{1})
	cursor := ch.ImageCount("chan1")
	ch.PostImage("chan1", []byte{2})
	ch.PostImage("chan1", []byte{3})

	imgs := ch.ImagesSince("chan1", cursor)
	is.Equal(2, len(imgs))
	is.Equal([]byte{2}, imgs[0])
	is.Equal([]byte{3}, imgs[1])
}

func TestChannel_ImagesSince_CursorAtEnd(t *testing.T) {
	is := is.New(t)
	ch := NewChannel()
	ch.PostImage("chan1", []byte{1})
	cursor := ch.ImageCount("chan1")
	imgs := ch.ImagesSince("chan1", cursor)
	is.Equal(0, len(imgs))
}

func TestChannel_ConcurrentPost(t *testing.T) {
	ch := NewChannel()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			ch.Post("chan1", "msg")
		}()
	}
	wg.Wait()

	msgs, _ := ch.History("chan1")
	if len(msgs) != 50 {
		t.Fatalf("expected 50 messages, got %d", len(msgs))
	}
}
