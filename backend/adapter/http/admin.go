// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package httpadapter

import (
	"encoding/json"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	. "github.com/FlyingOnion/consee/backend/common"
	"github.com/FlyingOnion/consee/backend/consul"
)

// checkAdminToken checks if the token provided has admin permission.
// validation of uuid is skipped, so it should be placed after checkUserToken
func (a *HTTPAdapter) CheckAdminToken(next http.Handler) http.Handler {
	return a.checkAdminToken(next)
}

func (a *HTTPAdapter) checkAdminToken(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utoken := r.Header.Get(ConseeTokenHeaderKey)
		ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
		err := a.aclService.CheckAdmin(ctx)
		if err != nil {
			errorResponse(w, err)
			return
		}
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}

func (a *HTTPAdapter) Import(w http.ResponseWriter, r *http.Request) {
	v := r.FormValue("dryrun")
	dryrun := v == "1"
	file, header, err := r.FormFile("file")
	if err != nil {
		errorResponse(w, StatusError{
			Status:  http.StatusBadRequest,
			Process: "parsing file",
			Err:     errInvalidFile,
		})
		return
	}
	fileExt := filepath.Ext(header.Filename)
	switch fileExt {
	case ".zip":
	case ".json":
	default:
		errorResponse(w, StatusError{
			Status:  http.StatusBadRequest,
			Process: "parsing file",
			Err:     errInvalidFileFormat,
		})
		return
	}
	b := make([]byte, header.Size)
	file.Read(b)
	file.Close()

	utoken := r.Header.Get(ConseeTokenHeaderKey)
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	ctx = consul.ContextWithWriteOptions(ctx, &consul.WriteOptions{Token: utoken})

	resp, err := a.a2.Import(ctx, &ImportRequest{
		Format:      fileExt[1:],
		Dryrun:      dryrun,
		FileContent: b,
	})
	if err != nil {
		errorResponse(w, err)
		return
	}
	response(w, resp)
}

func (a *HTTPAdapter) Export(w http.ResponseWriter, r *http.Request) {
	var req ExportRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		errorResponse(w, &StatusError{Err: err, Process: "decoding body", Status: http.StatusBadRequest})
		return
	}

	utoken := r.Header.Get(ConseeTokenHeaderKey)
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})

	b, err := a.a2.Export(ctx, &req)
	if err != nil {
		errorResponse(w, err)
		return
	}
	now := time.Now().Format("20060102-150405")
	// 设置响应头
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=consee-export-"+now+".zip")
	w.Header().Set("Content-Length", strconv.Itoa(len(b)))
	w.Write(b)
}
