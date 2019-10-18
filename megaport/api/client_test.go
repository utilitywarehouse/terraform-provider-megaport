package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
)

func testClientServer(handler func(w http.ResponseWriter, r *http.Request)) (*Client, *httptest.Server) {
	s := httptest.NewServer(http.HandlerFunc(handler))
	c := NewClient(s.URL)
	return c, s
}

func TestClient_Login(t *testing.T) {
	username := acctest.RandString(10)
	password := acctest.RandString(10)
	totp := acctest.RandStringFromCharSet(6, "0123456789")
	token := uuid.New().String()
	c, s := testClientServer(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("TestClient_Login: %v", err)
		}
		if string(body) != "" {
			t.Errorf("TestClient_Login: unexpected body: got '%s', expected nothing", body)
		}
		if err := r.ParseForm(); err != nil {
			t.Errorf("TestClient_Login: %v", err)
		}
		fu := r.Form.Get("username")
		if fu != username {
			t.Errorf("TestClient_Login: unexpected 'username' value: got '%s', expected '%s'", fu, username)
		}
		fp := r.Form.Get("password")
		if fp != password {
			t.Errorf("TestClient_Login: unexpected 'password' value: got '%s', expected '%s'", fp, password)
		}
		ft := r.Form.Get("oneTimePassword")
		if ft != totp {
			t.Errorf("TestClient_Login: unexpected 'oneTimePassword' value: got '%s', expected '%s'", ft, totp)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"data":{"token":"%s"}}`, token)
	})
	defer s.Close()
	if err := c.Login(username, password, totp); err != nil {
		t.Errorf("TestClient_Login: %v", err)
	}
	if c.Token != token {
		t.Errorf("TestClient_Login: unexpected token received: got '%s', expected '%s'", c.Token, token)
	}
}

func TestClient_Logout(t *testing.T) {
	token := uuid.New().String()
	c, s := testClientServer(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Errorf("TestClient_Logout: %v", err)
		}
		if string(body) != "" {
			t.Errorf("TestClient_Logout: unexpected body: got '%s', expected nothing", body)
		}
		ft := r.Header.Get("X-Auth-Token")
		if ft != token {
			t.Errorf("TestClient_Logout: unexpected 'token' value: got '%s', expected '%s'", ft, token)
		}
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{}`)
	})
	c.Token = token
	defer s.Close()
	if err := c.Logout(); err != nil {
		t.Errorf("TestClient_Logout: %v", err)
	}
}
