package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/burrbd/dip/bot"
)

// Update is a Telegram Bot API update payload.
type Update struct {
	UpdateID int      `json:"update_id"`
	Message  *Message `json:"message,omitempty"`
}

// Message is a Telegram message.
type Message struct {
	MessageID int    `json:"message_id"`
	From      User   `json:"from"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text"`
}

// User is a Telegram user.
type User struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

// Chat is a Telegram chat.
type Chat struct {
	ID   int64  `json:"id"`
	Type string `json:"type"` // "private", "group", "supergroup", "channel"
}

// Channel implements events.Channel for Telegram. Text messages are persisted
// to a local JSONL file store because the Telegram Bot API does not expose
// historical messages. Images are sent via sendPhoto without local storage.
type Channel struct {
	apiURL string
	store  *Store
	client *http.Client

	userChannelMu  sync.Mutex
	userChannelMap map[string]string // userID → game channelID

	// jsonMarshalFn is injectable to cover the marshal error branch in apiPost.
	jsonMarshalFn func(any) ([]byte, error)
	// newRequestFn is injectable to cover the http.NewRequest error branches.
	newRequestFn func(method, url string, body io.Reader) (*http.Request, error)
}

// New returns a Channel that posts to Telegram using the given bot token
// and persists history in the given store.
func New(token string, store *Store) *Channel {
	ch := newWith("https://api.telegram.org/bot"+token, store, &http.Client{})
	return ch
}

// newWith is the testable constructor with injectable API URL and HTTP client.
func newWith(apiURL string, store *Store, client *http.Client) *Channel {
	return &Channel{
		apiURL:         apiURL,
		store:          store,
		client:         client,
		userChannelMap: make(map[string]string),
		jsonMarshalFn:  json.Marshal,
		newRequestFn:   http.NewRequest,
	}
}

// Post sends a text message to channelID and persists it to the local store.
func (c *Channel) Post(channelID, text string) error {
	if err := c.sendMessage(channelID, text); err != nil {
		return err
	}
	return c.store.Append("ch_"+channelID, text)
}

// History returns all text messages previously posted to channelID.
func (c *Channel) History(channelID string) ([]string, error) {
	return c.store.ReadAll("ch_" + channelID)
}

// SendDM sends a private text message to userID and persists it to the local store.
func (c *Channel) SendDM(userID, text string) error {
	if err := c.sendMessage(userID, text); err != nil {
		return err
	}
	return c.store.Append("dm_"+userID, text)
}

// DMHistory returns all DM text messages previously sent to userID.
func (c *Channel) DMHistory(userID string) ([]string, error) {
	return c.store.ReadAll("dm_" + userID)
}

// PostImage sends a binary image to channelID via Telegram sendPhoto.
func (c *Channel) PostImage(channelID string, data []byte) error {
	return c.sendPhoto(channelID, data)
}

// ParseUpdate parses a raw Telegram webhook payload into a bot.Command.
// Returns the command and true when the update contains a bot command (text
// starting with "/"). Non-command messages and malformed payloads return false.
//
// When the message is from a group chat, the user→channel mapping is recorded
// so that subsequent DM commands from the same user resolve to the correct
// game channel.
func (c *Channel) ParseUpdate(body []byte) (bot.Command, bool) {
	var upd Update
	if err := json.Unmarshal(body, &upd); err != nil || upd.Message == nil {
		return bot.Command{}, false
	}
	msg := upd.Message
	if !strings.HasPrefix(msg.Text, "/") {
		return bot.Command{}, false
	}

	userID := strconv.FormatInt(msg.From.ID, 10)
	channelID := strconv.FormatInt(msg.Chat.ID, 10)
	isDM := msg.Chat.Type == "private"

	if !isDM {
		c.setUserChannel(userID, channelID)
	}

	tokens := strings.Fields(msg.Text)
	name := strings.TrimPrefix(tokens[0], "/")
	// Strip @botname suffix: "/order@mybotname" → "order"
	if i := strings.Index(name, "@"); i > 0 {
		name = name[:i]
	}

	gameChannelID := ""
	if isDM {
		gameChannelID, _ = c.getUserChannel(userID)
	}

	return bot.Command{
		Name:          name,
		Args:          tokens[1:],
		UserID:        userID,
		ChannelID:     channelID,
		IsDM:          isDM,
		GameChannelID: gameChannelID,
	}, true
}

func (c *Channel) setUserChannel(userID, channelID string) {
	c.userChannelMu.Lock()
	c.userChannelMap[userID] = channelID
	c.userChannelMu.Unlock()
}

func (c *Channel) getUserChannel(userID string) (string, bool) {
	c.userChannelMu.Lock()
	v, ok := c.userChannelMap[userID]
	c.userChannelMu.Unlock()
	return v, ok
}

// sendMessage calls the Telegram sendMessage API.
func (c *Channel) sendMessage(chatID, text string) error {
	return c.apiPost("sendMessage", map[string]any{"chat_id": chatID, "text": text})
}

// sendPhoto calls the Telegram sendPhoto API using multipart/form-data.
// WriteField, CreateFormFile, fw.Write, and w.Close all write to a bytes.Buffer
// and cannot return errors; they are called without error checks.
func (c *Channel) sendPhoto(chatID string, data []byte) error {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	w.WriteField("chat_id", chatID)   //nolint:errcheck // bytes.Buffer never fails
	fw, _ := w.CreateFormFile("photo", "map.jpg")
	fw.Write(data)  //nolint:errcheck // bytes.Buffer never fails
	w.Close()       //nolint:errcheck // bytes.Buffer never fails
	req, err := c.newRequestFn(http.MethodPost, c.apiURL+"/sendPhoto", &body)
	if err != nil {
		return fmt.Errorf("telegram: build sendPhoto request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())
	return c.doRequest(req)
}

// apiPost marshals payload as JSON and POSTs it to the given Telegram API method.
func (c *Channel) apiPost(method string, payload any) error {
	data, err := c.jsonMarshalFn(payload)
	if err != nil {
		return fmt.Errorf("telegram: marshal %s payload: %w", method, err)
	}
	req, err := c.newRequestFn(http.MethodPost, c.apiURL+"/"+method, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("telegram: build %s request: %w", method, err)
	}
	req.Header.Set("Content-Type", "application/json")
	return c.doRequest(req)
}

// doRequest executes req, drains the body, and returns an error for non-200 responses.
func (c *Channel) doRequest(req *http.Request) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("telegram: %s: %w", req.URL.Path, err)
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body) //nolint:errcheck
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram: %s: status %d", req.URL.Path, resp.StatusCode)
	}
	return nil
}
