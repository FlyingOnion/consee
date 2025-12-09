// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package consul

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	http2 "github.com/FlyingOnion/consee/backend/http"
)

type Client struct {
	addr       string
	scheme     string
	prefix     string
	httpClient *http.Client

	mu  sync.Mutex
	kv  *KV
	acl *ACL
}

func NewClient(options ...ClientOption) *Client {
	client := &Client{
		addr:       "localhost:8500",
		scheme:     "http",
		httpClient: http.DefaultClient,
	}
	for _, op := range options {
		op(client)
	}
	return client
}

type ClientOption func(*Client)

func WithAddress(addr string) ClientOption {
	return func(c *Client) {
		c.addr = addr
	}
}

func WithHTTPS() ClientOption {
	return func(c *Client) { c.scheme = "https" }
}

func WithPrefix(prefix string) ClientOption {
	return func(c *Client) { c.prefix = prefix }
}

type request struct {
	header http.Header
	body   []byte
	query  url.Values
}

type requestOption func(*request)

func requestOptionNoop(*request) {}

func reqWithBody(v []byte) requestOption {
	return func(r *request) {
		r.body = make([]byte, len(v))
		copy(r.body, v)
	}
}

func reqWithQuery(k, v string) requestOption {
	return func(r *request) { r.query.Add(k, v) }
}

func reqWithToken(token string) requestOption {
	if token == "" {
		return requestOptionNoop
	}
	return func(r *request) { r.header.Set("X-Consul-Token", token) }
}

func reqWithDataCenter(dc string) requestOption {
	if dc == "" {
		return requestOptionNoop
	}
	return reqWithQuery("dc", dc)
}

func reqWithNamespace(namespace string) requestOption {
	if namespace == "" {
		return requestOptionNoop
	}
	return reqWithQuery("ns", namespace)
}

func reqWithPartition(partition string) requestOption {
	if partition == "" {
		return requestOptionNoop
	}
	return reqWithQuery("partition", partition)
}

func reqWithSamenessGroup(sg string) requestOption {
	if sg == "" {
		return requestOptionNoop
	}
	return reqWithQuery("sameness-group", sg)
}

func reqWithPeer(peer string) requestOption {
	if peer == "" {
		return requestOptionNoop
	}
	return reqWithQuery("peer", peer)
}

func reqWithStale(stale bool) requestOption {
	if !stale {
		return requestOptionNoop
	}
	return reqWithQuery("stale", "")
}

func reqWithConsistent(requireConsistent bool) requestOption {
	if !requireConsistent {
		return requestOptionNoop
	}
	return reqWithQuery("consistent", "")
}

func reqWithIndex(index uint64) requestOption {
	if index == 0 {
		return requestOptionNoop
	}
	return reqWithQuery("index", strconv.FormatUint(index, 10))
}

func durToMsec(dur time.Duration) string {
	ms := dur / time.Millisecond
	if dur > 0 && ms == 0 {
		ms = 1
	}
	return strconv.FormatInt(int64(ms), 10) + "ms"
}

func reqWithWaitTime(duration time.Duration) requestOption {
	if duration == 0 {
		return requestOptionNoop
	}
	return reqWithQuery("wait", durToMsec(duration))
}

func reqWithWaitHash(hash string) requestOption {
	if hash == "" {
		return requestOptionNoop
	}
	return reqWithQuery("hash", hash)
}

func reqWithNear(near string) requestOption {
	if near == "" {
		return requestOptionNoop
	}
	return reqWithQuery("near", near)
}

func reqWithFilter(filter string) requestOption {
	if filter == "" {
		return requestOptionNoop
	}
	return reqWithQuery("filter", filter)
}

func reqWithNodeMeta(meta map[string]string) requestOption {
	if len(meta) == 0 {
		return requestOptionNoop
	}
	return func(r *request) {
		for key, value := range meta {
			r.query.Add("node-meta", key+":"+value)
		}
	}
}

func reqWithRelayFactor(rf uint8) requestOption {
	if rf == 0 {
		return requestOptionNoop
	}
	return reqWithQuery("relay-factor", strconv.FormatUint(uint64(rf), 10))
}

func reqWithLocalOnly(b bool) requestOption {
	if !b {
		return requestOptionNoop
	}
	return reqWithQuery("local-only", "true")
}

func reqWithConnect(c bool) requestOption {
	if !c {
		return requestOptionNoop
	}
	return reqWithQuery("connect", "true")
}

func reqWithCache(useCache, requireConsistent bool, maxAge, staleIfError time.Duration) requestOption {
	if !useCache || requireConsistent {
		return requestOptionNoop
	}
	return func(r *request) {
		r.query.Add("cache", "")
		cc := make([]string, 0, 2)
		if maxAge > 0 {
			cc = append(cc, "max-age="+strconv.FormatInt(int64(maxAge.Seconds()), 10))
		}
		if staleIfError > 0 {
			cc = append(cc, "stale-if-error="+strconv.FormatInt(int64(staleIfError.Seconds()), 10))
		}
		r.header.Add("Cache-Control", strings.Join(cc, ", "))
	}
}

func reqWithMergeCentralConfig(b bool) requestOption {
	if !b {
		return requestOptionNoop
	}
	return reqWithQuery("merge_central_config", "")
}

func reqWithGlobal(b bool) requestOption {
	if !b {
		return requestOptionNoop
	}
	return reqWithQuery("global", "")
}

// func reqWithRecurse() requestOption {
// 	return reqWithQuery("recurse", "")
// }

// func reqWithKeys() requestOption {
// 	return reqWithQuery("keys", "")
// }

func reqWithContentType(contentType string) requestOption {
	return func(r *request) { r.header.Set("Content-Type", contentType) }
}

func noBody() (io.ReadCloser, error) { return http.NoBody, nil }

func (c *Client) newRequest(ctx context.Context, method, path string, options ...requestOption) *http.Request {
	req := &request{header: http.Header{}, query: url.Values{}}
	for _, opt := range options {
		opt(req)
	}
	header := req.header.Clone()
	var (
		body    io.ReadCloser                 = http.NoBody
		getBody func() (io.ReadCloser, error) = noBody
	)
	if len(req.body) > 0 {
		body = io.NopCloser(bytes.NewReader(req.body))
		getBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(req.body)), nil
		}
		if header.Get("Content-Type") == "" {
			header.Set("Content-Type", "application/json")
		}
	}
	httpReq := &http.Request{
		Header:        header,
		Body:          body,
		GetBody:       getBody,
		ContentLength: int64(len(req.body)),
		Method:        method,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		URL: &url.URL{
			Host:     c.addr,
			Path:     c.prefix + path,
			Scheme:   c.scheme,
			RawQuery: req.query.Encode(),
		},
	}
	switch ctx {
	case nil:
	case context.Background():
	case context.TODO():
	default:
		httpReq = httpReq.WithContext(ctx)
	}
	return httpReq
}

func responseDirectly[T any](client *http.Client, httpReq *http.Request, decode func([]byte) (T, error)) (*Response[T], error) {
	t := time.Now()
	httpResponse, err := client.Do(httpReq)
	if httpResponse != nil && httpResponse.Body != nil {
		defer httpResponse.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	resp := &Response[T]{
		Duration: time.Since(t),
		Status:   httpResponse.StatusCode,
	}

	if httpResponse.Request.Method == http.MethodGet {
		resp.Metadata = &Metadata{}
		err := http2.ParseHeader(httpResponse.Header, resp.Metadata)
		if err != nil {
			resp.Err = err
			return resp, nil
		}
	}

	if httpResponse.ContentLength > 0 {
		// parse body anyway
		b := make([]byte, httpResponse.ContentLength)
		_, err := httpResponse.Body.Read(b)
		// there may be partial data in results
		resp.RawBody = b
		if err != nil && err != io.EOF {
			resp.Err = err
			return resp, nil
		}
		if httpResponse.StatusCode == http.StatusOK {
			resp.Body, resp.Err = decode(b)
		}
		return resp, nil
	}
	buf := make([]byte, 1024)
	for {
		n, err := httpResponse.Body.Read(buf)
		if n > 0 {
			resp.RawBody = append(resp.RawBody, buf[:n]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			resp.Err = err
			return resp, nil
		}
	}
	if len(resp.RawBody) > 0 {
		resp.Body, resp.Err = decode(resp.RawBody)
	}
	return resp, nil
}

// func newResponse[T any](httpResponse *http.Response, duration time.Duration, decode func([]byte) (T, error)) (*Response[T], error) {
// 	resp := &Response[T]{
// 		Duration: duration,
// 		Status:   httpResponse.StatusCode,
// 	}

// 	if httpResponse.Request.Method == http.MethodGet {
// 		resp.Metadata = &Metadata{}
// 		err := http2.ParseHeader(httpResponse.Header, resp.Metadata)
// 		if err != nil {
// 			resp.Err = err
// 			return resp, nil
// 		}
// 	}

// 	if httpResponse.ContentLength > 0 {
// 		// parse body anyway
// 		b := make([]byte, httpResponse.ContentLength)
// 		_, err := httpResponse.Body.Read(b)
// 		// there may be partial data in results
// 		resp.RawBody = b
// 		if err != nil && err != io.EOF {
// 			resp.Err = err
// 			return resp, nil
// 		}
// 		if httpResponse.StatusCode == http.StatusOK {
// 			resp.Body, resp.Err = decode(b)
// 		}
// 	}
// 	return resp, nil

// }
