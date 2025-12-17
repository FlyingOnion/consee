// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package repo

import (
	"context"

	"github.com/FlyingOnion/consee/backend/consul"
)

type ACLRepo interface {
	ListTokens(ctx context.Context) (*consul.Response[[]*consul.ACLToken], error)
	ListTokensFiltered(ctx context.Context, t consul.ACLTokenFilterOptions) (*consul.Response[[]*consul.ACLToken], error)
	ReadToken(ctx context.Context, id string) (*consul.Response[*consul.ACLToken], error)
	ReadSelf(ctx context.Context) (*consul.Response[*consul.ACLToken], error)
	CreateToken(ctx context.Context, req *consul.ACLToken) (*consul.Response[*consul.ACLToken], error)
	UpdateToken(ctx context.Context, req *consul.ACLToken) (*consul.Response[*consul.ACLToken], error)
	DeleteToken(ctx context.Context, id string) (*consul.Response[bool], error)

	ListPolicies(ctx context.Context) (*consul.Response[[]*consul.ACLPolicy], error)
	ReadPolicy(ctx context.Context, id string) (*consul.Response[*consul.ACLPolicy], error)
	ReadPolicyByName(ctx context.Context, name string) (*consul.Response[*consul.ACLPolicy], error)
	CreatePolicy(ctx context.Context, req *consul.ACLPolicy) (*consul.Response[*consul.ACLPolicy], error)
	UpdatePolicy(ctx context.Context, req *consul.ACLPolicy) (*consul.Response[*consul.ACLPolicy], error)
	DeletePolicy(ctx context.Context, id string) (*consul.Response[bool], error)

	ListRoles(ctx context.Context) (*consul.Response[[]*consul.ACLRole], error)
	ReadRole(ctx context.Context, id string) (*consul.Response[*consul.ACLRole], error)
	ReadRoleByName(ctx context.Context, name string) (*consul.Response[*consul.ACLRole], error)
	CreateRole(ctx context.Context, req *consul.ACLRole) (*consul.Response[*consul.ACLRole], error)
	UpdateRole(ctx context.Context, req *consul.ACLRole) (*consul.Response[*consul.ACLRole], error)
	DeleteRole(ctx context.Context, id string) (*consul.Response[bool], error)
}
