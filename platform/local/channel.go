// Package local provides an in-memory implementation of events.Channel for
// use in the QA REPL and unit tests. No external platform is required.
package local

import "sync"

// Channel implements events.Channel using in-memory slices. It is safe for
// concurrent use. Cursor helpers allow callers to read only new messages added
// since a previous point in time.
type Channel struct {
	mu   sync.Mutex
	msgs map[string][]string // channelID → messages
	dms  map[string][]string // userID → DM messages
	imgs map[string][][]byte // channelID → image byte slices
}

// NewChannel returns a ready-to-use Channel.
func NewChannel() *Channel {
	return &Channel{
		msgs: make(map[string][]string),
		dms:  make(map[string][]string),
		imgs: make(map[string][][]byte),
	}
}

// Post appends text to the channel identified by channelID.
func (c *Channel) Post(channelID, text string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.msgs[channelID] = append(c.msgs[channelID], text)
	return nil
}

// History returns all messages in the channel in chronological order.
func (c *Channel) History(channelID string) ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	msgs := c.msgs[channelID]
	result := make([]string, len(msgs))
	copy(result, msgs)
	return result, nil
}

// SendDM appends text to the DM thread for userID.
func (c *Channel) SendDM(userID, text string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dms[userID] = append(c.dms[userID], text)
	return nil
}

// DMHistory returns all messages in the user's DM thread in chronological order.
func (c *Channel) DMHistory(userID string) ([]string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	msgs := c.dms[userID]
	result := make([]string, len(msgs))
	copy(result, msgs)
	return result, nil
}

// PostImage appends a binary image to the image list for channelID.
func (c *Channel) PostImage(channelID string, data []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	img := make([]byte, len(data))
	copy(img, data)
	c.imgs[channelID] = append(c.imgs[channelID], img)
	return nil
}

// MessageCount returns the total number of messages posted to channelID.
func (c *Channel) MessageCount(channelID string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.msgs[channelID])
}

// MessagesSince returns messages posted to channelID at positions >= cursor.
func (c *Channel) MessagesSince(channelID string, cursor int) []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	msgs := c.msgs[channelID]
	if cursor >= len(msgs) {
		return nil
	}
	result := make([]string, len(msgs)-cursor)
	copy(result, msgs[cursor:])
	return result
}

// DMCount returns the total number of DM messages sent to userID.
func (c *Channel) DMCount(userID string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.dms[userID])
}

// DMsSince returns DM messages for userID at positions >= cursor.
func (c *Channel) DMsSince(userID string, cursor int) []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	msgs := c.dms[userID]
	if cursor >= len(msgs) {
		return nil
	}
	result := make([]string, len(msgs)-cursor)
	copy(result, msgs[cursor:])
	return result
}

// ImageCount returns the total number of images posted to channelID.
func (c *Channel) ImageCount(channelID string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.imgs[channelID])
}

// ImagesSince returns images posted to channelID at positions >= cursor.
func (c *Channel) ImagesSince(channelID string, cursor int) [][]byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	imgs := c.imgs[channelID]
	if cursor >= len(imgs) {
		return nil
	}
	result := make([][]byte, len(imgs)-cursor)
	for i, img := range imgs[cursor:] {
		cp := make([]byte, len(img))
		copy(cp, img)
		result[i] = cp
	}
	return result
}
