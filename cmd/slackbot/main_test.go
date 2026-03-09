package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/cheekybits/is"
	"github.com/burrbd/dip/bot"
	"github.com/burrbd/dip/engine"
	"github.com/burrbd/dip/platform/slack"
)

// ---- test helpers -----------------------------------------------------------

const testSecret = "test-signing-secret"

func signBody(secret, timestamp string, body []byte) string {
	base := "v0:" + timestamp + ":" + string(body)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(base)) //nolint:errcheck
	return "v0=" + hex.EncodeToString(mac.Sum(nil))
}

// newTestSetup creates a Channel backed by a mock Slack API server, and a
// Dispatcher wired to it. Returns the channel, dispatcher, and mock server.
func newTestSetup(t *testing.T, apiStatusCode int, apiBody string) (*slack.Channel, *bot.Dispatcher, *httptest.Server) {
	t.Helper()
	apiSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(apiStatusCode)
		w.Write([]byte(apiBody)) //nolint:errcheck
	}))
	t.Cleanup(apiSrv.Close)

	ch, err := slack.New("test-token", t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	notifier := slack.NewNotifier(ch)
	d := bot.New(ch, notifier, engine.Load, engine.New)
	return ch, d, apiSrv
}

func makeSlashRequest(t *testing.T, secret, command, text, channelID, userID, channelType string) *http.Request {
	t.Helper()
	v := url.Values{}
	v.Set("command", command)
	v.Set("text", text)
	v.Set("channel_id", channelID)
	v.Set("user_id", userID)
	v.Set("channel_type", channelType)
	body := []byte(v.Encode())
	ts := "1609459200"
	sig := signBody(secret, ts, body)
	req := httptest.NewRequest(http.MethodPost, "/slash", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-Slack-Request-Timestamp", ts)
	req.Header.Set("X-Slack-Signature", sig)
	return req
}

// ---- makeSlashHandler tests -------------------------------------------------

func TestSlashHandler_MethodNotAllowed(t *testing.T) {
	is := is.New(t)
	ch, d, _ := newTestSetup(t, http.StatusOK, `{"ok":true}`)
	h := makeSlashHandler(ch, d, testSecret)

	req := httptest.NewRequest(http.MethodGet, "/slash", nil)
	rr := httptest.NewRecorder()
	h(rr, req)
	is.Equal(rr.Code, http.StatusMethodNotAllowed)
}

func TestSlashHandler_InvalidSignature_Unauthorized(t *testing.T) {
	is := is.New(t)
	ch, d, _ := newTestSetup(t, http.StatusOK, `{"ok":true}`)
	h := makeSlashHandler(ch, d, testSecret)

	body := []byte("command=%2Fstart&channel_id=C100&user_id=U1&channel_type=channel")
	req := httptest.NewRequest(http.MethodPost, "/slash", strings.NewReader(string(body)))
	req.Header.Set("X-Slack-Signature", "v0=badsig")
	req.Header.Set("X-Slack-Request-Timestamp", "1609459200")
	rr := httptest.NewRecorder()
	h(rr, req)
	is.Equal(rr.Code, http.StatusUnauthorized)
}

func TestSlashHandler_NonCommand_Returns200(t *testing.T) {
	is := is.New(t)
	ch, d, _ := newTestSetup(t, http.StatusOK, `{"ok":true}`)
	h := makeSlashHandler(ch, d, testSecret)

	body := []byte("command=notaslash&channel_id=C100&user_id=U1&channel_type=channel")
	ts := "1609459200"
	sig := signBody(testSecret, ts, body)
	req := httptest.NewRequest(http.MethodPost, "/slash", strings.NewReader(string(body)))
	req.Header.Set("X-Slack-Request-Timestamp", ts)
	req.Header.Set("X-Slack-Signature", sig)
	rr := httptest.NewRecorder()
	h(rr, req)
	is.Equal(rr.Code, http.StatusOK)
}

func TestSlashHandler_UnknownCommand_PostsError(t *testing.T) {
	is := is.New(t)
	ch, d, _ := newTestSetup(t, http.StatusOK, `{"ok":true}`)
	h := makeSlashHandler(ch, d, testSecret)

	req := makeSlashRequest(t, testSecret, "/unknown", "", "C100", "U1", "channel")
	rr := httptest.NewRecorder()
	h(rr, req)
	is.Equal(rr.Code, http.StatusOK)
}

func TestSlashHandler_ValidCommand_PostsResponse(t *testing.T) {
	is := is.New(t)
	ch, d, _ := newTestSetup(t, http.StatusOK, `{"ok":true}`)
	h := makeSlashHandler(ch, d, testSecret)

	// /help always returns a response without needing game state
	req := makeSlashRequest(t, testSecret, "/help", "", "C100", "U1", "channel")
	rr := httptest.NewRecorder()
	h(rr, req)
	is.Equal(rr.Code, http.StatusOK)
}

// ---- makeEventsHandler tests ------------------------------------------------

func TestEventsHandler_MethodNotAllowed(t *testing.T) {
	is := is.New(t)
	h := makeEventsHandler(testSecret)

	req := httptest.NewRequest(http.MethodGet, "/events", nil)
	rr := httptest.NewRecorder()
	h(rr, req)
	is.Equal(rr.Code, http.StatusMethodNotAllowed)
}

func TestEventsHandler_InvalidSignature_Unauthorized(t *testing.T) {
	is := is.New(t)
	h := makeEventsHandler(testSecret)

	body := []byte(`{"type":"url_verification","challenge":"abc"}`)
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(string(body)))
	req.Header.Set("X-Slack-Signature", "v0=badsig")
	req.Header.Set("X-Slack-Request-Timestamp", "1609459200")
	rr := httptest.NewRecorder()
	h(rr, req)
	is.Equal(rr.Code, http.StatusUnauthorized)
}

func TestEventsHandler_URLVerification_ReturnsChallenge(t *testing.T) {
	is := is.New(t)
	h := makeEventsHandler(testSecret)

	body := []byte(`{"type":"url_verification","challenge":"my-challenge-token"}`)
	ts := "1609459200"
	sig := signBody(testSecret, ts, body)
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(string(body)))
	req.Header.Set("X-Slack-Request-Timestamp", ts)
	req.Header.Set("X-Slack-Signature", sig)
	rr := httptest.NewRecorder()
	h(rr, req)

	is.Equal(rr.Code, http.StatusOK)
	var resp map[string]string
	json.NewDecoder(rr.Body).Decode(&resp) //nolint:errcheck
	is.Equal(resp["challenge"], "my-challenge-token")
}

func TestEventsHandler_MalformedJSON_ReturnsBadRequest(t *testing.T) {
	is := is.New(t)
	h := makeEventsHandler(testSecret)

	body := []byte(`{bad json`)
	ts := "1609459200"
	sig := signBody(testSecret, ts, body)
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(string(body)))
	req.Header.Set("X-Slack-Request-Timestamp", ts)
	req.Header.Set("X-Slack-Signature", sig)
	rr := httptest.NewRecorder()
	h(rr, req)
	is.Equal(rr.Code, http.StatusBadRequest)
}

func TestEventsHandler_OtherEventType_Returns200(t *testing.T) {
	is := is.New(t)
	h := makeEventsHandler(testSecret)

	body := []byte(`{"type":"event_callback"}`)
	ts := "1609459200"
	sig := signBody(testSecret, ts, body)
	req := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(string(body)))
	req.Header.Set("X-Slack-Request-Timestamp", ts)
	req.Header.Set("X-Slack-Signature", sig)
	rr := httptest.NewRecorder()
	h(rr, req)
	is.Equal(rr.Code, http.StatusOK)
}

// ---- body read error paths --------------------------------------------------

// errReader is an io.Reader that always returns an error.
type errReader struct{}

func (errReader) Read(_ []byte) (int, error) { return 0, errors.New("read error") }

func TestSlashHandler_ReadBodyError_ReturnsBadRequest(t *testing.T) {
	is := is.New(t)
	ch, d, _ := newTestSetup(t, http.StatusOK, `{"ok":true}`)
	h := makeSlashHandler(ch, d, testSecret)

	req := httptest.NewRequest(http.MethodPost, "/slash", errReader{})
	rr := httptest.NewRecorder()
	h(rr, req)
	is.Equal(rr.Code, http.StatusBadRequest)
}

func TestEventsHandler_ReadBodyError_ReturnsBadRequest(t *testing.T) {
	is := is.New(t)
	h := makeEventsHandler(testSecret)

	req := httptest.NewRequest(http.MethodPost, "/events", errReader{})
	rr := httptest.NewRecorder()
	h(rr, req)
	is.Equal(rr.Code, http.StatusBadRequest)
}

// ---- envOr ------------------------------------------------------------------

func TestEnvOr_SetVar_ReturnsValue(t *testing.T) {
	is := is.New(t)
	t.Setenv("TEST_KEY_SLACKBOT", "myvalue")
	is.Equal(envOr("TEST_KEY_SLACKBOT", "default"), "myvalue")
}

func TestEnvOr_UnsetVar_ReturnsDefault(t *testing.T) {
	is := is.New(t)
	is.Equal(envOr("TEST_KEY_SLACKBOT_UNSET", "default"), "default")
}
