// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package infra

import (
	"context"

	"github.com/FlyingOnion/consee/backend/consul"
)

type acl struct {
	client *consul.Client
}

func (a *acl) ListTokens(ctx context.Context) (*consul.Response[[]*consul.ACLToken], error) {
	return a.client.ACL().TokenList(ctx, consul.QueryOptionsFromContext(ctx))
}

func (a *acl) ListTokensFiltered(ctx context.Context, t consul.ACLTokenFilterOptions) (*consul.Response[[]*consul.ACLToken], error) {
	return a.client.ACL().TokenListFiltered(ctx, consul.QueryOptionsFromContext(ctx), t)
}

func (a *acl) ReadToken(ctx context.Context, id string) (*consul.Response[*consul.ACLToken], error) {
	return a.client.ACL().TokenRead(ctx, id, consul.QueryOptionsFromContext(ctx))
}

func (a *acl) ReadSelf(ctx context.Context) (*consul.Response[*consul.ACLToken], error) {
	return a.client.ACL().TokenReadSelf(ctx, consul.QueryOptionsFromContext(ctx))
}

func (a *acl) CreateToken(ctx context.Context, req *consul.ACLToken) (*consul.Response[*consul.ACLToken], error) {
	return a.client.ACL().TokenCreate(ctx, req, consul.WriteOptionsFromContext(ctx))
}

func (a *acl) UpdateToken(ctx context.Context, req *consul.ACLToken) (*consul.Response[*consul.ACLToken], error) {
	return a.client.ACL().TokenUpdate(ctx, req, consul.WriteOptionsFromContext(ctx))
}

func (a *acl) DeleteToken(ctx context.Context, id string) (*consul.Response[bool], error) {
	return a.client.ACL().TokenDelete(ctx, id, consul.WriteOptionsFromContext(ctx))
}

func (a *acl) ListPolicies(ctx context.Context) (*consul.Response[[]*consul.ACLPolicy], error) {
	return a.client.ACL().PolicyList(ctx, consul.QueryOptionsFromContext(ctx))
}

func (a *acl) ReadPolicy(ctx context.Context, id string) (*consul.Response[*consul.ACLPolicy], error) {
	return a.client.ACL().PolicyRead(ctx, id, consul.QueryOptionsFromContext(ctx))
}

func (a *acl) ReadPolicyByName(ctx context.Context, name string) (*consul.Response[*consul.ACLPolicy], error) {
	return a.client.ACL().PolicyReadByName(ctx, name, consul.QueryOptionsFromContext(ctx))
}

func (a *acl) CreatePolicy(ctx context.Context, req *consul.ACLPolicy) (*consul.Response[*consul.ACLPolicy], error) {
	return a.client.ACL().PolicyCreate(ctx, req, consul.WriteOptionsFromContext(ctx))
}

func (a *acl) UpdatePolicy(ctx context.Context, req *consul.ACLPolicy) (*consul.Response[*consul.ACLPolicy], error) {
	return a.client.ACL().PolicyUpdate(ctx, req, consul.WriteOptionsFromContext(ctx))
}

func (a *acl) DeletePolicy(ctx context.Context, id string) (*consul.Response[bool], error) {
	return a.client.ACL().PolicyDelete(ctx, id, consul.WriteOptionsFromContext(ctx))
}

func (a *acl) ListRoles(ctx context.Context) (*consul.Response[[]*consul.ACLRole], error) {
	return a.client.ACL().RoleList(ctx, consul.QueryOptionsFromContext(ctx))
}

func (a *acl) ReadRole(ctx context.Context, id string) (*consul.Response[*consul.ACLRole], error) {
	return a.client.ACL().RoleRead(ctx, id, consul.QueryOptionsFromContext(ctx))
}

func (a *acl) ReadRoleByName(ctx context.Context, name string) (*consul.Response[*consul.ACLRole], error) {
	return a.client.ACL().RoleReadByName(ctx, name, consul.QueryOptionsFromContext(ctx))
}

func (a *acl) CreateRole(ctx context.Context, req *consul.ACLRole) (*consul.Response[*consul.ACLRole], error) {
	return a.client.ACL().RoleCreate(ctx, req, consul.WriteOptionsFromContext(ctx))
}

func (a *acl) UpdateRole(ctx context.Context, req *consul.ACLRole) (*consul.Response[*consul.ACLRole], error) {
	return a.client.ACL().RoleUpdate(ctx, req, consul.WriteOptionsFromContext(ctx))
}

func (a *acl) DeleteRole(ctx context.Context, id string) (*consul.Response[bool], error) {
	return a.client.ACL().RoleDelete(ctx, id, consul.WriteOptionsFromContext(ctx))
}
