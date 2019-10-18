package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
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
		fmt.Fprint(w, `{}`)
	})
	c.Token = token
	defer s.Close()
	if err := c.Logout(); err != nil {
		t.Errorf("TestClient_Logout: %v", err)
	}
}

func TestClient_responseDataToError(t *testing.T) {
	testCases := []struct {
		p string
		e string
	}{
		{
			p: `{"message":"foo"}`,
			e: `megaport-api: foo`,
		},
		{
			p: `{"message":"foo","data":"bar"}`,
			e: `megaport-api: foo: bar`,
		},
		{
			p: `{"message":"foo","data":["bar","baz"]}`,
			e: `megaport-api: foo: 2 errors: ['bar', 'baz']`,
		},
		{
			p: `{"message":"foo","data":{"a":"b","c":5}}`,
			e: `megaport-api: foo: a="b" c=5`,
		},
		{
			p: `{"message":"foo","data":true}`,
			e: `megaport-api: foo: cannot process error data of type bool: true`,
		},
	}
	c, s := testClientServer(func(w http.ResponseWriter, r *http.Request) {
		tc, err := strconv.Atoi(strings.Trim(r.URL.Path, `/`))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, `{}`)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, testCases[tc].p)
	})
	req, err := http.NewRequest(http.MethodGet, s.URL, nil)
	if err != nil {
		t.Errorf("TestClient_responseDataToError: %v", err)
	}
	e := `megaport-api: not found`
	if err := c.do(req, nil); err.Error() != e {
		t.Errorf("TestClient_responseDataToError: unexpected error:\n\tgot     : %v\n\texpected: %s", err, e)
	}
	for i, tc := range testCases {
		req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%d", s.URL, i), nil)
		if err != nil {
			t.Errorf("TestClient_responseDataToError: %v", err)
		}
		if err := c.do(req, nil); err.Error() != tc.e {
			t.Errorf("TestClient_responseDataToError: unexpected error in test case #%d:\n\tgot     : %v\n\texpected: %s", i, err, tc.e)
		}
	}
}
