// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	. "github.com/FlyingOnion/consee/backend/common"
	"github.com/FlyingOnion/consee/backend/repo"
)

var _ AdminService = &adminService{}

type AdminService interface {
	// AdminRepo is used for admin operations.
	// Other than initializations, it should not be used outside service package.
	AdminRepo() repo.AdminRepo

	GetValueType(ctx context.Context, b64key string) (string, error)
	WriteValueType(ctx context.Context, b64key, vt string) error
	DeleteValueType(ctx context.Context, b64key string) error
	GetKVHistory(ctx context.Context, b64key string) ([]string, error)
	AddNewHistoryVersion(ctx context.Context, b64key, version, oldValue string) error
	GetKVHistoryValue(ctx context.Context, b64key, version string) (string, error)
	// GetKVMeta(ctx context.Context, key string) (*KVMeta, error)
	// WriteKVMeta(ctx context.Context, key string, meta *KVMeta) error
	// DeleteKVMeta(ctx context.Context, key string) error

	ListNotifications(ctx context.Context) (*ListNotificationsResponse, error)
	GetOpenNotificationsCount(ctx context.Context) (int, error)
	WriteNotification(ctx context.Context, n *Notification) error

	// CheckAdmin(ctx context.Context, token string) error
	GetTokenMetadata(ctx context.Context, accessorId string) (*TokenMetadata, error)
	GetTokenName(ctx context.Context, accessorId string) (string, error)
	GetTokenIdByName(ctx context.Context, name string) (accessorId string, err error)
	// WriteIdNameMapping writes both id-name and name-id mapping
	WriteIdNameMapping(ctx context.Context, accessorId, name string) error
	WriteTokenMetadata(ctx context.Context, accessorId string, metadata *TokenMetadata) error
	DeleteTokenMetadata(ctx context.Context, accessorId, name string) error
}

type adminService struct {
	admin repo.AdminRepo
}

func NewAdminService(admin repo.AdminRepo) AdminService {
	return &adminService{admin}
}

func (a *adminService) AdminRepo() repo.AdminRepo {
	return a.admin
}

func (a *adminService) GetValueType(ctx context.Context, b64key string) (string, error) {
	resp, err := a.admin.Read(ctx, ConseeInternalKeyPrefix+"kvmeta/valuetype/"+b64key)
	if err != nil {
		slog.Error("failed to get value type", "b64key", b64key, "error", err)
		return "", errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return "", errAdminPermissionDenied
	}
	if resp.Body != nil && len(resp.Body.Value) > 0 {
		return string(resp.Body.Value), nil
	}
	return "", nil
}

func (a *adminService) WriteValueType(ctx context.Context, b64key, vt string) error {
	if vt == "" {
		vt = "plaintext"
	}
	resp, err := a.admin.Write(ctx, ConseeInternalKeyPrefix+"kvmeta/valuetype/"+b64key, vt)
	if err != nil {
		slog.Error("failed to write value type", "b64key", b64key, "valueType", vt, "error", err)
		return errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return errAdminPermissionDenied
	}
	return nil
}

func (a *adminService) DeleteValueType(ctx context.Context, b64key string) error {
	resp, err := a.admin.Delete(ctx, ConseeInternalKeyPrefix+"kvmeta/valuetype/"+b64key)
	if err != nil {
		slog.Error("failed to delete value type", "b64key", b64key, "error", err)
		return errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return errAdminPermissionDenied
	}
	return nil
}

func (a *adminService) GetKVHistory(ctx context.Context, b64key string) ([]string, error) {
	return []string{}, errNotImplemented
}

func (a *adminService) AddNewHistoryVersion(ctx context.Context, b64key, version, oldValue string) error {
	return errNotImplemented
}

func (a *adminService) GetKVHistoryValue(ctx context.Context, b64key, version string) (string, error) {
	return "", errNotImplemented
}

func (a *adminService) ListNotifications(ctx context.Context) (*ListNotificationsResponse, error) {
	return nil, errNotImplemented
}

func (a *adminService) GetOpenNotificationsCount(ctx context.Context) (int, error) {
	return 0, errNotImplemented
}

func (a *adminService) WriteNotification(ctx context.Context, n *Notification) error {
	return errNotImplemented
}

func (a *adminService) GetTokenMetadata(ctx context.Context, accessorId string) (*TokenMetadata, error) {
	resp, err := a.admin.Read(ctx, ConseeInternalKeyPrefix+"acl-token/metadata/"+accessorId)
	if err != nil {
		slog.Error("failed to get token metadata", "accessorId", accessorId, "error", err)
		return nil, errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return nil, errAdminPermissionDenied
	}
	if resp.Status == http.StatusNotFound {
		return nil, &DomainError{Code: DomainErrorCodeInternalError, Message: "metadata not found"}
	}
	if resp.Err != nil {
		slog.Error("failed to parse token metadata response", "accessorId", accessorId, "error", resp.Err)
		return nil, errFailedToParse
	}
	var meta TokenMetadata
	err = json.Unmarshal(resp.Body.Value, &meta)
	if err != nil {
		slog.Error("failed to unmarshal token metadata", "accessorId", accessorId, "error", err)
		return nil, errFailedToParse
	}
	return &meta, nil
}

func (a *adminService) GetTokenIdByName(ctx context.Context, name string) (accessorId string, err error) {
	resp, err := a.admin.Read(ctx, ConseeInternalKeyPrefix+"acl-token/name-id/"+name)
	if err != nil {
		slog.Error("failed to get token id by name", "name", name, "error", err)
		return "", errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return "", errAdminPermissionDenied
	}
	if resp.Status == http.StatusNotFound {
		return "", &DomainError{Code: DomainErrorCodeNotFound, Message: "token not found"}
	}
	if resp.Err != nil {
		slog.Error("failed to parse token id response", "name", name, "error", resp.Err)
		return "", errFailedToParse
	}
	if resp.Body == nil || resp.Body.Value == nil {
		return "", nil
	}
	return string(resp.Body.Value), nil
}

func (a *adminService) GetTokenName(ctx context.Context, id string) (string, error) {
	resp, err := a.admin.Read(ctx, ConseeInternalKeyPrefix+"acl-token/id-name/"+id)
	if err != nil {
		slog.Error("failed to get token name by id", "id", id, "error", err)
		return "", errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return "", errAdminPermissionDenied
	}
	if resp.Status == http.StatusNotFound {
		return "", &DomainError{Code: DomainErrorCodeNotFound, Message: "token not found"}
	}
	if resp.Err != nil {
		slog.Error("failed to parse token name response", "id", id, "error", resp.Err)
		return "", errFailedToParse
	}
	if resp.Body == nil || resp.Body.Value == nil {
		return "", nil
	}
	return string(resp.Body.Value), nil
}

func (a *adminService) WriteIdNameMapping(ctx context.Context, id, name string) (err error) {
	err = a.writeIdNameMapping(ctx, id, name)
	if err != nil {
		return
	}
	return a.writeNameIdMapping(ctx, name, id)
}

func (a *adminService) writeIdNameMapping(ctx context.Context, id, name string) error {
	resp, err := a.admin.Write(ctx, ConseeInternalKeyPrefix+"acl-token/id-name/"+id, name)
	if err != nil {
		slog.Error("failed to write id-name mapping", "id", id, "name", name, "error", err)
		return errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return errAdminPermissionDenied
	}
	if resp.Err != nil {
		slog.Error("failed to parse id-name mapping response", "id", id, "name", name, "error", resp.Err)
		return errFailedToParse
	}
	return nil
}

func (a *adminService) writeNameIdMapping(ctx context.Context, name, id string) error {
	resp, err := a.admin.Write(ctx, ConseeInternalKeyPrefix+"acl-token/name-id/"+name, id)
	if err != nil {
		slog.Error("failed to write name-id mapping", "name", name, "id", id, "error", err)
		return errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return errAdminPermissionDenied
	}
	if resp.Err != nil {
		slog.Error("failed to parse name-id mapping response", "name", name, "id", id, "error", resp.Err)
		return errFailedToParse
	}
	return nil
}

func (a *adminService) WriteTokenMetadata(ctx context.Context, id string, metadata *TokenMetadata) error {
	b, _ := metadata.MarshalJSON()
	resp, err := a.admin.Write(ctx, ConseeInternalKeyPrefix+"acl-token/metadata/"+id, string(b))
	if err != nil {
		slog.Error("failed to write token metadata", "id", id, "error", err)
		return errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return errAdminPermissionDenied
	}
	if resp.Err != nil {
		slog.Error("failed to parse token metadata response", "id", id, "error", resp.Err)
		return errFailedToParse
	}
	return nil
}

func (a *adminService) DeleteTokenMetadata(ctx context.Context, id, name string) error {
	a.admin.Delete(ctx, ConseeInternalKeyPrefix+"acl-token/id-name/"+id)
	a.admin.Delete(ctx, ConseeInternalKeyPrefix+"acl-token/name-id/"+name)
	a.admin.Delete(ctx, ConseeInternalKeyPrefix+"acl-token/metadata/"+id)
	return nil
}
