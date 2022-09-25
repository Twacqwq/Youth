package pkg

import (
	"io"
	"net/http"
	"strings"
	"time"
)

type requestOption struct {
	timeout time.Duration
	data    string
	headers map[string]string
}

type Option interface {
	apply(option *requestOption) error
}

type optionFunc func(opts *requestOption) error

func (f optionFunc) apply(opts *requestOption) error {
	return f(opts)
}

func httpRequest(method, url string, options ...Option) (code int, details string, err error) {
	reqOpt := defaultOptions()
	for _, option := range options {
		option.apply(reqOpt)
	}

	req, err := http.NewRequest(method, url, strings.NewReader(reqOpt.data))
	if err != nil {
		return
	}
	defer req.Body.Close()

	if len(reqOpt.headers) != 0 {
		for k, v := range reqOpt.headers {
			req.Header.Add(k, v)
		}
	}
	client := http.Client{Timeout: reqOpt.timeout}
	res, err := client.Do(req)
	if err != nil {
		return
	}

	defer res.Body.Close()
	code = res.StatusCode
	result, _ := io.ReadAll(res.Body)
	details = string(result)

	return
}

func Get(url string, options ...Option) (code int, details string, err error) {
	return httpRequest("GET", url, options...)
}

func Post(url string, options ...Option) (code int, details string, err error) {
	return httpRequest("POST", url, options...)
}

// default option
func defaultOptions() *requestOption {
	return &requestOption{
		timeout: 5 * time.Second,
		data:    "",
		headers: make(map[string]string),
	}
}

func WithTimeout(timeout time.Duration) Option {
	return optionFunc(func(opts *requestOption) (err error) {
		opts.timeout, err = timeout, nil
		return
	})
}

func WithHeaders(headers map[string]string) Option {
	return optionFunc(func(opts *requestOption) (err error) {
		for k, v := range headers {
			opts.headers[k] = v
		}
		return
	})
}

func WithData(data string) Option {
	return optionFunc(func(opts *requestOption) (err error) {
		opts.data, err = data, nil
		return
	})
}
