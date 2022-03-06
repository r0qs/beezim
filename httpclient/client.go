// Copyright 2020 Ethersphere. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// This code is based on the beekeeper beeclient api

package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const contentType = "application/json; charset=utf-8"

type Client struct {
	Host       string
	HTTPClient *http.Client
}

type ClientOptions struct {
	HTTPClient *http.Client
}

func NewClient(u *url.URL, o *ClientOptions) (c *Client, err error) {
	c = &Client{}

	if o == nil {
		o = new(ClientOptions)
	}
	if o.HTTPClient == nil {
		o.HTTPClient = new(http.Client)
	}
	c.HTTPClient = httpClientWithTransport(u, o.HTTPClient)
	c.Host = u.Host

	return c, nil
}

func httpClientWithTransport(baseURL *url.URL, httpc *http.Client) *http.Client {
	if httpc == nil {
		httpc = new(http.Client)
	}

	transport := httpc.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	httpc.Transport = roundTripper(baseURL, transport)

	return httpc
}

func roundTripper(baseURL *url.URL, transport http.RoundTripper) http.RoundTripper {
	if !strings.HasSuffix(baseURL.Path, "/") {
		baseURL.Path += "/"
	}

	return roundTripperFunc(func(r *http.Request) (resp *http.Response, err error) {
		u, err := baseURL.Parse(r.URL.String())
		if err != nil {
			return nil, err
		}
		r.URL = u
		return transport.RoundTrip(r)
	})
}

// requestJSON handles the HTTP request response cycle. It JSON encodes the request
// body, creates an HTTP request with provided method on a path with required
// headers and decodes request body if the v argument is not nil and content type is
// application/json.
func (c *Client) RequestJSON(ctx context.Context, method, path string, body, v interface{}) (err error) {
	var bodyBuffer io.ReadWriter
	if body != nil {
		bodyBuffer = new(bytes.Buffer)
		if err = encodeJSON(bodyBuffer, body); err != nil {
			return err
		}
	}

	return c.Request(ctx, method, path, bodyBuffer, v)
}

// request handles the HTTP request response cycle.
func (c *Client) Request(ctx context.Context, method, path string, body io.Reader, v interface{}) (err error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Accept", contentType)
	req.Header.Set("User-Agent", "Mozilla/5.0 Firefox/86.0")

	r, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer drain(r.Body)

	if err = responseErrorHandler(r); err != nil {
		return err
	}

	if v != nil && strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return json.NewDecoder(r.Body).Decode(&v)
	}

	return nil
}

func (c *Client) RequestWithResponseHeader(ctx context.Context, method, path string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Accept", contentType)
	req.Header.Set("User-Agent", "Mozilla/5.0 Firefox/86.0")

	r, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer drain(r.Body)

	if err = responseErrorHandler(r); err != nil {
		return nil, err
	}

	return r, nil
}

// encodeJSON writes a JSON-encoded v object to the provided writer with
// SetEscapeHTML set to false.
func encodeJSON(w io.Writer, v interface{}) (err error) {
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	return enc.Encode(v)
}

// requestData handles the HTTP request response cycle.
func (c *Client) RequestData(ctx context.Context, method, path string, body io.Reader) (resp io.ReadCloser, err error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)

	if body != nil {
		req.Header.Set("Content-Type", contentType)
	}
	req.Header.Set("Accept", contentType)
	req.Header.Set("User-Agent", "Mozilla/5.0 Firefox/86.0")

	r, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if err = responseErrorHandler(r); err != nil {
		return nil, err
	}

	return r.Body, nil
}

// requestWithHeader handles the HTTP request response cycle.
func (c *Client) RequestWithHeader(ctx context.Context, method, path string, header http.Header, body io.Reader, v interface{}) (err error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)

	req.Header = header
	req.Header.Add("Accept", contentType)
	req.Header.Set("User-Agent", "Mozilla/5.0 Firefox/86.0")

	r, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	if err = responseErrorHandler(r); err != nil {
		return err
	}

	if v != nil && strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return json.NewDecoder(r.Body).Decode(&v)
	}

	return
}

// drain discards all of the remaining data from the reader and closes it,
// asynchronously.
func drain(r io.ReadCloser) {
	go func() {
		// Panicking here does not put data in
		// an inconsistent state.
		defer func() {
			_ = recover()
		}()

		_, _ = io.Copy(ioutil.Discard, r)
		r.Close()
	}()
}

// responseErrorHandler returns an error based on the HTTP status code or nil if
// the status code is from 200 to 299.
func responseErrorHandler(r *http.Response) (err error) {
	if r.StatusCode == http.StatusAccepted {
		return ErrRecoveryInitiated
	}
	if r.StatusCode/100 == 2 {
		return nil
	}
	switch r.StatusCode {
	case http.StatusBadRequest:
		return decodeBadRequest(r)
	case http.StatusUnauthorized:
		return ErrUnauthorized
	case http.StatusForbidden:
		return ErrForbidden
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusTooManyRequests:
		return ErrTooManyRequests
	case http.StatusInternalServerError:
		return ErrInternalServerError
	case http.StatusServiceUnavailable:
		return ErrServiceUnavailable
	default:
		return errors.New(strings.ToLower(r.Status))
	}
}

// decodeBadRequest parses the body of HTTP response that contains a list of
// errors as the result of bad request data.
func decodeBadRequest(r *http.Response) (err error) {
	type badRequestResponse struct {
		Errors []string `json:"errors"`
	}

	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		return NewBadRequestError("bad request")
	}
	var e badRequestResponse
	if err = json.NewDecoder(r.Body).Decode(&e); err != nil {
		if err == io.EOF {
			return NewBadRequestError("bad request")
		}
		return err
	}
	return NewBadRequestError(e.Errors...)
}

// roundTripperFunc type is an adapter to allow the use of ordinary functions as
// http.RoundTripper interfaces. If f is a function with the appropriate
// signature, roundTripperFunc(f) is a http.RoundTripper that calls f.
type roundTripperFunc func(*http.Request) (*http.Response, error)

// RoundTrip calls f(r).
func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}
