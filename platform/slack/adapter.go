// Package slack adapts incoming Slack slash commands and Events API payloads
// into bot.Command values, and sends bot responses back to Slack channels.
package slack

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"github.com/burrbd/dip/bot"
)

// Channel implements events.Channel for Slack. Text messages are persisted to
// a local JSONL file store so that event history can be replayed. Images are
// uploaded via files.upload without local storage.
type Channel struct {
	apiURL string
	token  string
	st     *store
	client *http.Client

	userChannelMu  sync.Mutex
	userChannelMap map[string]string // userID → game channelID

	// injectable for testing error branches
	jsonMarshalFn   func(any) ([]byte, error)
	jsonUnmarshalFn func([]byte, any) error
	newRequestFn    func(method, url string, body io.Reader) (*http.Request, error)
}

// New returns a Channel that posts to Slack using the given bot token and
// persists history in the given data directory.
func New(token, dataDir string) (*Channel, error) {
	st, err := newStore(dataDir)
	if err != nil {
		return nil, err
	}
	return newWith("https://slack.com/api", token, st, &http.Client{}), nil
}

// newWith is the testable constructor with injectable API URL and HTTP client.
func newWith(apiURL, token string, st *store, client *http.Client) *Channel {
	return &Channel{
		apiURL:          apiURL,
		token:           token,
		st:              st,
		client:          client,
		userChannelMap:  make(map[string]string),
		jsonMarshalFn:   json.Marshal,
		jsonUnmarshalFn: json.Unmarshal,
		newRequestFn:    http.NewRequest,
	}
}

// Post sends a text message to channelID via chat.postMessage and persists it.
func (c *Channel) Post(channelID, text string) error {
	if err := c.chatPostMessage(channelID, text); err != nil {
		return err
	}
	return c.st.append("ch_"+channelID, text)
}

// History returns all text messages previously posted to channelID.
func (c *Channel) History(channelID string) ([]string, error) {
	return c.st.readAll("ch_" + channelID)
}

// SendDM sends a private text message to userID via chat.postMessage and persists it.
func (c *Channel) SendDM(userID, text string) error {
	if err := c.chatPostMessage(userID, text); err != nil {
		return err
	}
	return c.st.append("dm_"+userID, text)
}

// DMHistory returns all DM text messages previously sent to userID.
func (c *Channel) DMHistory(userID string) ([]string, error) {
	return c.st.readAll("dm_" + userID)
}

// PostImage uploads a binary image to channelID via files.upload.
func (c *Channel) PostImage(channelID string, data []byte) error {
	return c.filesUpload(channelID, data)
}

// ParseSlashCommand parses a raw Slack slash-command webhook body into a
// bot.Command. Returns the command and true when the body contains a valid
// slash command. Malformed bodies or non-command payloads return false.
//
// When the command comes from a public channel, the user→channel mapping is
// recorded so that subsequent DM commands from the same user can resolve the
// game channel.
func (c *Channel) ParseSlashCommand(body []byte) (bot.Command, bool) {
	vals, err := url.ParseQuery(string(body))
	if err != nil {
		return bot.Command{}, false
	}

	command := vals.Get("command")
	if !strings.HasPrefix(command, "/") {
		return bot.Command{}, false
	}

	name := strings.TrimPrefix(command, "/")
	text := vals.Get("text")
	channelID := vals.Get("channel_id")
	userID := vals.Get("user_id")
	channelType := vals.Get("channel_type")
	isDM := channelType == "im"

	var args []string
	if text != "" {
		args = strings.Fields(text)
	}

	if !isDM {
		c.setUserChannel(userID, channelID)
	}

	gameChannelID := ""
	if isDM {
		gameChannelID, _ = c.getUserChannel(userID)
	}

	return bot.Command{
		Name:          name,
		Args:          args,
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

// chatPostMessage calls the Slack chat.postMessage API.
func (c *Channel) chatPostMessage(channelID, text string) error {
	return c.apiPost("chat.postMessage", map[string]any{"channel": channelID, "text": text})
}

// filesUpload uploads data as an image to channelID via files.upload.
func (c *Channel) filesUpload(channelID string, data []byte) error {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	w.WriteField("channels", channelID) //nolint:errcheck // bytes.Buffer never fails
	fw, _ := w.CreateFormFile("file", "map.jpg")
	fw.Write(data) //nolint:errcheck // bytes.Buffer never fails
	w.Close()      //nolint:errcheck // bytes.Buffer never fails

	req, err := c.newRequestFn(http.MethodPost, c.apiURL+"/files.upload", &body)
	if err != nil {
		return fmt.Errorf("slack: build files.upload request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return c.doRequest(req)
}

// apiPost marshals payload as JSON and POSTs it to the given Slack API method.
func (c *Channel) apiPost(method string, payload any) error {
	data, err := c.jsonMarshalFn(payload)
	if err != nil {
		return fmt.Errorf("slack: marshal %s payload: %w", method, err)
	}
	req, err := c.newRequestFn(http.MethodPost, c.apiURL+"/"+method, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("slack: build %s request: %w", method, err)
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	return c.doRequest(req)
}

// slackResponse is the common JSON envelope returned by Slack API calls.
type slackResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error"`
}

// doRequest executes req, drains the body, and returns an error for non-200
// responses or when the Slack API returns ok:false.
func (c *Channel) doRequest(req *http.Request) error {
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("slack: %s: %w", req.URL.Path, err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("slack: %s: status %d", req.URL.Path, resp.StatusCode)
	}
	var sr slackResponse
	if err := c.jsonUnmarshalFn(body, &sr); err != nil {
		return fmt.Errorf("slack: %s: decode response: %w", req.URL.Path, err)
	}
	if !sr.OK {
		return fmt.Errorf("slack: %s: api error: %s", req.URL.Path, sr.Error)
	}
	return nil
}

// VerifyRequest validates a Slack request using the signing secret.
// It returns true when the signature matches the expected HMAC-SHA256 digest.
func VerifyRequest(signingSecret, timestamp string, body []byte, signature string) bool {
	if !strings.HasPrefix(signature, "v0=") {
		return false
	}
	base := "v0:" + timestamp + ":" + string(body)
	mac := hmac.New(sha256.New, []byte(signingSecret))
	mac.Write([]byte(base)) //nolint:errcheck // hash.Hash.Write never returns an error
	expected := "v0=" + hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(signature), []byte(expected))
}
