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
	"github.com/go-chi/chi/v5"
)

func (a *HTTPAdapter) ListKeys(w http.ResponseWriter, r *http.Request) {
	utoken := r.Header.Get(ConseeTokenHeaderKey)
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	keys, err := a.kvService.ListKeys(ctx)
	if err != nil {
		errorResponse(w, err)
		return
	}
	response(w, keys)
}

func (a *HTTPAdapter) GetKV(w http.ResponseWriter, r *http.Request) {
	utoken := r.Header.Get(ConseeTokenHeaderKey)
	b64key := chi.URLParam(r, "b64key")

	k, err := base64.StdEncoding.DecodeString(b64key)
	if err != nil {
		errorResponse(w, &StatusError{Err: err, Process: "decoding b64key", Status: http.StatusBadRequest})
		return
	}
	// for history reading, we should check if the user has access to this kv
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	kv, err := a.kvService.Get(ctx, string(k))
	if err != nil {
		errorResponse(w, err)
		return
	}
	v := r.URL.Query().Get("v")
	if v != "" {
		hv, err := a.adminService.GetKVHistoryValue(ctx, b64key, v)
		if err != nil {
			errorResponse(w, err)
			return
		}
		kv.Value = hv
	}
	response(w, kv)
}

func (a *HTTPAdapter) GetValueType(w http.ResponseWriter, r *http.Request) {
	b64key := chi.URLParam(r, "b64key")
	vt, err := a.adminService.GetValueType(r.Context(), b64key)
	if err != nil {
		errorResponse(w, err)
		return
	}
	response(w, vt)
}

func (a *HTTPAdapter) UpdateValueType(w http.ResponseWriter, r *http.Request) {
	b64key := chi.URLParam(r, "b64key")
	b, _ := io.ReadAll(r.Body)
	err := a.kvService.UpdateType(r.Context(), b64key, string(b))
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *HTTPAdapter) BatchUpdateKV(w http.ResponseWriter, r *http.Request) {
	utoken := r.Header.Get(ConseeTokenHeaderKey)

	var req BatchUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorResponse(w, &StatusError{Err: err, Process: "decoding body", Status: http.StatusBadRequest})
		return
	}
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	ctx = consul.ContextWithWriteOptions(ctx, &consul.WriteOptions{Token: utoken})
	err = a.kvService.BatchUpdate(ctx, &req)
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *HTTPAdapter) CreateKV(w http.ResponseWriter, r *http.Request) {
	utoken := r.Header.Get(ConseeTokenHeaderKey)
	var req CreateKeyValueRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorResponse(w, &StatusError{Err: err, Process: "decoding body", Status: http.StatusBadRequest})
		return
	}
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	ctx = consul.ContextWithWriteOptions(ctx, &consul.WriteOptions{Token: utoken})
	err = a.kvService.Create(ctx, &req)
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *HTTPAdapter) UpdateKV(w http.ResponseWriter, r *http.Request) {
	utoken := r.Header.Get(ConseeTokenHeaderKey)
	b64key := chi.URLParam(r, "b64key")

	k, err := base64.StdEncoding.DecodeString(b64key)
	if err != nil {
		errorResponse(w, &StatusError{Err: err, Process: "decoding b64key", Status: http.StatusBadRequest})
		return
	}

	var req UpdateValueRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorResponse(w, &StatusError{Err: err, Process: "decoding body", Status: http.StatusBadRequest})
		return
	}
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	ctx = consul.ContextWithWriteOptions(ctx, &consul.WriteOptions{Token: utoken})
	err = a.kvService.Update(ctx, string(k), req.Value)
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (a *HTTPAdapter) UpdateKVValueType(w http.ResponseWriter, r *http.Request) {
	// utoken := r.Header.Get(ConseeTokenHeaderKey)
}

func (a *HTTPAdapter) DeleteKV(w http.ResponseWriter, r *http.Request) {
	utoken := r.Header.Get(ConseeTokenHeaderKey)
	b64key := chi.URLParam(r, "b64key")

	k, err := base64.StdEncoding.DecodeString(b64key)
	if err != nil {
		errorResponse(w, &StatusError{Err: err, Process: "decoding b64key", Status: http.StatusBadRequest})
		return
	}

	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	ctx = consul.ContextWithWriteOptions(ctx, &consul.WriteOptions{Token: utoken})
	err = a.kvService.Delete(ctx, string(k))
	if err != nil {
		errorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
