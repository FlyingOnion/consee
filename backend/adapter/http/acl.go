// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package httpadapter

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"

	. "github.com/FlyingOnion/consee/backend/common"
	"github.com/FlyingOnion/consee/backend/consul"
	"github.com/FlyingOnion/consee/backend/service"
	"github.com/go-chi/chi/v5"
)

func (a *HTTPAdapter) ApplyToken(w http.ResponseWriter, r *http.Request) {
	var req TokenApplicationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorResponse(w, &StatusError{Err: err, Process: "decoding body", Status: http.StatusBadRequest})
		return
	}
	utoken := r.Header.Get(ConseeTokenHeaderKey)
	if utoken != req.SecretID {
		errorResponse(w, &StatusError{Err: errInvalidToken, Status: http.StatusForbidden})
		return
	}
	resp, err := a.aclService.CreateTokenApplicationRequest(r.Context(), &req)
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	response(w, resp)
}

func (a *HTTPAdapter) ParseRule(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		errorResponse(w, &StatusError{Err: err, Process: "decoding body", Status: http.StatusBadRequest})
		return
	}
	var ruleList service.HCLRuleList
	// TODO: this call seems weird, but it's ok for now.
	err = service.ParseHCLRules(string(b), &ruleList)
	if err != nil {
		errorResponse(w, &StatusError{Err: err, Process: "parsing rule", Status: http.StatusBadRequest})
		return
	}
	response(w, ruleList.ToParsedRuleList())
}

func (a *HTTPAdapter) HandleTokenApplication(w http.ResponseWriter, r *http.Request) {
	var req HandleTokenApplicationRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorResponse(w, &StatusError{Err: err, Process: "decoding body", Status: http.StatusBadRequest})
		return
	}
	accessorId := chi.URLParam(r, "id")
	utoken := r.Header.Get(ConseeTokenHeaderKey)
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	ctx = consul.ContextWithWriteOptions(ctx, &consul.WriteOptions{Token: utoken})
	err = a.aclService.ReviewTokenApplicationRequest(ctx, accessorId, &req)
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *HTTPAdapter) ListACLTokens(w http.ResponseWriter, r *http.Request) {
	utoken := r.Header.Get(ConseeTokenHeaderKey)
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})

	tokens, err := a.aclService.ListTokens(ctx)
	if err != nil {
		errorResponse(w, err)
		return
	}
	response(w, tokens)
}

func (a *HTTPAdapter) ReadACLToken(w http.ResponseWriter, r *http.Request) {
	utoken := r.Header.Get(ConseeTokenHeaderKey)
	accessorId := chi.URLParam(r, "id")
	if accessorId == "" {
		errorResponse(w, &StatusError{
			Process: "reading token accessor id",
			Status:  http.StatusBadRequest,
			Err:     errIdEmpty,
		})
		return
	}
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	t, err := a.aclService.ReadToken(ctx, accessorId)
	if err != nil {
		errorResponse(w, err)
		return
	}
	response(w, t)
}

func (a *HTTPAdapter) CreateACLToken(w http.ResponseWriter, r *http.Request) {
	var req CreateTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorResponse(w, &StatusError{Err: err, Process: "decoding body", Status: http.StatusBadRequest})
		return
	}

	utoken := r.Header.Get(ConseeTokenHeaderKey)
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	ctx = consul.ContextWithWriteOptions(ctx, &consul.WriteOptions{Token: utoken})
	err = a.aclService.CreateToken(ctx, &req)
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *HTTPAdapter) UpdateACLToken(w http.ResponseWriter, r *http.Request) {
	var req UpdateTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorResponse(w, &StatusError{Err: err, Process: "decoding body", Status: http.StatusBadRequest})
		return
	}

	accessorId := chi.URLParam(r, "id")
	utoken := r.Header.Get(ConseeTokenHeaderKey)
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	ctx = consul.ContextWithWriteOptions(ctx, &consul.WriteOptions{Token: utoken})
	err = a.aclService.UpdateToken(ctx, accessorId, &req)
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *HTTPAdapter) DeleteACLToken(w http.ResponseWriter, r *http.Request) {
	accessorId := chi.URLParam(r, "id")
	if accessorId == "" {
		errorResponse(w, &StatusError{
			Process: "reading token accessor id",
			Status:  http.StatusBadRequest,
			Err:     errIdEmpty,
		})
		return
	}
	utoken := r.Header.Get(ConseeTokenHeaderKey)
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	ctx = consul.ContextWithWriteOptions(ctx, &consul.WriteOptions{Token: utoken})
	err := a.aclService.DeleteToken(ctx, accessorId)
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *HTTPAdapter) ListACLPolicies(w http.ResponseWriter, r *http.Request) {
	utoken := r.Header.Get(ConseeTokenHeaderKey)
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})

	v, ok := r.URL.Query()["exclusive"]
	exclusive := ""
	if ok {
		if len(v) == 0 || v[0] == "" || v[0] == "1" {
			exclusive = "1"
		} else {
			exclusive = "0"
		}
	}

	options := ListPoliciesOptions{Exclusive: exclusive}
	policies, err := a.aclService.ListPolicies(ctx, options)
	if err != nil {
		errorResponse(w, err)
		return
	}
	response(w, policies)
}

func (a *HTTPAdapter) CreateACLPolicy(w http.ResponseWriter, r *http.Request) {
	var req CreatePolicyRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorResponse(w, &StatusError{Err: err, Process: "decoding body", Status: http.StatusBadRequest})
		return
	}

	utoken := r.Header.Get(ConseeTokenHeaderKey)
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	ctx = consul.ContextWithWriteOptions(ctx, &consul.WriteOptions{Token: utoken})
	err = a.aclService.CreatePolicy(ctx, &req)
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *HTTPAdapter) ReadACLPolicy(w http.ResponseWriter, r *http.Request) {
	b64name := chi.URLParam(r, "b64name")
	name, err := base64.StdEncoding.DecodeString(b64name)
	if err != nil {
		errorResponse(w, &StatusError{
			Process: "decoding policy name",
			Status:  http.StatusBadRequest,
			Err:     errIdEmpty,
		})
		return
	}
	utoken := r.Header.Get(ConseeTokenHeaderKey)
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	resp, err := a.aclService.ReadPolicy(ctx, string(name))
	if err != nil {
		errorResponse(w, err)
		return
	}
	response(w, resp)
}

func (a *HTTPAdapter) UpdatePolicyRule(w http.ResponseWriter, r *http.Request) {
	b64name := chi.URLParam(r, "b64name")
	name, err := base64.StdEncoding.DecodeString(b64name)
	if err != nil {
		errorResponse(w, &StatusError{
			Process: "decoding policy name",
			Status:  http.StatusBadRequest,
			Err:     errIdEmpty,
		})
		return
	}
	newRule, _ := io.ReadAll(r.Body)
	utoken := r.Header.Get(ConseeTokenHeaderKey)
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	ctx = consul.ContextWithWriteOptions(ctx, &consul.WriteOptions{Token: utoken})
	err = a.aclService.UpdatePolicyRule(ctx, string(name), string(newRule))
	if err != nil {
		errorResponse(w, NewStatusError(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *HTTPAdapter) DeleteACLPolicy(w http.ResponseWriter, r *http.Request) {
	b64name := chi.URLParam(r, "b64name")
	name, err := base64.StdEncoding.DecodeString(b64name)
	if err != nil {
		errorResponse(w, &StatusError{
			Process: "decoding policy name",
			Status:  http.StatusBadRequest,
			Err:     errIdEmpty,
		})
		return
	}
	utoken := r.Header.Get(ConseeTokenHeaderKey)
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	ctx = consul.ContextWithWriteOptions(ctx, &consul.WriteOptions{Token: utoken})
	err = a.aclService.DeletePolicy(ctx, string(name))
	if err != nil {
		errorResponse(w, NewStatusError(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *HTTPAdapter) ListACLRoles(w http.ResponseWriter, r *http.Request) {
	utoken := r.Header.Get(ConseeTokenHeaderKey)
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	
	roles, err := a.aclService.ListRoles(ctx)
	if err != nil {
		errorResponse(w, err)
		return
	}
	response(w, roles)
}

func (a *HTTPAdapter) CreateACLRole(w http.ResponseWriter, r *http.Request) {
	var req CreateRoleRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorResponse(w, &StatusError{Err: err, Process: "decoding body", Status: http.StatusBadRequest})
		return
	}

	utoken := r.Header.Get(ConseeTokenHeaderKey)
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	ctx = consul.ContextWithWriteOptions(ctx, &consul.WriteOptions{Token: utoken})
	err = a.aclService.CreateRole(ctx, &req)
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *HTTPAdapter) ReadACLRole(w http.ResponseWriter, r *http.Request) {
	b64name := chi.URLParam(r, "b64name")
	name, err := base64.StdEncoding.DecodeString(b64name)
	if err != nil {
		errorResponse(w, &StatusError{
			Process: "decoding role name",
			Status:  http.StatusBadRequest,
			Err:     err,
		})
		return
	}
	
	utoken := r.Header.Get(ConseeTokenHeaderKey)
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	role, err := a.aclService.ReadRole(ctx, string(name))
	if err != nil {
		errorResponse(w, err)
		return
	}
	response(w, role)
}

func (a *HTTPAdapter) UpdateACLRole(w http.ResponseWriter, r *http.Request) {
	b64name := chi.URLParam(r, "b64name")
	name, err := base64.StdEncoding.DecodeString(b64name)
	if err != nil {
		errorResponse(w, &StatusError{
			Process: "decoding role name",
			Status:  http.StatusBadRequest,
			Err:     err,
		})
		return
	}
	
	var req UpdateRoleRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorResponse(w, &StatusError{Err: err, Process: "decoding body", Status: http.StatusBadRequest})
		return
	}
	
	utoken := r.Header.Get(ConseeTokenHeaderKey)
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	ctx = consul.ContextWithWriteOptions(ctx, &consul.WriteOptions{Token: utoken})
	err = a.aclService.UpdateRole(ctx, string(name), &req)
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *HTTPAdapter) DeleteACLRole(w http.ResponseWriter, r *http.Request) {
	b64name := chi.URLParam(r, "b64name")
	name, err := base64.StdEncoding.DecodeString(b64name)
	if err != nil {
		errorResponse(w, &StatusError{
			Process: "decoding role name",
			Status:  http.StatusBadRequest,
			Err:     err,
		})
		return
	}
	
	utoken := r.Header.Get(ConseeTokenHeaderKey)
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	ctx = consul.ContextWithWriteOptions(ctx, &consul.WriteOptions{Token: utoken})
	err = a.aclService.DeleteRole(ctx, string(name))
	if err != nil {
		errorResponse(w, NewStatusError(err))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
