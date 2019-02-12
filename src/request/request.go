package request

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

// M key value map
type M map[string]interface{}

// Option request option
type Option struct {
	Method  string
	URL     string
	Headers M
	QS      M
	Form    M
	Body    interface{}
	JSON    bool
	Timeout time.Duration
}

// Request request
func Request(opt Option) (*http.Response, error) {
	var err error

	if opt.Method == "" {
		opt.Method = http.MethodGet // default method: get
	}

	if opt.Timeout.Nanoseconds() == 0 {
		opt.Timeout = time.Second * 30 // default timeout: 30s
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Timeout:   opt.Timeout,
		Transport: tr,
	}

	header := http.Header{}
	for key, value := range opt.Headers {
		v, ok := value.(string)
		if ok {
			header.Set(key, v)
		}
	}

	var bs []byte

	if opt.Body != nil {
		bs, err = json.Marshal(opt.Body)
		if err != nil {
			return nil, err
		}
	} else if opt.Form != nil {
		header.Set("Content-Type", "application/x-www-form-urlencoded")
		form := url.Values{}
		for key, value := range opt.Form {
			if v, ok := value.(string); ok {
				form.Set(key, v)
			} else if vs, ok := value.([]string); ok {
				for _, v := range vs {
					form.Add(key, v)
				}
			}

			bs = []byte(form.Encode())
		}
	}

	if opt.JSON {
		header.Set("Content-Type", "application/json")
	}

	header.Set("User-Agent", UserAgent) // user agent

	req, err := http.NewRequest(opt.Method, opt.URL, bytes.NewReader(bs))
	if err != nil {
		return nil, err
	}

	req.Header = header
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
