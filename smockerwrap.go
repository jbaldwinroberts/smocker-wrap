package smocker_wrap

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
)

const (
	HttpScheme                       = "http"
	Host                             = "localhost:8081"
	HeaderContentTypeApplicationJson = "application/json"
)

var UnexpectedStatusCodeError = errors.New("received an unexpected status code")

type MockRequest struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

type MockResponse struct {
	Status int         `json:"status"`
	Body   interface{} `json:"body,omitempty"`
}

type MockContext struct {
	Times int `json:"times,omitempty"`
}

type Mock struct {
	Request  MockRequest   `json:"request,omitempty"`
	Response *MockResponse `json:"response,omitempty"`
	Context  *MockContext  `json:"context,omitempty"`
}

type VerifyResult struct {
	Mocks struct {
		Verified bool   `json:"verified"`
		AllUsed  bool   `json:"all_used"`
	} `json:"mocks"`
}

func Reset(force bool) error {
	u := url.URL{
		Scheme: HttpScheme,
		Host:   Host,
		Path:   "/reset",
	}

	q := u.Query()
	q.Set("force", strconv.FormatBool(force))
	u.RawQuery = q.Encode()

	resp, err := http.Post(u.String(), "", nil)
	if err != nil {
		fmt.Printf("Reset request encountered an error: %+v", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Reset request received an unexpected status code: %d", resp.StatusCode)
		return UnexpectedStatusCodeError
	}

	return nil
}

func AddMock(reset bool, session string, request MockRequest, response *MockResponse) error {
	mocks := []*Mock{
		{
			Request:  request,
			Response: response,
			Context: &MockContext{Times: 1},
		},
	}

	return AddMocks(reset, session, mocks)
}

func AddMocks(reset bool, session string, mocks []*Mock) error {
	u := url.URL{
		Scheme: HttpScheme,
		Host:   Host,
		Path:   "/mocks",
	}

	q := u.Query()
	q.Set("reset", strconv.FormatBool(reset))
	q.Set("session", session)
	u.RawQuery = q.Encode()

	b, err := json.Marshal(mocks)
	if err != nil {
		fmt.Printf("Add mocks encountered an error marshaling the mocks: %+v", err)
		return err
	}

	resp, err := http.Post(u.String(), HeaderContentTypeApplicationJson, bytes.NewBuffer(b))
	if err != nil {
		fmt.Printf("Add mocks encountered an error making the request: %+v", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Add mocks received an unexpected status code: %d", resp.StatusCode)
		return UnexpectedStatusCodeError
	}

	return nil
}

func VerifyMocks(session string) (VerifyResult, error) {
	u := url.URL{
		Scheme: HttpScheme,
		Host:   Host,
		Path:   "/sessions/verify",
	}

	// q := u.Query()
	// q.Set("session", session)
	// u.RawQuery = q.Encode()

	resp, err := http.Post(u.String(), HeaderContentTypeApplicationJson, nil)
	if err != nil {
		fmt.Printf("Verify mocks encountered an error making the request: %+v", err)
		return VerifyResult{}, err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Verify mocks received an unexpected status code: %d", resp.StatusCode)
		return VerifyResult{}, UnexpectedStatusCodeError
	}

	var verifyResult VerifyResult
	if err := json.NewDecoder(resp.Body).Decode(&verifyResult); err != nil {
		fmt.Printf("Verify mocks encountered an error decoding the response: %+v", err)
		return VerifyResult{}, err
	}

	return verifyResult, nil
}
