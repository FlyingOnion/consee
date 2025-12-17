// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package consul

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

type ACLLink struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ACLToken struct {
	CreateIndex uint64
	ModifyIndex uint64
	AccessorID  string
	SecretID    string
	Description string
	Policies    []*ACLLink `json:",omitempty"`
	Roles       []*ACLLink `json:",omitempty"`
	// ServiceIdentities []*ACLServiceIdentity `json:",omitempty"`
	// NodeIdentities    []*ACLNodeIdentity    `json:",omitempty"`
	// TemplatedPolicies []*ACLTemplatedPolicy `json:",omitempty"`
	Local          bool
	AuthMethod     string        `json:",omitempty"`
	ExpirationTTL  time.Duration `json:",omitempty"`
	ExpirationTime *time.Time    `json:",omitempty"`
	CreateTime     time.Time     `json:",omitempty"`
	Hash           []byte        `json:",omitempty"`

	// DEPRECATED (ACL-Legacy-Compat)
	// Rules are an artifact of legacy tokens deprecated in Consul 1.4
	Rules string `json:"-"`

	// Namespace is the namespace the ACLToken is associated with.
	// Namespaces are a Consul Enterprise feature.
	Namespace string `json:",omitempty"`

	// Partition is the partition the ACLToken is associated with.
	// Partitions are a Consul Enterprise feature.
	Partition string `json:",omitempty"`

	// AuthMethodNamespace is the namespace the token's AuthMethod is associated with.
	// Namespacing is a Consul Enterprise feature.
	AuthMethodNamespace string `json:",omitempty"`
}

// ACLPolicy represents an ACL Policy.
type ACLPolicy struct {
	ID          string
	Name        string
	Description string
	Rules       string
	Datacenters []string
	Hash        []byte
	CreateIndex uint64
	ModifyIndex uint64

	// Namespace is the namespace the ACLPolicy is associated with.
	// Namespacing is a Consul Enterprise feature.
	Namespace string `json:",omitempty"`

	// Partition is the partition the ACLPolicy is associated with.
	// Partitions are a Consul Enterprise feature.
	Partition string `json:",omitempty"`
}

// ACLRole represents an ACL Role.
type ACLRole struct {
	ID          string
	Name        string
	Description string
	Policies    []*ACLLink `json:",omitempty"`
	// ServiceIdentities []*ACLServiceIdentity `json:",omitempty"`
	// NodeIdentities    []*ACLNodeIdentity    `json:",omitempty"`
	// TemplatedPolicies []*ACLTemplatedPolicy `json:",omitempty"`
	Hash        []byte
	CreateIndex uint64
	ModifyIndex uint64

	// Namespace is the namespace the ACLRole is associated with.
	// Namespacing is a Consul Enterprise feature.
	Namespace string `json:",omitempty"`

	// Partition is the partition the ACLRole is associated with.
	// Partitions are a Consul Enterprise feature.
	Partition string `json:",omitempty"`
}

type ACL struct {
	c *Client
}

func (c *Client) ACL() *ACL {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.acl == nil {
		c.acl = &ACL{c}
	}
	return c.acl
}

func (a *ACL) TokenReadSelf(ctx context.Context, q *QueryOptions) (*Response[*ACLToken], error) {
	httpReq := a.c.newRequest(ctx, http.MethodGet, "/v1/acl/token/self", q.toRequestOptions()...)
	return responseDirectly(a.c.httpClient, httpReq, decodeACLToken)
}

func (a *ACL) TokenList(ctx context.Context, q *QueryOptions) (*Response[[]*ACLToken], error) {
	httpReq := a.c.newRequest(ctx, http.MethodGet, "/v1/acl/tokens", q.toRequestOptions()...)
	return responseDirectly(a.c.httpClient, httpReq, decodeACLTokenList)
}

func (a *ACL) TokenListFiltered(ctx context.Context, q *QueryOptions, t ACLTokenFilterOptions) (*Response[[]*ACLToken], error) {
	options := append(q.toRequestOptions(), t.toRequestOptions()...)
	httpReq := a.c.newRequest(ctx, http.MethodGet, "/v1/acl/tokens", options...)
	return responseDirectly(a.c.httpClient, httpReq, decodeACLTokenList)
}

func (a *ACL) TokenRead(ctx context.Context, id string, q *QueryOptions) (*Response[*ACLToken], error) {
	httpReq := a.c.newRequest(ctx, http.MethodGet, "/v1/acl/token/"+id, q.toRequestOptions()...)
	return responseDirectly(a.c.httpClient, httpReq, decodeACLToken)
}

func (a *ACL) TokenCreate(ctx context.Context, req *ACLToken, w *WriteOptions) (*Response[*ACLToken], error) {
	b, _ := json.Marshal(req)
	options := append(w.toRequestOptions(),
		reqWithContentType("application/json"),
		reqWithBody(b),
	)
	httpReq := a.c.newRequest(ctx, http.MethodPut, "/v1/acl/token", options...)
	return responseDirectly(a.c.httpClient, httpReq, decodeACLToken)
}

func (a *ACL) TokenUpdate(ctx context.Context, req *ACLToken, w *WriteOptions) (*Response[*ACLToken], error) {
	b, _ := json.Marshal(req)
	options := append(w.toRequestOptions(),
		reqWithContentType("application/json"),
		reqWithBody(b),
	)
	httpReq := a.c.newRequest(ctx, http.MethodPut, "/v1/acl/token/"+req.AccessorID, options...)
	return responseDirectly(a.c.httpClient, httpReq, decodeACLToken)
}

func (a *ACL) TokenDelete(ctx context.Context, id string, w *WriteOptions) (*Response[bool], error) {
	httpReq := a.c.newRequest(ctx, http.MethodDelete, "/v1/acl/token/"+id, w.toRequestOptions()...)
	return responseDirectly(a.c.httpClient, httpReq, decodeTrue)
}

func (a *ACL) PolicyList(ctx context.Context, q *QueryOptions) (*Response[[]*ACLPolicy], error) {
	httpReq := a.c.newRequest(ctx, http.MethodGet, "/v1/acl/policies", q.toRequestOptions()...)
	return responseDirectly(a.c.httpClient, httpReq, decodeACLPolicyList)
}

func (a *ACL) PolicyRead(ctx context.Context, id string, q *QueryOptions) (*Response[*ACLPolicy], error) {
	httpReq := a.c.newRequest(ctx, http.MethodGet, "/v1/acl/policy/"+id, q.toRequestOptions()...)
	return responseDirectly(a.c.httpClient, httpReq, decodeACLPolicy)
}

func (a *ACL) PolicyReadByName(ctx context.Context, name string, q *QueryOptions) (*Response[*ACLPolicy], error) {
	httpReq := a.c.newRequest(ctx, http.MethodGet, "/v1/acl/policy/name/"+url.QueryEscape(name), q.toRequestOptions()...)
	return responseDirectly(a.c.httpClient, httpReq, decodeACLPolicy)
}

func (a *ACL) PolicyCreate(ctx context.Context, req *ACLPolicy, w *WriteOptions) (*Response[*ACLPolicy], error) {
	b, _ := json.Marshal(req)
	options := append(w.toRequestOptions(),
		reqWithBody(b),
	)
	httpReq := a.c.newRequest(ctx, http.MethodPut, "/v1/acl/policy", options...)
	return responseDirectly(a.c.httpClient, httpReq, decodeACLPolicy)
}

func (a *ACL) PolicyUpdate(ctx context.Context, req *ACLPolicy, w *WriteOptions) (*Response[*ACLPolicy], error) {
	b, _ := json.Marshal(req)
	options := append(w.toRequestOptions(),
		reqWithBody(b),
	)
	httpReq := a.c.newRequest(ctx, http.MethodPut, "/v1/acl/policy/"+req.ID, options...)
	return responseDirectly(a.c.httpClient, httpReq, decodeACLPolicy)
}

func (a *ACL) PolicyDelete(ctx context.Context, id string, w *WriteOptions) (*Response[bool], error) {
	httpReq := a.c.newRequest(ctx, http.MethodDelete, "/v1/acl/policy/"+id, w.toRequestOptions()...)
	return responseDirectly(a.c.httpClient, httpReq, decodeTrue)
}

func (a *ACL) RoleList(ctx context.Context, q *QueryOptions) (*Response[[]*ACLRole], error) {
	httpReq := a.c.newRequest(ctx, http.MethodGet, "/v1/acl/roles", q.toRequestOptions()...)
	return responseDirectly(a.c.httpClient, httpReq, decodeACLRoleList)
}

func (a *ACL) RoleRead(ctx context.Context, id string, q *QueryOptions) (*Response[*ACLRole], error) {
	httpReq := a.c.newRequest(ctx, http.MethodGet, "/v1/acl/role/"+id, q.toRequestOptions()...)
	return responseDirectly(a.c.httpClient, httpReq, decodeACLRole)
}

func (a *ACL) RoleReadByName(ctx context.Context, name string, q *QueryOptions) (*Response[*ACLRole], error) {
	httpReq := a.c.newRequest(ctx, http.MethodGet, "/v1/acl/role/name/"+url.QueryEscape(name), q.toRequestOptions()...)
	return responseDirectly(a.c.httpClient, httpReq, decodeACLRole)
}

func (a *ACL) RoleCreate(ctx context.Context, req *ACLRole, w *WriteOptions) (*Response[*ACLRole], error) {
	b, _ := json.Marshal(req)
	options := append(w.toRequestOptions(),
		reqWithBody(b),
	)
	httpReq := a.c.newRequest(ctx, http.MethodPut, "/v1/acl/role", options...)
	return responseDirectly(a.c.httpClient, httpReq, decodeACLRole)
}

func (a *ACL) RoleUpdate(ctx context.Context, req *ACLRole, w *WriteOptions) (*Response[*ACLRole], error) {
	b, _ := json.Marshal(req)
	options := append(w.toRequestOptions(),
		reqWithBody(b),
	)
	httpReq := a.c.newRequest(ctx, http.MethodPut, "/v1/acl/role/"+req.ID, options...)
	return responseDirectly(a.c.httpClient, httpReq, decodeACLRole)
}

func (a *ACL) RoleDelete(ctx context.Context, id string, w *WriteOptions) (*Response[bool], error) {
	httpReq := a.c.newRequest(ctx, http.MethodDelete, "/v1/acl/role/"+id, w.toRequestOptions()...)
	return responseDirectly(a.c.httpClient, httpReq, decodeTrue)
}
