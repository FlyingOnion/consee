// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package repo

import (
	"context"

	"github.com/FlyingOnion/consee/backend/consul"
)

type KVRepo interface {
	ListKeys(ctx context.Context, prefix, sep string) (*consul.Response[[]string], error)
	List(ctx context.Context, prefix string) (*consul.Response[[]*consul.KVPair], error)
	Read(ctx context.Context, key string) (*consul.Response[*consul.KVPair], error)
	Write(ctx context.Context, key, value string) (*consul.Response[bool], error)
	Delete(ctx context.Context, key string) (*consul.Response[bool], error)
	WatchKeys(ctx context.Context, prefix string, onResponse func(*consul.Response[[]string], error) (stop bool))
}
