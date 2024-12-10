package esclient

import (
	"encoding/json"
	"io"
	"net/http"
)

type Client interface {
	Ping
	Index
	Bulk
	Search
	Count
}

type client struct {
	httpClient *http.Client
	baseUrl    string
	username   string
	password   string
}

type ClientOption func(*client)

func WithHttpClient(httpClient *http.Client) ClientOption {
	return func(c *client) {
		c.httpClient = httpClient
	}
}

func WithBasicAuth(username, password string) ClientOption {
	return func(c *client) {
		c.username = username
		c.password = password
	}
}

func NewClient(baseUrl string, options ...ClientOption) Client {
	c := &client{
		baseUrl: baseUrl,
	}

	for _, option := range options {
		option(c)
	}

	if c.httpClient == nil {
		c.httpClient = &http.Client{}
	}

	return c
}

func (c *client) do(req *http.Request) (*http.Response, error) {
	if len(req.Header) == 0 {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.username != "" && c.password != "" {
		req.SetBasicAuth(c.username, c.password)
	}

	// TODO(Michael): add retry mechanism
	return c.httpClient.Do(req)
}

type Response[T any] struct {
	StatusCode   int
	ErrorMessage string
	Result       *T
}

func (r *Response[T]) SetBody(body io.ReadCloser) error {
	if body == nil {
		return nil
	}

	if r.IsError() {
		errorMessages, err := parseCloserToString(body)
		if err != nil {
			return err
		}

		r.ErrorMessage = errorMessages
	}

	var result T
	if err := json.NewDecoder(body).Decode(&result); err != nil {
		return err
	}
	r.Result = &result

	return nil
}

func (r *Response[T]) String() string {
	str, _ := json.Marshal(r.Result)
	return string(str)
}

func (r *Response[T]) IsError() bool {
	return r.StatusCode > 299
}

func parseCloserToString(r io.ReadCloser) (string, error) {
	body, err := io.ReadAll(r)
	if err != nil {
		return "", nil
	}

	var data interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}
