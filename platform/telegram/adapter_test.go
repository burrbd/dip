package telegram

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/cheekybits/is"
)

// ---- Store ------------------------------------------------------------------

func TestNewStore_CreatesDirectory(t *testing.T) {
	is := is.New(t)
	dir := filepath.Join(t.TempDir(), "sub", "dir")
	_, err := NewStore(dir)
	is.NoErr(err)
	_, err = os.Stat(dir)
	is.NoErr(err)
}

func TestNewStore_InvalidPath_ReturnsError(t *testing.T) {
	is := is.New(t)
	// Use an existing file as the target directory so MkdirAll fails.
	f, err := os.CreateTemp("", "notadir")
	is.NoErr(err)
	f.Close()
	defer os.Remove(f.Name())
	_, err = NewStore(filepath.Join(f.Name(), "sub"))
	is.NotNil(err)
}

func TestStore_AppendAndReadAll_RoundTrip(t *testing.T) {
	is := is.New(t)
	st, err := NewStore(t.TempDir())
	is.NoErr(err)

	is.NoErr(st.Append("key", "line1"))
	is.NoErr(st.Append("key", "line2"))

	lines, err := st.ReadAll("key")
	is.NoErr(err)
	is.Equal(len(lines), 2)
	is.Equal(lines[0], "line1")
	is.Equal(lines[1], "line2")
}

func TestStore_ReadAll_MissingFile_ReturnsNil(t *testing.T) {
	is := is.New(t)
	st, _ := NewStore(t.TempDir())
	lines, err := st.ReadAll("nonexistent")
	is.NoErr(err)
	is.Equal(len(lines), 0)
}

func TestStore_ReadAll_SkipsEmptyLines(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	st, _ := NewStore(dir)
	os.WriteFile(filepath.Join(dir, "key.jsonl"), []byte("line1\n\n  \nline2\n"), 0o644)

	lines, err := st.ReadAll("key")
	is.NoErr(err)
	is.Equal(len(lines), 2)
}

func TestStore_Append_CannotOpenFile_ReturnsError(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	st, _ := NewStore(dir)
	// Replace the expected file path with a directory so OpenFile fails.
	os.MkdirAll(filepath.Join(dir, "key.jsonl"), 0o755)
	err := st.Append("key", "data")
	is.NotNil(err)
}

func TestStore_ReadAll_OpenError_ReturnsError(t *testing.T) {
	is := is.New(t)
	st, _ := NewStore(t.TempDir())
	// Inject a failing open function to cover the os.Open error branch regardless
	// of the OS user running the test (root can open any file by permissions).
	st.openFn = func(string) (*os.File, error) { return nil, errors.New("open fail") }
	// First write a file so IsNotExist is false, then trigger open error via injection.
	st.openFn = func(name string) (*os.File, error) { return nil, &os.PathError{Op: "open", Path: name, Err: os.ErrPermission} }
	_, err := st.ReadAll("somekey")
	is.NotNil(err)
}

// ---- helpers ----------------------------------------------------------------

// newTestChannel returns a Channel wired to srv with a temp-dir Store.
func newTestChannel(t *testing.T, srv *httptest.Server) *Channel {
	t.Helper()
	st, err := NewStore(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return newWith(srv.URL, st, srv.Client())
}

// mockServer returns an httptest server that always responds with statusCode.
func mockServer(t *testing.T, statusCode int) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(statusCode)
	}))
	t.Cleanup(srv.Close)
	return srv
}

// ---- Channel.Post -----------------------------------------------------------

func TestChannel_Post_SendsMessageAndStores(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK)
	ch := newTestChannel(t, srv)

	is.NoErr(ch.Post("-100", "hello"))

	msgs, err := ch.History("-100")
	is.NoErr(err)
	is.Equal(len(msgs), 1)
	is.Equal(msgs[0], "hello")
}

func TestChannel_Post_APIError_ReturnsError(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusInternalServerError)
	ch := newTestChannel(t, srv)
	is.NotNil(ch.Post("-100", "hello"))
}

func TestChannel_Post_MarshalError_ReturnsError(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK)
	ch := newTestChannel(t, srv)
	ch.jsonMarshalFn = func(any) ([]byte, error) { return nil, errors.New("marshal fail") }
	is.NotNil(ch.Post("-100", "text"))
}

// ---- Channel.History --------------------------------------------------------

func TestChannel_History_EmptyChannel_ReturnsNil(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK)
	ch := newTestChannel(t, srv)
	msgs, err := ch.History("unknown")
	is.NoErr(err)
	is.Equal(len(msgs), 0)
}

// ---- Channel.SendDM ---------------------------------------------------------

func TestChannel_SendDM_SendsAndStores(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK)
	ch := newTestChannel(t, srv)

	is.NoErr(ch.SendDM("42", "dm text"))

	msgs, err := ch.DMHistory("42")
	is.NoErr(err)
	is.Equal(len(msgs), 1)
	is.Equal(msgs[0], "dm text")
}

func TestChannel_SendDM_APIError_ReturnsError(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusBadGateway)
	ch := newTestChannel(t, srv)
	is.NotNil(ch.SendDM("42", "text"))
}

// ---- Channel.DMHistory ------------------------------------------------------

func TestChannel_DMHistory_EmptyUser_ReturnsNil(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK)
	ch := newTestChannel(t, srv)
	msgs, err := ch.DMHistory("nobody")
	is.NoErr(err)
	is.Equal(len(msgs), 0)
}

// ---- Channel.PostImage ------------------------------------------------------

func TestChannel_PostImage_SendsPhoto(t *testing.T) {
	is := is.New(t)
	var receivedPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedPath = r.URL.Path
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)
	st, _ := NewStore(t.TempDir())
	ch := newWith(srv.URL, st, srv.Client())

	is.NoErr(ch.PostImage("-100", []byte("fake-jpeg-bytes")))
	if !strings.Contains(receivedPath, "sendPhoto") {
		t.Errorf("expected sendPhoto request, got path %q", receivedPath)
	}
}

func TestChannel_PostImage_APIError_ReturnsError(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusBadRequest)
	ch := newTestChannel(t, srv)
	is.NotNil(ch.PostImage("-100", []byte("data")))
}

func TestChannel_PostImage_NetworkError_ReturnsError(t *testing.T) {
	is := is.New(t)
	st, _ := NewStore(t.TempDir())
	ch := newWith("http://127.0.0.1:0", st, &http.Client{})
	is.NotNil(ch.PostImage("-100", []byte("data")))
}

// ---- doRequest network error ------------------------------------------------

func TestDoRequest_NetworkError_ReturnsError(t *testing.T) {
	is := is.New(t)
	st, _ := NewStore(t.TempDir())
	ch := newWith("http://127.0.0.1:0", st, &http.Client{})
	is.NotNil(ch.Post("-100", "text"))
}

// ---- ParseUpdate ------------------------------------------------------------

func makeUpdate(chatType string, chatID, userID int64, text string) []byte {
	u := Update{
		UpdateID: 1,
		Message: &Message{
			MessageID: 1,
			From:      User{ID: userID},
			Chat:      Chat{ID: chatID, Type: chatType},
			Text:      text,
		},
	}
	b, _ := json.Marshal(u)
	return b
}

func TestParseUpdate_GroupCommand_ReturnsBotCommand(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK)
	ch := newTestChannel(t, srv)

	cmd, ok := ch.ParseUpdate(makeUpdate("group", -100, 42, "/start"))
	is.Equal(ok, true)
	is.Equal(cmd.Name, "start")
	is.Equal(cmd.UserID, "42")
	is.Equal(cmd.ChannelID, "-100")
	is.Equal(cmd.IsDM, false)
	is.Equal(cmd.GameChannelID, "")
}

func TestParseUpdate_GroupCommand_RecordsUserChannel(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK)
	ch := newTestChannel(t, srv)

	// First a group message to register user→channel.
	ch.ParseUpdate(makeUpdate("group", -100, 42, "/join England"))

	// Now a DM from the same user resolves the game channel.
	cmd, ok := ch.ParseUpdate(makeUpdate("private", 42, 42, "/order A Vie-Bud"))
	is.Equal(ok, true)
	is.Equal(cmd.IsDM, true)
	is.Equal(cmd.GameChannelID, "-100")
}

func TestParseUpdate_DMCommand_NoChannelRecord_EmptyGameChannelID(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK)
	ch := newTestChannel(t, srv)

	cmd, ok := ch.ParseUpdate(makeUpdate("private", 42, 42, "/orders"))
	is.Equal(ok, true)
	is.Equal(cmd.IsDM, true)
	is.Equal(cmd.GameChannelID, "")
}

func TestParseUpdate_CommandWithBotnameSuffix_StripsName(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK)
	ch := newTestChannel(t, srv)

	cmd, ok := ch.ParseUpdate(makeUpdate("group", -100, 42, "/order@dipbot A Vie-Bud"))
	is.Equal(ok, true)
	is.Equal(cmd.Name, "order")
	is.Equal(len(cmd.Args), 2) // "/order@dipbot A Vie-Bud" → ["A", "Vie-Bud"]
}

func TestParseUpdate_NonCommandText_ReturnsFalse(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK)
	ch := newTestChannel(t, srv)

	_, ok := ch.ParseUpdate(makeUpdate("group", -100, 42, "hello world"))
	is.Equal(ok, false)
}

func TestParseUpdate_NilMessage_ReturnsFalse(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK)
	ch := newTestChannel(t, srv)

	_, ok := ch.ParseUpdate([]byte(`{"update_id":1}`))
	is.Equal(ok, false)
}

func TestParseUpdate_MalformedJSON_ReturnsFalse(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK)
	ch := newTestChannel(t, srv)

	_, ok := ch.ParseUpdate([]byte(`{bad json`))
	is.Equal(ok, false)
}

func TestParseUpdate_CommandWithArgs_ParsesArgs(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK)
	ch := newTestChannel(t, srv)

	cmd, ok := ch.ParseUpdate(makeUpdate("group", -100, 42, "/join England"))
	is.Equal(ok, true)
	is.Equal(cmd.Name, "join")
	is.Equal(len(cmd.Args), 1)
	is.Equal(cmd.Args[0], "England")
}

// ---- New (production constructor) -------------------------------------------

func TestNew_ReturnsNonNilChannel(t *testing.T) {
	is := is.New(t)
	st, _ := NewStore(t.TempDir())
	ch := New("test-token", st)
	is.NotNil(ch)
}

// ---- Notifier ---------------------------------------------------------------

func TestNewNotifier_And_Notify(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK)
	ch := newTestChannel(t, srv)
	n := NewNotifier(ch)
	is.NotNil(n)
	is.NoErr(n.Notify("-100", "phase resolved"))
	// Verify message was stored.
	msgs, err := ch.History("-100")
	is.NoErr(err)
	is.Equal(len(msgs), 1)
}

// ---- Store.Append write error -----------------------------------------------

func TestStore_Append_WriteError_ReturnsError(t *testing.T) {
	is := is.New(t)
	st, _ := NewStore(t.TempDir())
	st.fprintlnFn = func(_ io.Writer, _ string) error { return errors.New("disk full") }
	err := st.Append("key", "value")
	is.NotNil(err)
}

// ---- Store.ReadAll scanner error (line too long) ----------------------------

func TestStore_ReadAll_ScannerError_ReturnsError(t *testing.T) {
	is := is.New(t)
	dir := t.TempDir()
	st, _ := NewStore(dir)
	// bufio.Scanner default max token = 65536 bytes. A longer line triggers an error.
	bigLine := strings.Repeat("x", 65537)
	os.WriteFile(filepath.Join(dir, "big.jsonl"), []byte(bigLine+"\n"), 0o644)
	_, err := st.ReadAll("big")
	is.NotNil(err)
}

// ---- http.NewRequest error paths -------------------------------------------

func TestChannel_Post_NewRequestError_ReturnsError(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK)
	ch := newTestChannel(t, srv)
	ch.newRequestFn = func(_, _ string, _ io.Reader) (*http.Request, error) {
		return nil, errors.New("bad URL")
	}
	is.NotNil(ch.Post("-100", "text"))
}

func TestChannel_PostImage_NewRequestError_ReturnsError(t *testing.T) {
	is := is.New(t)
	srv := mockServer(t, http.StatusOK)
	ch := newTestChannel(t, srv)
	ch.newRequestFn = func(_, _ string, _ io.Reader) (*http.Request, error) {
		return nil, errors.New("bad URL")
	}
	is.NotNil(ch.PostImage("-100", []byte("data")))
}
