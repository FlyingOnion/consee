// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package infra

import (
	"context"

	"github.com/FlyingOnion/consee/backend/consul"
)

type admin struct {
	q      *consul.QueryOptions
	w      *consul.WriteOptions
	client *consul.Client
}

func (a *admin) ListKeys(ctx context.Context, prefix, sep string) (*consul.Response[[]string], error) {
	return a.client.KV().Keys(ctx, prefix, sep, a.q)
}

func (a *admin) List(ctx context.Context, prefix string) (*consul.Response[[]*consul.KVPair], error) {
	return a.client.KV().List(ctx, prefix, a.q)
}

func (a *admin) Read(ctx context.Context, key string) (*consul.Response[*consul.KVPair], error) {
	return a.client.KV().Get(ctx, key, a.q)
}

func (a *admin) Write(ctx context.Context, key, value string) (*consul.Response[bool], error) {
	return a.client.KV().Put(ctx, &consul.KVPair{Key: key, Value: []byte(value)}, a.w)
}

func (a *admin) Delete(ctx context.Context, key string) (*consul.Response[bool], error) {
	if key[len(key)-1] == '/' {
		return a.client.KV().DeleteTree(ctx, key, consul.WriteOptionsFromContext(ctx))
	}
	return a.client.KV().Delete(ctx, key, a.w)
}

func (a *admin) WatchKeys(ctx context.Context, prefix string, onResponse func(*consul.Response[[]string], error) (stop bool)) {
	a.client.KV().WatchKeys(ctx, prefix, a.q, onResponse)
}

func (a *admin) ReadSelf(ctx context.Context) (*consul.Response[*consul.ACLToken], error) {
	return a.client.ACL().TokenReadSelf(ctx, a.q)
}

func (a *admin) ListTokens(ctx context.Context) (*consul.Response[[]*consul.ACLToken], error) {
	return a.client.ACL().TokenList(ctx, a.q)
}

func (a *admin) ListTokensFiltered(ctx context.Context, t consul.ACLTokenFilterOptions) (*consul.Response[[]*consul.ACLToken], error) {
	return a.client.ACL().TokenListFiltered(ctx, a.q, t)
}

func (a *admin) ReadToken(ctx context.Context, id string) (*consul.Response[*consul.ACLToken], error) {
	return a.client.ACL().TokenRead(ctx, id, a.q)
}

func (a *admin) CreateToken(ctx context.Context, req *consul.ACLToken) (*consul.Response[*consul.ACLToken], error) {
	return a.client.ACL().TokenCreate(ctx, req, a.w)
}

func (a *admin) UpdateToken(ctx context.Context, req *consul.ACLToken) (*consul.Response[*consul.ACLToken], error) {
	return a.client.ACL().TokenUpdate(ctx, req, a.w)
}

func (a *admin) DeleteToken(ctx context.Context, id string) (*consul.Response[bool], error) {
	return a.client.ACL().TokenDelete(ctx, id, a.w)
}

func (a *admin) ListPolicies(ctx context.Context) (*consul.Response[[]*consul.ACLPolicy], error) {
	return a.client.ACL().PolicyList(ctx, a.q)
}

func (a *admin) ReadPolicy(ctx context.Context, id string) (*consul.Response[*consul.ACLPolicy], error) {
	return a.client.ACL().PolicyRead(ctx, id, a.q)
}

func (a *admin) ReadPolicyByName(ctx context.Context, name string) (*consul.Response[*consul.ACLPolicy], error) {
	return a.client.ACL().PolicyReadByName(ctx, name, a.q)
}

func (a *admin) CreatePolicy(ctx context.Context, req *consul.ACLPolicy) (*consul.Response[*consul.ACLPolicy], error) {
	return a.client.ACL().PolicyCreate(ctx, req, a.w)
}

func (a *admin) UpdatePolicy(ctx context.Context, req *consul.ACLPolicy) (*consul.Response[*consul.ACLPolicy], error) {
	return a.client.ACL().PolicyUpdate(ctx, req, a.w)
}

func (a *admin) DeletePolicy(ctx context.Context, id string) (*consul.Response[bool], error) {
	return a.client.ACL().PolicyDelete(ctx, id, a.w)
}
