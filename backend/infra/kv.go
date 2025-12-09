// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package infra

import (
	"context"

	"github.com/FlyingOnion/consee/backend/consul"
)

type kv struct {
	client *consul.Client
}

func (kv *kv) ListKeys(ctx context.Context, prefix, sep string) (*consul.Response[[]string], error) {
	return kv.client.KV().Keys(ctx, prefix, sep, consul.QueryOptionsFromContext(ctx))
}

func (kv *kv) List(ctx context.Context, prefix string) (*consul.Response[[]*consul.KVPair], error) {
	return kv.client.KV().List(ctx, prefix, consul.QueryOptionsFromContext(ctx))
}

func (kv *kv) Read(ctx context.Context, key string) (*consul.Response[*consul.KVPair], error) {
	return kv.client.KV().Get(ctx, key, consul.QueryOptionsFromContext(ctx))
}

func (kv *kv) Write(ctx context.Context, key, value string) (*consul.Response[bool], error) {
	return kv.client.KV().Put(ctx, &consul.KVPair{Key: key, Value: []byte(value)}, consul.WriteOptionsFromContext(ctx))
}

func (kv *kv) Delete(ctx context.Context, key string) (*consul.Response[bool], error) {
	if key[len(key)-1] == '/' {
		return kv.client.KV().DeleteTree(ctx, key, consul.WriteOptionsFromContext(ctx))
	}
	return kv.client.KV().Delete(ctx, key, consul.WriteOptionsFromContext(ctx))
}

func (kv *kv) WatchKeys(ctx context.Context, prefix string, onResponse func(*consul.Response[[]string], error) (stop bool)) {
	kv.client.KV().WatchKeys(ctx, prefix, consul.QueryOptionsFromContext(ctx), onResponse)
}
