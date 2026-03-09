package slack

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cheekybits/is"
)

// ---- store ------------------------------------------------------------------

func TestNewStore_CreatesDirectory(t *testing.T) {
	is := is.New(t)
	dir := filepath.Join(t.TempDir(), "sub", "dir")
	_, err := newStore(dir)
	is.NoErr(err)
	_, err = os.Stat(dir)
	is.NoErr(err)
}

func TestNewStore_InvalidPath_ReturnsError(t *testing.T) {
	is := is.New(t)
	f, err := os.CreateTemp("", "notadir")
	is.NoErr(err)
	f.Close()
	defer os.Remove(f.Name())
	_, err = newStore(filepath.Join(f.Name(), "sub"))
	is.NotNil(err)
}

func TestStore_AppendAndReadAll_RoundTrip(t *testing.T) {
	is := is.New(t)
	st, err := newStore(t.TempDir())
	is.NoErr(err)

	is.NoErr(st.append("key", "line1"))
	is.NoErr(st.append("key", "line2"))

	lines, err := st.readAll("key")
	is.NoErr(err)
	is.Equal(len(lines), 2)
	is.Equal(lines[0], "line1")
	is.Equal(lines[1], "line2")
}

func TestStore_ReadAll_MissingFile_ReturnsNil(t *testing.T) {
	is := is.New(t)
	st, _ := newStore(t.TempDir())
	lines, err := st.readAll("nonexistent")
	is.NoErr(err)
	is.Equal(len(lines), 0)
}

func TestStore_ReadAll_SkipsEmptyLines(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	st, _ := newStore(dir)
	os.WriteFile(filepath.Join(dir, "key.jsonl"), []byte("line1\n\n  \nline2\n"), 0o644)

	lines, err := st.readAll("key")
	is.NoErr(err)
	is.Equal(len(lines), 2)
}

func TestStore_Append_CannotOpenFile_ReturnsError(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	st, _ := newStore(dir)
	os.MkdirAll(filepath.Join(dir, "key.jsonl"), 0o755)
	err := st.append("key", "data")
	is.NotNil(err)
}

func TestStore_Append_WriteError_ReturnsError(t *testing.T) {
	is := is.New(t)
	st, _ := newStore(t.TempDir())
	st.fprintlnFn = func(_ io.Writer, _ string) error { return errors.New("disk full") }
	err := st.append("key", "value")
	is.NotNil(err)
}

func TestStore_ReadAll_OpenError_ReturnsError(t *testing.T) {
	is := is.New(t)
	st, _ := newStore(t.TempDir())
	st.openFn = func(name string) (*os.File, error) {
		return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrPermission}
	}
	_, err := st.readAll("somekey")
	is.NotNil(err)
}

func TestStore_ReadAll_ScannerError_ReturnsError(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	st, _ := newStore(dir)
	bigLine := strings.Repeat("x", 65537)
	os.WriteFile(filepath.Join(dir, "big.jsonl"), []byte(bigLine+"\n"), 0o644)
	_, err := st.readAll("big")
	is.NotNil(err)
}

// ---- helpers ----------------------------------------------------------------

// okJSON is a valid Slack ok:true JSON response.
const okJSON = `{"ok":true}`

// newTestChannel returns a Channel wired to srv with a temp-dir store.
func newTestChannel(t *testing.T, srv *httptest.Server) *Channel {
	t.Helper()
	st, err := newStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return newWith(srv.URL, "test-token", st, srv.Client())
}

// mockServer returns an httptest server that always responds with statusCode
// and the given body.
func mockServer(t *testing.T, statusCode int, body string) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(statusCode)
		w.Write([]byte(body)) //nolint:errcheck
	}))
	t.Cleanup(srv.Close)
	return srv
}

// ---- New --------------------------------------------------------------------

func TestNew_InvalidDataDir_ReturnsError(t *testing.T) {
	is := is.New(t)
	f, err := os.CreateTemp("", "notadir")
	is.NoErr(err)
	f.Close()
	defer os.Remove(f.Name())
	_, err = New("token", filepath.Join(f.Name(), "sub"))
	is.NotNil(err)
}

func TestNew_ValidDir_ReturnsChannel(t *testing.T) {
	is := is.New(t)
	ch, err := New("token", t.TempDir())
	is.NoErr(err)
	is.NotNil(ch)
}

// ---- Channel.Post -----------------------------------------------------------

func TestChannel_Post_SendsMessageAndStores(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK, okJSON)
	ch := newTestChannel(t, srv)

	is.NoErr(ch.Post("C100", "hello"))

	msgs, err := ch.History("C100")
	is.NoErr(err)
	is.Equal(len(msgs), 1)
	is.Equal(msgs[0], "hello")
}

func TestChannel_Post_APIError_ReturnsError(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusInternalServerError, "")
	ch := newTestChannel(t, srv)
	is.NotNil(ch.Post("C100", "hello"))
}

func TestChannel_Post_SlackNotOK_ReturnsError(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK, `{"ok":false,"error":"channel_not_found"}`)
	ch := newTestChannel(t, srv)
	is.NotNil(ch.Post("C100", "hello"))
}

func TestChannel_Post_MarshalError_ReturnsError(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK, okJSON)
	ch := newTestChannel(t, srv)
	ch.jsonMarshalFn = func(any) ([]byte, error) { return nil, errors.New("marshal fail") }
	is.NotNil(ch.Post("C100", "text"))
}

func TestChannel_Post_NewRequestError_ReturnsError(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK, okJSON)
	ch := newTestChannel(t, srv)
	ch.newRequestFn = func(_, _ string, _ io.Reader) (*http.Request, error) {
		return nil, errors.New("bad URL")
	}
	is.NotNil(ch.Post("C100", "text"))
}

func TestChannel_Post_UnmarshalError_ReturnsError(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK, `{bad json`)
	ch := newTestChannel(t, srv)
	is.NotNil(ch.Post("C100", "text"))
}

func TestChannel_Post_NetworkError_ReturnsError(t *testing.T) {
	is := is.New(t)
	st, _ := newStore(t.TempDir())
	ch := newWith("http://127.0.0.1:0", "tok", st, &http.Client{})
	is.NotNil(ch.Post("C100", "text"))
}

// ---- Channel.History --------------------------------------------------------

func TestChannel_History_EmptyChannel_ReturnsNil(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK, okJSON)
	ch := newTestChannel(t, srv)
	msgs, err := ch.History("unknown")
	is.NoErr(err)
	is.Equal(len(msgs), 0)
}

// ---- Channel.SendDM ---------------------------------------------------------

func TestChannel_SendDM_SendsAndStores(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK, okJSON)
	ch := newTestChannel(t, srv)

	is.NoErr(ch.SendDM("U42", "dm text"))

	msgs, err := ch.DMHistory("U42")
	is.NoErr(err)
	is.Equal(len(msgs), 1)
	is.Equal(msgs[0], "dm text")
}

func TestChannel_SendDM_APIError_ReturnsError(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusBadGateway, "")
	ch := newTestChannel(t, srv)
	is.NotNil(ch.SendDM("U42", "text"))
}

// ---- Channel.DMHistory ------------------------------------------------------

func TestChannel_DMHistory_EmptyUser_ReturnsNil(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK, okJSON)
	ch := newTestChannel(t, srv)
	msgs, err := ch.DMHistory("nobody")
	is.NoErr(err)
	is.Equal(len(msgs), 0)
}

// ---- Channel.PostImage ------------------------------------------------------

func TestChannel_PostImage_UploadsFile(t *testing.T) {
	is := is.New(t)
	var receivedPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(okJSON)) //nolint:errcheck
	}))
	t.Cleanup(srv.Close)
	st, _ := newStore(t.TempDir())
	ch := newWith(srv.URL, "tok", st, srv.Client())

	is.NoErr(ch.PostImage("C100", []byte("fake-jpeg-bytes")))
	if !strings.Contains(receivedPath, "files.upload") {
		t.Errorf("expected files.upload request, got path %q", receivedPath)
	}
}

func TestChannel_PostImage_APIError_ReturnsError(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusBadRequest, "")
	ch := newTestChannel(t, srv)
	is.NotNil(ch.PostImage("C100", []byte("data")))
}

func TestChannel_PostImage_NetworkError_ReturnsError(t *testing.T) {
	is := is.New(t)
	st, _ := newStore(t.TempDir())
	ch := newWith("http://127.0.0.1:0", "tok", st, &http.Client{})
	is.NotNil(ch.PostImage("C100", []byte("data")))
}

func TestChannel_PostImage_NewRequestError_ReturnsError(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK, okJSON)
	ch := newTestChannel(t, srv)
	ch.newRequestFn = func(_, _ string, _ io.Reader) (*http.Request, error) {
		return nil, errors.New("bad URL")
	}
	is.NotNil(ch.PostImage("C100", []byte("data")))
}

// ---- ParseSlashCommand ------------------------------------------------------

func makeSlashBody(command, text, channelID, userID, channelType string) []byte {
	v := url.Values{}
	v.Set("command", command)
	v.Set("text", text)
	v.Set("channel_id", channelID)
	v.Set("user_id", userID)
	v.Set("channel_type", channelType)
	return []byte(v.Encode())
}

func TestParseSlashCommand_GroupCommand_ReturnsBotCommand(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK, okJSON)
	ch := newTestChannel(t, srv)

	cmd, ok := ch.ParseSlashCommand(makeSlashBody("/start", "", "C100", "U42", "channel"))
	is.Equal(ok, true)
	is.Equal(cmd.Name, "start")
	is.Equal(cmd.UserID, "U42")
	is.Equal(cmd.ChannelID, "C100")
	is.Equal(cmd.IsDM, false)
	is.Equal(cmd.GameChannelID, "")
}

func TestParseSlashCommand_GroupCommand_RecordsUserChannel(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK, okJSON)
	ch := newTestChannel(t, srv)

	ch.ParseSlashCommand(makeSlashBody("/join", "England", "C100", "U42", "channel"))

	cmd, ok := ch.ParseSlashCommand(makeSlashBody("/order", "A Vie-Bud", "D99", "U42", "im"))
	is.Equal(ok, true)
	is.Equal(cmd.IsDM, true)
	is.Equal(cmd.GameChannelID, "C100")
}

func TestParseSlashCommand_DMCommand_NoChannelRecord_EmptyGameChannelID(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK, okJSON)
	ch := newTestChannel(t, srv)

	cmd, ok := ch.ParseSlashCommand(makeSlashBody("/orders", "", "D99", "U42", "im"))
	is.Equal(ok, true)
	is.Equal(cmd.IsDM, true)
	is.Equal(cmd.GameChannelID, "")
}

func TestParseSlashCommand_CommandWithArgs_ParsesArgs(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK, okJSON)
	ch := newTestChannel(t, srv)

	cmd, ok := ch.ParseSlashCommand(makeSlashBody("/join", "England", "C100", "U42", "channel"))
	is.Equal(ok, true)
	is.Equal(cmd.Name, "join")
	is.Equal(len(cmd.Args), 1)
	is.Equal(cmd.Args[0], "England")
}

func TestParseSlashCommand_EmptyText_ReturnsNilArgs(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK, okJSON)
	ch := newTestChannel(t, srv)

	cmd, ok := ch.ParseSlashCommand(makeSlashBody("/newgame", "", "C100", "U42", "channel"))
	is.Equal(ok, true)
	is.Equal(len(cmd.Args), 0)
}

func TestParseSlashCommand_NonCommand_ReturnsFalse(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK, okJSON)
	ch := newTestChannel(t, srv)

	// body with no "command" field starting with "/"
	body := []byte("command=hello&channel_id=C100&user_id=U42&channel_type=channel")
	_, ok := ch.ParseSlashCommand(body)
	is.Equal(ok, false)
}

func TestParseSlashCommand_MalformedQuery_ReturnsFalse(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK, okJSON)
	ch := newTestChannel(t, srv)

	// %zz is invalid percent-encoding; url.ParseQuery returns an error.
	_, ok := ch.ParseSlashCommand([]byte("command=%2Forder&bad=%zz"))
	is.Equal(ok, false)
}

// ---- Notifier ---------------------------------------------------------------

func TestNewNotifier_And_Notify(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK, okJSON)
	ch := newTestChannel(t, srv)
	n := NewNotifier(ch)
	is.NotNil(n)
	is.NoErr(n.Notify("C100", "phase resolved"))
	msgs, err := ch.History("C100")
	is.NoErr(err)
	is.Equal(len(msgs), 1)
}

// ---- VerifyRequest ----------------------------------------------------------

func makeSignature(secret, timestamp string, body []byte) string {
	base := "v0:" + timestamp + ":" + string(body)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(base)) //nolint:errcheck
	return "v0=" + hex.EncodeToString(mac.Sum(nil))
}

func TestVerifyRequest_ValidSignature_ReturnsTrue(t *testing.T) {
	is := is.New(t)
	body := []byte("command=%2Forder&text=A+Vie-Bud")
	sig := makeSignature("mysecret", "1609459200", body)
	is.Equal(VerifyRequest("mysecret", "1609459200", body, sig), true)
}

func TestVerifyRequest_WrongSecret_ReturnsFalse(t *testing.T) {
	is := is.New(t)
	body := []byte("command=%2Forder")
	sig := makeSignature("correct-secret", "1609459200", body)
	is.Equal(VerifyRequest("wrong-secret", "1609459200", body, sig), false)
}

func TestVerifyRequest_TamperedBody_ReturnsFalse(t *testing.T) {
	is := is.New(t)
	body := []byte("command=%2Forder")
	sig := makeSignature("mysecret", "1609459200", body)
	is.Equal(VerifyRequest("mysecret", "1609459200", []byte("tampered"), sig), false)
}

func TestVerifyRequest_NoV0Prefix_ReturnsFalse(t *testing.T) {
	is := is.New(t)
	body := []byte("command=%2Forder")
	is.Equal(VerifyRequest("mysecret", "1609459200", body, "sha256=abc"), false)
}
