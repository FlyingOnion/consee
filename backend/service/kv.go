// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package service

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/FlyingOnion/consee/backend/buffer"
	. "github.com/FlyingOnion/consee/backend/common"
	"github.com/FlyingOnion/consee/backend/consul"
	"github.com/FlyingOnion/consee/backend/repo"
)

// func ListKeys() (keys []string, err error)
// func Get(key string) (*GetValueResponse, error)
// func Create(req *CreateKeyValueRequest) error
// func Update(key, value string) error
// func UpdateType(key, valueType string) error
// func Delete(key string) error

type KVService interface {
	ListKeys(ctx context.Context) (keys []string, err error)
	Get(ctx context.Context, key string) (*GetValueResponse, error)
	Create(ctx context.Context, req *CreateKeyValueRequest) error
	Update(ctx context.Context, key, value string) error
	UpdateType(ctx context.Context, key, valueType string) error
	BatchUpdate(ctx context.Context, req *BatchUpdateRequest) error
	Delete(ctx context.Context, key string) error

	WatchOpenNotificationsCount(ctx context.Context, cb func(n int))
}

type kvService struct {
	kv    repo.KVRepo
	admin AdminService
}

func NewKVService(kv repo.KVRepo, admin AdminService) KVService {
	return &kvService{
		kv:    kv,
		admin: admin,
	}
}

func (s *kvService) ListKeys(ctx context.Context) (keys []string, err error) {
	resp, err := s.kv.ListKeys(ctx, "", "")
	// resp, err := s.client.KV().Keys(ctx, "", "", consul.QueryOptionsFromContext(ctx))
	if err != nil {
		slog.Error("failed to list keys", "error", err)
		return nil, err
	}
	if len(resp.Body) == 0 {
		return resp.Body, nil
	}
	keysExcludingInternal := make([]string, 0, len(resp.Body))
	for _, k := range resp.Body {
		if strings.HasPrefix(k, ConseeInternalKeyPrefix) {
			continue
		}
		keysExcludingInternal = append(keysExcludingInternal, k)
	}
	return keysExcludingInternal, nil
}

func (s *kvService) Get(ctx context.Context, key string) (*GetValueResponse, error) {
	resp, err := s.kv.Read(ctx, key)
	// resp, err := s.client.KV().Get(ctx, key, consul.QueryOptionsFromContext(ctx))
	if err != nil {
		slog.Error("kvGet: failed to get key", "key", key, "error", err)
		return nil, errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		slog.Error("kvGet: permission denied", "key", key, "status", resp.Status)
		return nil, errPermissionDenied
	}
	if resp.Status == http.StatusNotFound {
		slog.Error("kvGet: key not found", "key", key)
		return nil, &DomainError{Code: DomainErrorCodeNotFound, Message: "key not found"}
	}

	return &GetValueResponse{
		Key:   resp.Body.Key,
		Value: string(resp.Body.Value),
	}, resp.Err
}

func (s *kvService) Create(ctx context.Context, req *CreateKeyValueRequest) error {
	resp1, err := s.kv.Read(ctx, req.Key)
	if err != nil {
		slog.Error("kvCreate: failed to read key", "key", req.Key, "error", err)
		return errFailedToConnectConsul
	}
	if resp1.Body != nil {
		slog.Error("kvCreate: key already exists", "key", req.Key)
		return &DomainError{Code: DomainErrorCodeAlreadyExists, Message: "key already exists"}
	}

	resp, err := s.kv.Write(ctx, req.Key, req.Value)
	// resp, err := s.client.KV().Put(ctx, &consul.KVPair{Key: req.Key, Value: []byte(req.Value)}, consul.WriteOptionsFromContext(ctx))
	if err != nil {
		slog.Error("kvCreate: failed to create key", "key", req.Key, "error", err)
		return errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		slog.Error("kvCreate: permission denied", "key", req.Key, "status", resp.Status)
		return errPermissionDenied
	}
	if !resp.Body {
		return errUnknown
	}
	if req.Key[len(req.Key)-1] == '/' {
		return nil
	}
	return s.admin.WriteValueType(ctx, req.Key, req.ValueType)
}

func (s *kvService) Update(ctx context.Context, key, value string) error {
	resp1, err := s.kv.Read(ctx, key)
	// resp1, err := s.client.KV().Keys(ctx, key, "", consul.QueryOptionsFromContext(ctx))
	if err != nil {
		return errFailedToConnectConsul
	}
	if resp1.Body == nil {
		return &DomainError{Code: DomainErrorCodeNotFound, Message: "key not found"}
	}

	resp, err := s.kv.Write(ctx, key, value)
	if err != nil {
		return errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return errPermissionDenied
	}
	if !resp.Body {
		return errUnknown
	}
	return nil
}

func (s *kvService) UpdateType(ctx context.Context, key, valueType string) error {
	resp1, err := s.kv.ListKeys(ctx, key, "")
	if err != nil {
		slog.Error("kvUpdateType: failed to list keys", "key", key, "error", err)
		return errFailedToConnectConsul
	}
	if len(resp1.Body) == 0 {
		return &DomainError{Code: DomainErrorCodeNotFound, Message: "key not found"}
	}
	if key[len(key)-1] == '/' {
		return nil
	}
	return s.admin.WriteValueType(ctx, key, valueType)
}

type BatchUpdateErrorList struct {
	Key   string
	Error error
}

func (s *kvService) BatchUpdate(ctx context.Context, req *BatchUpdateRequest) error {
	nErr, errList := 0, []BatchUpdateErrorList{}
	for _, kv := range req.KeyValues {
		err := s.Update(ctx, kv.Key, kv.Value)
		if err != nil {
			nErr++
			errList = append(errList, BatchUpdateErrorList{kv.Key, err})
		}
	}
	if nErr == 0 {
		return nil
	}
	var b buffer.Buffer
	b.WriteInt(nErr).WriteString(" errors occured during batch update:")
	for i := 0; i < len(errList); i++ {
		b.WriteByte(' ').WriteString(errList[i].Key).WriteString(": ").WriteString(errList[i].Error.Error())
		if i == len(errList)-1 {
			b.WriteByte('.')
		} else {
			b.WriteByte(';')
		}
	}
	return &DomainError{
		Code:    DomainErrorCodeMultiple,
		Message: b.String(),
	}
}

func (s *kvService) Delete(ctx context.Context, key string) error {
	resp, err := s.kv.Delete(ctx, key)
	if err != nil {
		slog.Error("kvDelete: failed to delete key", "key", key, "error", err)
		return errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		slog.Error("kvDelete: permission denied", "key", key, "status", resp.Status)
		return errPermissionDenied
	}
	s.admin.DeleteValueType(ctx, key)
	return nil
}

func (s *kvService) WatchOpenNotificationsCount(ctx context.Context, cb func(n int)) {
	s.kv.WatchKeys(ctx, ConseeInternalKeyPrefix+"notifications/open/", func(r *consul.Response[[]string], err error) (stop bool) {
		if cb != nil && r != nil {
			slog.Info("open notifications count changed", "n", len(r.Body))
			cb(len(r.Body))
		}
		return false
	})
}
