package ghook

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	secret1 = []byte("d696f82431664d9ea93483789db0116c")
	secret2 = []byte("e177b4e2ce1d2604d05e79a8375031b6")
)

func testHookReq(url string, secret []byte) error {
	body := []byte("this is a test")
	mac := hmac.New(sha1.New, secret)
	mac.Write(body)
	sig := mac.Sum(nil)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-GitHub-Event", "ping")
	req.Header.Set("X-GitHub-Delivery", "empty")
	req.Header.Set("X-Hub-Signature", fmt.Sprintf("sha1=%s", hex.EncodeToString(sig)))

	// test proper digest
	res, err := http.DefaultClient.Do(req)
	if err != nil || res.StatusCode != http.StatusOK {
		msg, _ := ioutil.ReadAll(res.Body)
		return fmt.Errorf("err: %s, status: %d, msg: %s", err, res.StatusCode, string(msg))
	}

	return res.Body.Close()
}

func TestHook(t *testing.T) {
	hook := New(secret1, func(event *Event) error {
		fmt.Printf("received %s event!\n", event.Name)
		return nil
	})
	server := httptest.NewServer(hook)
	defer server.Close()

	// test good digest
	if err := testHookReq(server.URL, secret1); err != nil {
		t.Fatal(err)
	}

	// test bad digest
	if err := testHookReq(server.URL, secret2); err == nil {
		t.Fatal("bad digest validation")
	}
}
