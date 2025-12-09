// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package httpadapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	. "github.com/FlyingOnion/consee/backend/common"
	"github.com/FlyingOnion/consee/backend/consul"
	"github.com/FlyingOnion/consee/backend/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

const ConseeTokenHeaderKey = "G-Consee-Token"
const ConseeErrorHeaderKey = "G-Consee-Error"

func fuzz(s string) string {
	if len(s) <= 8 {
		return s
	}
	b := bytes.Repeat([]byte("*"), 8)
	b[0], b[1], b[len(b)-2], b[len(b)-1] = s[0], s[1], s[len(s)-2], s[len(s)-1]
	return string(b)
}

// errorResponse returns an error errorResponse to the client;
//
// err should be non-nil;
//
// response body should be empty;
//
// error messages are all in header.
func errorResponse(rw http.ResponseWriter, err error) {
	if statusErr, ok := err.(*StatusError); ok {
		rw.Header().Set(ConseeErrorHeaderKey, statusErr.Error())
		rw.WriteHeader(statusErr.Status)
		return
	}
	if statusErr, ok := err.(StatusError); ok {
		rw.Header().Set(ConseeErrorHeaderKey, statusErr.Error())
		rw.WriteHeader(statusErr.Status)
		return
	}
	if domainErr, ok := err.(*service.DomainError); ok {
		statusErr := NewStatusError(domainErr)
		rw.Header().Set(ConseeErrorHeaderKey, statusErr.Error())
		rw.WriteHeader(statusErr.Status)
		return
	}
	if domainErr, ok := err.(service.DomainError); ok {
		statusErr := NewStatusError(domainErr)
		rw.Header().Set(ConseeErrorHeaderKey, statusErr.Error())
		rw.WriteHeader(statusErr.Status)
		return
	}
	statusErr := StatusError{Err: err, Status: http.StatusInternalServerError}
	rw.Header().Set(ConseeErrorHeaderKey, statusErr.Error())
	rw.WriteHeader(http.StatusInternalServerError)
}

// response returns 200 and data.
func response(rw http.ResponseWriter, data any) {
	if data == nil {
		return
	}
	if b, ok := data.([]byte); ok {
		rw.Write(b)
		return
	}
	if s, ok := data.(string); ok {
		rw.Write([]byte(s))
		return
	}
	err := json.NewEncoder(rw).Encode(data)
	if err != nil {
		errorResponse(rw, StatusError{Err: err, Process: "writing response", Status: http.StatusInternalServerError})
	}
}

type HTTPAdapter struct {
	a2           service.All
	kvService    service.KVService
	aclService   service.ACLService
	adminService service.AdminService
}

func NewAdapter(a2 service.All, kvService service.KVService, aclService service.ACLService, adminService service.AdminService) *HTTPAdapter {
	return &HTTPAdapter{
		a2:           a2,
		kvService:    kvService,
		aclService:   aclService,
		adminService: adminService,
	}
}

// func (a *HTTPAdapter) KVHandler() http.Handler {
// 	kv := chi.NewRouter()
// 	kv.Use(a.CheckUserToken)
// 	kv.Get("/keys", a.ListKeys)
// 	kv.Get("/value/{b64key}", a.GetKV)
// 	kv.Get("/history/{b64key}", a.GetKVHistory)
// 	kv.Get("/valuetype/{b64key}", a.GetValueType)
// 	kv.Put("/valuetype/{b64key}", a.UpdateValueType)
// 	kv.Post("/value", a.CreateKV)
// 	kv.Put("/value/{b64key}", a.UpdateKV)
// 	kv.Put("/value-type/{b64key}", a.UpdateKVValueType)
// 	kv.Delete("/value/{b64key}", a.DeleteKV)
// 	kv.Put("/batch", a.checkAdminToken(http.HandlerFunc(a.BatchUpdateKV)))
// 	return kv
// }

// func (a *HTTPAdapter) ACLHandler() http.Handler {
// 	r := chi.NewRouter()
// 	r.Post("/acl/token-apply", a.ApplyToken)
// 	r.Post("/acl/hcl-rule", a.ParseRule)
// 	r.Group(func(r1 chi.Router) {
// 		r1.Use(a.CheckUserToken)
// 		r1.Put("/token-apply/{id}", a.checkAdminToken(http.HandlerFunc(a.HandleTokenApplication)))

// 		r1.Get("/tokens", a.ListACLTokens)
// 		r1.Get("/token/{id}", a.ReadACLToken)
// 		r1.Post("/token", a.CreateACLToken)
// 		r1.Put("/token/{id}", a.UpdateACLToken)
// 		r1.Delete("/token/{id}", a.DeleteACLToken)

// 		r1.Get("/policies", a.ListACLPolicies)
// 		r1.Post("/policy", a.CreateACLPolicy)
// 		r1.Get("/policy/{b64name}", a.ReadACLPolicy)
// 		r1.Delete("/policy/{b64name}", a.DeleteACLPolicy)
// 	})

// 	return r
// }

func (a *HTTPAdapter) HttpHandler() http.Handler {
	return a.chiHandler()
}

func (a *HTTPAdapter) chiHandler() *chi.Mux {
	r := chi.NewRouter()
	r.Route("/api", func(rApi chi.Router) {
		rApi.Route("/v0", func(rApiV0 chi.Router) {
			rApiV0.Post("/authenticate", a.Authenticate)
			rApiV0.Group(func(sub chi.Router) {
				sub.Use(a.CheckUserToken, a.CheckAdminToken)
				sub.Post("/export", a.Export)
				sub.Post("/import", a.Import)
			})
			rApiV0.Route("/kv", func(kv chi.Router) {
				kv.Use(a.CheckUserToken)
				kv.Get("/keys", a.ListKeys)
				kv.Get("/value/{b64key}", a.GetKV)
				kv.Get("/valuetype/{b64key}", a.GetValueType)
				kv.Put("/valuetype/{b64key}", a.UpdateValueType)
				kv.Post("/value", a.CreateKV)
				kv.Put("/value/{b64key}", a.UpdateKV)
				kv.Put("/value-type/{b64key}", a.UpdateKVValueType)
				kv.Delete("/value/{b64key}", a.DeleteKV)
				kv.Put("/batch", a.checkAdminToken(http.HandlerFunc(a.BatchUpdateKV)))
			})
			rApiV0.Route("/acl", func(acl chi.Router) {
				acl.Post("/token-request", a.ApplyToken)
				acl.Post("/hcl-rule", a.ParseRule)
				acl.Group(func(sub chi.Router) {
					sub.Use(a.CheckUserToken)

					sub.Get("/tokens", a.ListACLTokens)
					sub.Get("/token/{id}", a.ReadACLToken)
					sub.Post("/token", a.CreateACLToken)
					sub.Put("/token/{id}", a.UpdateACLToken)
					sub.Delete("/token/{id}", a.DeleteACLToken)

					sub.Get("/policies", a.ListACLPolicies)
					sub.Post("/policy", a.CreateACLPolicy)
					sub.Get("/policy/{b64name}", a.ReadACLPolicy)
					sub.Delete("/policy/{b64name}", a.DeleteACLPolicy)
				})
			})
		})
	})
	r.HandleFunc("/ui", a.ServeUI)
	r.HandleFunc("/ui/*", a.ServeUI)
	return r
}

func (a *HTTPAdapter) ServeUI(w http.ResponseWriter, r *http.Request) {
	srcPath := filepath.Join("public", strings.TrimPrefix(r.URL.Path, "/ui"))
	// SPA 回退：如果文件不存在则返回 index.html
	if _, err := os.Stat(srcPath); err == nil {
		http.ServeFile(w, r, srcPath)
	} else {
		http.ServeFile(w, r, "public/index.html")
	}
}

func (a *HTTPAdapter) Authenticate(w http.ResponseWriter, r *http.Request) {
	utoken := r.Header.Get(ConseeTokenHeaderKey)
	if utoken == "" {
		errorResponse(w, &StatusError{
			Err:    errTokenEmpty,
			Status: http.StatusBadRequest,
		})
		return
	}
	_, err := uuid.Parse(utoken)
	if err != nil {
		errorResponse(w, StatusError{
			Err:    fmt.Errorf("invalid token %s (should be a valid uuid)", fuzz(utoken)),
			Status: http.StatusBadRequest,
		})
		return
	}
	ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
	err = a.aclService.ValidateToken(ctx)
	if err != nil {
		errorResponse(w, StatusError{
			Process: "authentication",
			Status:  http.StatusUnauthorized,
			Err:     errInvalidToken,
		})
		return
	}
	err = a.aclService.CheckAdmin(ctx)
	if err != nil {
		response(w, AuthenticateResult{1, 0, 0})
		return
	}
	nOpenNotifications, err := a.adminService.GetOpenNotificationsCount(ctx)
	if err != nil {
		response(w, AuthenticateResult{1, 1, 0})
		return
	}
	response(w, AuthenticateResult{1, 1, nOpenNotifications})
}

func (a *HTTPAdapter) CheckUserToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		utoken := r.Header.Get(ConseeTokenHeaderKey)
		if utoken == "" {
			errorResponse(w, &StatusError{
				Err:    errTokenEmpty,
				Status: http.StatusBadRequest,
			})
			return
		}
		_, err := uuid.Parse(utoken)
		if err != nil {
			errorResponse(w, StatusError{
				Err:    fmt.Errorf("invalid token %s (should be a valid uuid)", fuzz(utoken)),
				Status: http.StatusBadRequest,
			})
			return
		}
		ctx := consul.ContextWithQueryOptions(r.Context(), &consul.QueryOptions{Token: utoken})
		err = a.aclService.ValidateToken(ctx)
		if err != nil {
			errorResponse(w, err)
			return
		}
		if next != nil {
			next.ServeHTTP(w, r)
		}
	})
}
