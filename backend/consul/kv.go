// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package consul

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// KVPair is used to represent a single K/V entry
type KVPair struct {
	// Key is the name of the key. It is also part of the URL path when accessed
	// via the API.
	Key string

	// CreateIndex holds the index corresponding the creation of this KVPair. This
	// is a read-only field.
	CreateIndex uint64

	// ModifyIndex is used for the Check-And-Set operations and can also be fed
	// back into the WaitIndex of the QueryOptions in order to perform blocking
	// queries.
	ModifyIndex uint64

	// LockIndex holds the index corresponding to a lock on this key, if any. This
	// is a read-only field.
	LockIndex uint64

	// Flags are any user-defined flags on the key. It is up to the implementer
	// to check these values, since Consul does not treat them specially.
	Flags uint64

	// Value is the value for the key. This can be any value, but it will be
	// base64 encoded upon transport.
	Value []byte

	// Session is a string representing the ID of the session. Any other
	// interactions with this key over the same session must specify the same
	// session ID.
	Session string

	// Namespace is the namespace the KVPair is associated with
	// Namespacing is a Consul Enterprise feature.
	Namespace string `json:",omitempty"`

	// Partition is the partition the KVPair is associated with
	// Admin Partition is a Consul Enterprise feature.
	Partition string `json:",omitempty"`
}

type Response[T any] struct {
	Status   int
	Duration time.Duration
	*Metadata
	RawBody []byte
	Body    T
	// Err describes any error that occurred during the parse of the response.
	// This is not the same as the error returned by the HTTP request.
	// No matter whether Err is nil or not, response from consul is returned.
	Err error
	// Decode  func(T) error
}

func (c *Client) KV() *KV {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.kv == nil {
		c.kv = &KV{c}
	}
	return c.kv
}

type KV struct {
	c *Client
}

func (kv *KV) Keys(ctx context.Context, prefix, sep string, q *QueryOptions) (*Response[[]string], error) {
	options := append(
		q.toRequestOptions(),
		reqWithQuery("keys", ""),
		reqWithQuery("separator", sep),
	)
	httpRequest := kv.c.newRequest(ctx, http.MethodGet, "/v1/kv/"+strings.TrimPrefix(prefix, "/"), options...)
	return responseDirectly(kv.c.httpClient, httpRequest, decodeStringSlice)
	// t := time.Now()
	// httpResponse, err := kv.c.httpClient.Do(httpRequest)
	// if httpResponse != nil && httpResponse.Body != nil {
	// 	defer httpResponse.Body.Close()
	// }
	// if err != nil {
	// 	return nil, err
	// }
	// duration := time.Since(t)
	// return newResponse[[]string](httpResponse, duration, func(b []byte) ([]string, error) {
	// 	keys := []string{}
	// 	json.Unmarshal(b, &keys)
	// 	return keys, nil
	// })

	// resp := &Response[[]string]{
	// 	Status:   httpResponse.StatusCode,
	// 	Duration: time.Since(t),
	// 	Metadata: &Metadata{},
	// }
	// err = http2.ParseHeader(httpResponse.Header, resp.Metadata)
	// resp.Err = err
	// // parse body anyway
	// b := make([]byte, httpResponse.ContentLength)
	// _, err = httpResponse.Body.Read(b)
	// resp.RawBody = b
	// if resp.Err != nil {
	// 	return resp, nil
	// }
	// if err != nil && err != io.EOF {
	// 	resp.Err = err
	// 	return resp, nil
	// }
	// if httpResponse.StatusCode != http.StatusOK {
	// 	return resp, nil
	// }
	// keys := []string{}
	// err = json.Unmarshal(b, &keys)
	// if err != nil {
	// 	resp.Err = err
	// 	return resp, nil
	// }
	// resp.Body = keys
	// return resp, nil
}

func (kv *KV) Get(ctx context.Context, key string, q *QueryOptions) (*Response[*KVPair], error) {
	options := q.toRequestOptions()
	httpRequest := kv.c.newRequest(ctx, http.MethodGet, "/v1/kv/"+strings.TrimPrefix(key, "/"), options...)
	return responseDirectly(kv.c.httpClient, httpRequest, decodeKVPair)
	// t := time.Now()
	// httpResponse, err := kv.c.httpClient.Do(httpRequest)
	// if httpResponse != nil && httpResponse.Body != nil {
	// 	defer httpResponse.Body.Close()
	// }
	// if err != nil {
	// 	return nil, err
	// }
	// duration := time.Since(t)
	// resp := &Response[*KVPair]{
	// 	Status:   httpResponse.StatusCode,
	// 	Duration: time.Since(t),
	// 	Metadata: &Metadata{},
	// }
	// err = http2.ParseHeader(httpResponse.Header, resp.Metadata)
	// resp.Err = err
	// // parse body anyway
	// b := make([]byte, httpResponse.ContentLength)
	// _, err = httpResponse.Body.Read(b)
	// resp.RawBody = b
	// if resp.Err != nil {
	// 	return resp, nil
	// }
	// if err != nil && err != io.EOF {
	// 	resp.Err = err
	// 	return resp, nil
	// }
	// if httpResponse.StatusCode != http.StatusOK {
	// 	return resp, nil
	// }
	// kvPairs := []*KVPair{}
	// err = json.Unmarshal(b, &kvPairs)
	// if err != nil {
	// 	resp.Err = err
	// 	return resp, nil
	// }
	// resp.Body = kvPairs[0]
	// return newResponse(httpResponse, duration, func(b []byte) (*KVPair, error) {
	// 	kvPairs := []*KVPair{}
	// 	e := json.Unmarshal(b, &kvPairs)
	// 	if e != nil {
	// 		return nil, e
	// 	}
	// 	return kvPairs[0], nil
	// })
}

func (kv *KV) Put(ctx context.Context, kvPair *KVPair, w *WriteOptions) (*Response[bool], error) {
	key := kvPair.Key
	if len(key) > 0 && key[0] == '/' {
		return nil, fmt.Errorf("Invalid key. Key must not begin with a '/': %s", key)
	}
	options := append(w.toRequestOptions(), reqWithBody(kvPair.Value), reqWithContentType("application/octet-stream"))
	httpRequest := kv.c.newRequest(ctx, http.MethodPut, "/v1/kv/"+key, options...)
	return responseDirectly(kv.c.httpClient, httpRequest, decodeTrue)
	// t := time.Now()
	// httpResponse, err := kv.c.httpClient.Do(httpRequest)
	// if httpResponse != nil && httpResponse.Body != nil {
	// 	defer httpResponse.Body.Close()
	// }
	// if err != nil {
	// 	return nil, err
	// }
	// duration := time.Since(t)
	// return newResponse(httpResponse, duration, func(b []byte) (bool, error) {
	// 	return string(b) == "true", nil
	// })
}

func (kv *KV) Delete(ctx context.Context, key string, w *WriteOptions) (*Response[bool], error) {
	httpRequest := kv.c.newRequest(ctx, http.MethodDelete, "/v1/kv/"+strings.TrimPrefix(key, "/"), w.toRequestOptions()...)
	return responseDirectly(kv.c.httpClient, httpRequest, decodeTrue)
}

func (kv *KV) DeleteTree(ctx context.Context, prefix string, w *WriteOptions) (*Response[bool], error) {
	if prefix[len(prefix)-1] != '/' {
		// consul allow a key and a prefix to have the same name, like "foo" and "foo/";
		// without trailing slash the prefix is a single key instead of a path;
		// to avoid ambiguity, just delete the key itself
		return kv.Delete(ctx, prefix, w)
	}
	options := append(w.toRequestOptions(), reqWithQuery("recurse", ""))
	httpRequest := kv.c.newRequest(ctx, http.MethodDelete, "/v1/kv/"+strings.TrimPrefix(prefix, "/"), options...)
	return responseDirectly(kv.c.httpClient, httpRequest, decodeTrue)
}

func (kv *KV) List(ctx context.Context, prefix string, q *QueryOptions) (*Response[[]*KVPair], error) {
	options := append(q.toRequestOptions(), reqWithQuery("recurse", ""))
	httpRequest := kv.c.newRequest(ctx, http.MethodGet, "/v1/kv/"+strings.TrimPrefix(prefix, "/"), options...)
	return responseDirectly(kv.c.httpClient, httpRequest, decodeKVPairs)
}

func (kv *KV) WatchKeys(ctx context.Context, prefix string, q *QueryOptions, onResponse func(*Response[[]string], error) (stop bool)) {
	q1 := q.Copy()
	resp, err := kv.c.KV().Keys(ctx, prefix, "", q1)
	if onResponse == nil {
		onResponse = func(_ *Response[[]string], _ error) (stop bool) { return false }
	}
	stop := onResponse(resp, err)

	if resp == nil {
		slog.Info("initial response is nil", "err", err)
	} else {
		slog.Info("initial response", "status", resp.Status, "nkeys", len(resp.Body))
		for !stop && context.Cause(ctx) == nil {
			q1.WaitIndex = resp.Metadata.LastIndex
			resp, err = kv.c.KV().Keys(ctx, prefix, "", q1)
			stop = onResponse(resp, err)
		}
	}
	slog.Info("watch stopped", "context_status", context.Cause(ctx))
}
