package local

import (
	"testing"

	"github.com/cheekybits/is"
)

func TestNotifier_PostsToChannel(t *testing.T) {
	is := is.New(t)
	ch := NewChannel()
	n := NewNotifier(ch)

	err := n.Notify("chan1", "hello world")
	is.NoErr(err)

	msgs, err := ch.History("chan1")
	is.NoErr(err)
	is.Equal(1, len(msgs))
	is.Equal("[notify] hello world", msgs[0])
}
