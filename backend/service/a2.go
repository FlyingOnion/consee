// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package service

import (
	"archive/zip"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/FlyingOnion/consee/backend/buffer"
	. "github.com/FlyingOnion/consee/backend/common"
)

type All interface {
	Initialize(ctx context.Context) error
	Import(ctx context.Context, req *ImportRequest) (*ImportResponse, error)
	Export(ctx context.Context, req *ExportRequest) ([]byte, error)
}

type a2 struct {
	kv    KVService
	acl   ACLService
	admin AdminService
}

func NewA2(kv KVService, acl ACLService, admin AdminService) All {
	return &a2{kv, acl, admin}
}

func (s *a2) Initialize(ctx context.Context) (err error) {
	resp0, err := s.admin.AdminRepo().ReadSelf(ctx)
	if err != nil {
		slog.Error("failed to read self during initialization", "error", err)
		return
	}
	if resp0.Status == http.StatusForbidden {
		return errAdminPermissionDenied
	}
	if resp0.Status == http.StatusNotFound {
		return &DomainError{Code: DomainErrorCodeNotFound, Message: "token not found"}
	}
	if resp0.Err != nil {
		slog.Error("failed to parse self response during initialization", "error", resp0.Err)
		return errFailedToParse
	}
	for _, policy := range resp0.Body.Policies {
		if policy.Name == PolicyNameGlobalManagement {
			goto validateIdNameMapping
		}
	}
	return errNotAdmin

validateIdNameMapping:
	resp1, _ := s.admin.AdminRepo().Read(ctx, ConseeInternalKeyPrefix+"acl-token/id-name/"+resp0.Body.AccessorID)
	now := time.Now().Format(time.DateTime)
	// have checked permissions before
	// 403 is impossible here
	if resp1.Status == http.StatusNotFound || resp1.Body == nil {
		s.admin.WriteIdNameMapping(ctx, resp0.Body.AccessorID, ConseeAdmin)
		s.admin.WriteTokenMetadata(ctx, resp0.Body.AccessorID, &TokenMetadata{
			CreatedAt:     now,
			CreatedBy:     "initializer",
			LastUpdatedAt: now,
			LastUpdatedBy: "initializer",
			Version:       now,
		})
	}

	return
}

func (s *a2) Export(ctx context.Context, req *ExportRequest) (data []byte, err error) {
	if req.Format != "zip" {
		return nil, &DomainError{Code: DomainErrorCodeInvalidInput, Message: "unsupported format, only zip is supported"}
	}

	// 如果没有指定keys，则导出所有key
	keys := req.Keys
	if len(keys) == 0 {
		allKeys, err := s.kv.ListKeys(ctx)
		if err != nil {
			slog.Error("failed to list all keys during export", "error", err)
			return nil, err
		}
		keys = allKeys
	}

	// 创建zip文件
	var buf buffer.Buffer
	zipWriter := zip.NewWriter(&buf)
	defer func() {
		if err != nil {
			zipWriter.Close()
		}
	}()

	kvMeta := make([]*ExportedKVMeta, 0, len(keys))
	// 导出每个key的值
	for _, key := range keys {
		kv, err := s.kv.Get(ctx, key)
		// err := s.writeKVToZip(ctx, zipWriter, key)
		if err != nil {
			slog.Error("failed to get key during export", "key", key, "error", err)
			return nil, err
		}
		b64key := base64.StdEncoding.EncodeToString([]byte(key))
		f, err := zipWriter.Create("kv/" + b64key + "/latest")
		if err != nil {
			slog.Error("failed to create zip entry for key", "key", key, "b64key", b64key, "error", err)
			return nil, err
		}
		f.Write([]byte(kv.Value))

		vtkv, err := s.admin.GetValueType(ctx, b64key)
		if err != nil {
			dErr := err.(*DomainError)
			if dErr.Code == DomainErrorCodeNotFound {
				vtkv = "plaintext"
			} else {
				slog.Error("failed to get value type during export", "key", key, "b64key", b64key, "error", err)
				return nil, err
			}
		}

		historyKeys, err := s.admin.GetKVHistory(ctx, b64key)
		if err != nil {
			slog.Error("failed to get kv history during export", "key", key, "b64key", b64key, "error", err)
			return nil, errFailedToConnectConsul
		}
		for _, ht := range historyKeys {
			htvalue, err := s.admin.GetKVHistoryValue(ctx, b64key, ht)
			if err != nil {
				slog.Error("failed to get kv history value during export", "key", key, "b64key", b64key, "history", ht, "error", err)
				return nil, err
			}
			f, err := zipWriter.Create("kv/" + b64key + "/" + ht)
			if err != nil {
				slog.Error("failed to create zip entry for history", "key", key, "b64key", b64key, "history", ht, "error", err)
				return nil, err
			}
			f.Write([]byte(htvalue))
		}
		kvMeta = append(kvMeta, &ExportedKVMeta{
			Name:            key,
			ValueType:       vtkv,
			HistoryVersions: historyKeys,
		})
	}

	e := ExportMetadata{
		Keys:     kvMeta,
		Tokens:   []ACLLink{},
		Policies: []string{},
	}

	if req.ACL {
		tokens, _ := s.acl.ListTokens(ctx)
		for _, t := range tokens {
			token, err := s.acl.ReadToken(ctx, t.ID)
			if err != nil {
				slog.Error("failed to read token during export", "tokenId", t.ID, "tokenName", t.Name, "error", err)
				return nil, err
			}
			policyMode := "common"
			rules := ""
			policies := make([]string, 0, len(token.Policies))
			if len(token.Policies) == 1 && token.Policies[0].Name == "--"+token.AccessorID {
				// exclusive，连带专有策略的规则一起保存
				policyMode = "exclusive"
				policy, err := s.acl.ReadPolicy(ctx, token.Policies[0].ID)
				if err != nil {
					return nil, err
				}
				rules = policy.Rules
			} else {
				// common token，需要保存所有策略名字
				for _, p := range token.Policies {
					policies = append(policies, p.Name)
				}
			}
			f, err := zipWriter.Create("tokens/" + t.ID)
			if err != nil {
				slog.Error("failed to create zip entry for token", "tokenId", t.ID, "tokenName", t.Name, "error", err)
				return nil, err
			}
			b, _ := json.Marshal(CreateTokenRequest{
				AccessorID: token.AccessorID,
				SecretID:   token.SecretID,
				Name:       token.Name,
				PolicyMode: policyMode,
				Policies:   policies,
				Rules:      rules,
			})
			f.Write(b)
		}
		e.Tokens = tokens

		policies, _ := s.acl.ListPolicies(ctx, ListPoliciesOptions{Exclusive: "0"})
		policyNames := make([]string, 0, len(policies))
		for _, p := range policies {
			if p.Name == PolicyNameGlobalManagement || p.Name == PolicyNameBuiltinGlobalReadonly {
				continue
			}
			policy, err := s.acl.ReadPolicy(ctx, p.ID)
			if err != nil {
				slog.Error("failed to read policy during export", "policyId", p.ID, "policyName", p.Name, "error", err)
				return nil, err
			}
			b64PolicyName := base64.StdEncoding.EncodeToString([]byte(p.Name))
			f, err := zipWriter.Create("policies/" + b64PolicyName)
			if err != nil {
				slog.Error("failed to create zip entry for policy", "policyId", p.ID, "policyName", p.Name, "b64PolicyName", b64PolicyName, "error", err)
				return nil, err
			}
			b, _ := json.Marshal(CreatePolicyRequest{
				Name:        policy.Name,
				Description: policy.Description,
				Rules:       policy.Rules,
			})
			f.Write(b)
			policyNames = append(policyNames, p.Name)
		}
		e.Policies = policyNames

	}

	b, _ := json.Marshal(e)
	w, err := zipWriter.Create("metadata.json")
	if err != nil {
		slog.Error("failed to create zip entry for metadata", "error", err)
		return nil, err
	}
	w.Write(b)
	zipWriter.Flush()
	zipWriter.Close()
	return buf.Bytes(), nil
}

func (s *a2) Import(ctx context.Context, req *ImportRequest) (*ImportResponse, error) {
	if binary.LittleEndian.Uint32(req.FileContent[:4]) != 0x04034b50 {
		return nil, &DomainError{Code: DomainErrorCodeInvalidInput, Message: "invalid file format"}
	}
	r, err := zip.NewReader(bytes.NewReader(req.FileContent), int64(len(req.FileContent)))
	if err != nil {
		slog.Error("failed to create zip reader for import", "error", err)
		return nil, errUnknown
	}
	// read metadata.json
	f, err := r.Open("metadata.json")
	if err != nil {
		slog.Error("failed to open metadata.json during import", "error", err)
		return nil, &DomainError{Code: DomainErrorCodeInvalidInput, Message: "invalid file format: metadata.json not found"}
	}
	var meta ExportMetadata
	err = json.NewDecoder(f).Decode(&meta)
	if err != nil {
		slog.Error("failed to decode metadata.json during import", "error", err)
		return nil, &DomainError{Code: DomainErrorCodeInvalidInput, Message: "invalid file format: metadata.json is invalid"}
	}
	f.Close()

	resp := s.ImportDryrun(ctx, &meta)
	if !req.Dryrun {
		return resp, nil
	}

	// 实际导入数据
	return s.doImport(ctx, r, &meta, req.OnConflict), nil
}

func (s *a2) doImport(ctx context.Context, r *zip.Reader, meta *ExportMetadata, onConflict OnConflictPolicy) *ImportResponse {
	resp := &ImportResponse{
		Conflicts: Conflicts{
			Keys:     []*ExportedKVMeta{},
			Tokens:   []ACLLink{},
			Policies: []string{},
		},
		Errors: []ImportError{},
	}
	// 导入KV数据
	for _, kv := range meta.Keys {
		conflict := false
		b64key := base64.StdEncoding.EncodeToString([]byte(kv.Name))
		// 读取最新的值
		f, err := r.Open("kv/" + b64key + "/latest")
		if err != nil {
			resp.Errors = append(resp.Errors, ImportError{
				Kind:  "kv",
				Param: kv.Name,
				Cause: err.Error(),
			})
			continue
		}

		valueBytes, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			resp.Errors = append(resp.Errors, ImportError{
				Kind:  "kv",
				Param: kv.Name,
				Cause: err.Error(),
			})
			continue
		}

		// 创建或更新KV
		existingKV, _ := s.kv.Get(ctx, kv.Name)
		if existingKV == nil {
			// 如果key不存在，创建新的
			err = s.kv.Create(ctx, &CreateKeyValueRequest{
				Key:       kv.Name,
				Value:     string(valueBytes),
				ValueType: kv.ValueType,
			})
			if err != nil {
				resp.Errors = append(resp.Errors, ImportError{
					Kind:  "kv",
					Param: kv.Name,
					Cause: err.Error(),
				})
			}
			goto importHistory
		}
		// 值一样时直接跳过更新，否则记录冲突
		if existingKV.Value != string(valueBytes) {
			conflict = true
			if onConflict == OnConflictPolicyReplace {
				// 如果key存在，更新值
				err = s.kv.Update(ctx, kv.Name, string(valueBytes))
				if err != nil {
					resp.Errors = append(resp.Errors, ImportError{
						Kind:  "kv",
						Param: kv.Name,
						Cause: err.Error(),
					})
				}
				s.admin.WriteValueType(ctx, b64key, kv.ValueType)
			}
		}
	importHistory:
		// 导入历史版本（如果有）
		for _, history := range kv.HistoryVersions {
			hf, err := r.Open("kv/" + b64key + "/" + history)
			if err != nil {
				resp.Errors = append(resp.Errors, ImportError{
					Kind:  "kv-history",
					Param: kv.Name + ":" + history,
					Cause: "history version not found",
				})
				continue
			}
			historyValue, err := io.ReadAll(hf)
			hf.Close()
			if err != nil {
				resp.Errors = append(resp.Errors, ImportError{
					Kind:  "kv-history",
					Param: kv.Name + ":" + history,
					Cause: "failed to read history version",
				})
				continue
			}

			// 保存历史版本到admin服务
			s.admin.AddNewHistoryVersion(ctx,
				b64key,
				history,
				string(historyValue),
			)
		}
		if conflict {
			resp.Conflicts.Keys = append(resp.Conflicts.Keys, kv)
		} else {
			resp.Successes.Keys++
		}
	}

	// 导入policies
	for _, policyName := range meta.Policies {
		if policyName == PolicyNameGlobalManagement || policyName == PolicyNameBuiltinGlobalReadonly {
			// 应该是不可能触发的，导出时已经排除了内置策略
			continue // 跳过内置策略
		}

		b64PolicyName := base64.StdEncoding.EncodeToString([]byte(policyName))
		f, err := r.Open("policies/" + b64PolicyName)
		if err != nil {
			resp.Errors = append(resp.Errors, ImportError{
				Kind:  "policy",
				Param: policyName,
				Cause: "policy not found",
			})
			continue // 跳过找不到的策略
		}

		var policyReq CreatePolicyRequest
		err = json.NewDecoder(f).Decode(&policyReq)
		f.Close()
		if err != nil {
			resp.Errors = append(resp.Errors, ImportError{
				Kind:  "policy",
				Param: policyName,
				Cause: "invalid policy information",
			})
			continue
		}

		conflict := false
		// 检查策略是否已存在
		existingPolicy, _ := s.acl.ReadPolicy(ctx, policyName)
		if existingPolicy == nil {
			// 创建新策略
			err = s.acl.CreatePolicy(ctx, &policyReq)
			if err != nil {
				resp.Errors = append(resp.Errors, ImportError{
					Kind:  "policy",
					Param: policyName,
					Cause: err.Error(),
				})
			}
		} else {
			conflict = true
			// 更新现有策略的规则
			err = s.acl.UpdatePolicyRule(ctx, policyName, policyReq.Rules)
			if err != nil {
				resp.Errors = append(resp.Errors, ImportError{
					Kind:  "policy",
					Param: policyName,
					Cause: err.Error(),
				})
			}
		}
		if conflict {
			resp.Conflicts.Policies = append(resp.Conflicts.Policies, policyName)
		} else {
			resp.Successes.Policies++
		}
	}

	// 导入tokens
	for _, token := range meta.Tokens {
		f, err := r.Open("tokens/" + token.ID)
		if err != nil {
			resp.Errors = append(resp.Errors, ImportError{
				Kind:  "token",
				Param: token.Name + "(ID:" + token.ID + ")",
				Cause: "token not found",
			})
			continue // 跳过找不到的token
		}

		var tokenReq CreateTokenRequest
		err = json.NewDecoder(f).Decode(&tokenReq)
		f.Close()
		if err != nil {
			resp.Errors = append(resp.Errors, ImportError{
				Kind:  "token",
				Param: token.Name + "(ID:" + token.ID + ")",
				Cause: "invalid token information",
			})
			continue
		}

		conflict := false
		// 检查token是否已存在
		existingToken, _ := s.acl.ReadToken(ctx, token.ID)
		if existingToken == nil {
			// 创建新token
			err = s.acl.CreateToken(ctx, &tokenReq)
			if err != nil {
				resp.Errors = append(resp.Errors, ImportError{
					Kind:  "token",
					Param: token.Name + "(ID:" + token.ID + ")",
					Cause: err.Error(),
				})
			}
		} else {
			conflict = true
			// 更新现有token
			err = s.acl.UpdateToken(ctx, token.ID, &UpdateTokenRequest{
				Policies: tokenReq.Policies,
			})
			if err != nil {
				resp.Errors = append(resp.Errors, ImportError{
					Kind:  "token",
					Param: token.Name + "(ID:" + token.ID + ")",
					Cause: err.Error(),
				})
			}
		}
		if conflict {
			resp.Conflicts.Tokens = append(resp.Conflicts.Tokens, token)
		} else {
			resp.Successes.Tokens++
		}
	}

	return resp
}

func (s *a2) ImportDryrun(ctx context.Context, meta *ExportMetadata) *ImportResponse {
	resp := &ImportResponse{
		Conflicts: Conflicts{
			Keys:     []*ExportedKVMeta{},
			Tokens:   []ACLLink{},
			Policies: []string{},
		},
		Errors: []ImportError{},
	}

	// 检查KV数据冲突
	for _, kv := range meta.Keys {
		existingKV, err := s.kv.Get(ctx, kv.Name)
		if err == nil && existingKV != nil {
			// Key已存在，记录冲突
			resp.Conflicts.Keys = append(resp.Conflicts.Keys, kv)
			continue
		}
		dErr := err.(*DomainError)
		if dErr.Code == DomainErrorCodeNotFound {
			// Key不存在，可以安全导入
			resp.Successes.Keys++
			continue
		}
		resp.Errors = append(resp.Errors, ImportError{
			Kind:  "kv",
			Param: kv.Name,
			Cause: dErr.Message,
		})
	}

	// 检查policies冲突
	for _, policyName := range meta.Policies {
		if policyName == PolicyNameGlobalManagement || policyName == PolicyNameBuiltinGlobalReadonly {
			// 应该是不可能触发的，导出时已经排除了内置策略
			continue // 跳过内置策略
		}

		existingPolicy, err := s.acl.ReadPolicy(ctx, policyName)
		if err == nil && existingPolicy != nil {
			// Policy已存在，记录冲突
			resp.Conflicts.Policies = append(resp.Conflicts.Policies, policyName)
			continue
		}
		dErr := err.(*DomainError)
		if dErr.Code == DomainErrorCodeNotFound {
			// Policy不存在，可以安全导入
			resp.Successes.Policies++
			continue
		}
		resp.Errors = append(resp.Errors, ImportError{
			Kind:  "policy",
			Param: policyName,
			Cause: dErr.Message,
		})
	}

	// 检查tokens冲突
	for _, token := range meta.Tokens {
		existingToken, err := s.acl.ReadToken(ctx, token.ID)

		if err == nil && existingToken != nil {
			// Token已存在，记录冲突
			resp.Conflicts.Tokens = append(resp.Conflicts.Tokens, token)
			continue
		}
		dErr := err.(*DomainError)
		if dErr.Code == DomainErrorCodeNotFound {
			// Token不存在，可以安全导入
			resp.Successes.Tokens++
		}
		resp.Errors = append(resp.Errors, ImportError{
			Kind:  "token",
			Param: token.Name + "(ID:" + token.ID + ")",
			Cause: dErr.Message,
		})
	}

	return resp
}
