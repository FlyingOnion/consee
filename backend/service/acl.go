// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package service

import (
	"context"
	"log/slog"
	"net/http"
	"regexp"
	"slices"
	"time"

	. "github.com/FlyingOnion/consee/backend/common"
	"github.com/FlyingOnion/consee/backend/consul"
	"github.com/FlyingOnion/consee/backend/repo"
	"github.com/google/uuid"
)

var ConseeExclusivePolicyNameRegexp = regexp.MustCompile(`^--[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

type ACLService interface {
	ValidateToken(ctx context.Context) error
	// CheckAdmin checks whether the context (with token) is valid and has admin permission.
	// Returns nil if the token is valid and has admin permission.
	// Returns a non-nil DomainError otherwise.
	CheckAdmin(ctx context.Context) error

	CreateTokenApplicationRequest(ctx context.Context, req *TokenApplicationRequest) (*TokenApplicationResponse, error)
	ReviewTokenApplicationRequest(ctx context.Context, id string, req *HandleTokenApplicationRequest) error

	ListTokens(ctx context.Context) ([]ACLLink, error)
	ReadToken(ctx context.Context, id string) (*ReadTokenResponse, error)
	CreateToken(ctx context.Context, req *CreateTokenRequest) error
	UpdateToken(ctx context.Context, id string, req *UpdateTokenRequest) error
	DeleteToken(ctx context.Context, id string) error

	ValidateHCLRules(rules string) *ValidateHCLRulesResponse
	ListPolicies(ctx context.Context, options ListPoliciesOptions) ([]ACLLink, error)
	CreatePolicy(ctx context.Context, req *CreatePolicyRequest) error
	ReadPolicy(ctx context.Context, name string) (*ReadPolicyResponse, error)
	UpdatePolicyRule(ctx context.Context, name, rules string) error
	DeletePolicy(ctx context.Context, name string) error
}

type aclService struct {
	acl   repo.ACLRepo
	admin AdminService
}

func NewACLService(acl repo.ACLRepo, admin AdminService) ACLService {
	return &aclService{
		acl:   acl,
		admin: admin,
	}
}

func (s *aclService) ValidateToken(ctx context.Context) error {
	resp, err := s.acl.ReadSelf(ctx)
	if err != nil {
		slog.Error("failed to read self during token validation", "error", err)
		return errFailedToConnectConsul
	}
	if resp.Status == http.StatusNotFound {
		return &DomainError{Code: DomainErrorCodeNotFound, Message: "token not found"}
	}
	if resp.Err != nil {
		slog.Error("failed to parse self response during token validation", "error", resp.Err)
		return errFailedToParse
	}
	return nil
}

func (s *aclService) CheckAdmin(ctx context.Context) error {
	resp, err := s.acl.ReadSelf(ctx)
	if err != nil {
		slog.Error("failed to read self during admin check", "error", err)
		return errFailedToConnectConsul
	}
	if resp.Status == http.StatusNotFound {
		return &DomainError{Code: DomainErrorCodeNotFound, Message: "token not found"}
	}
	if resp.Err != nil {
		slog.Error("failed to parse self response during admin check", "error", resp.Err)
		return errFailedToParse
	}
	for _, policy := range resp.Body.Policies {
		if policy.Name == PolicyNameGlobalManagement {
			return nil
		}
	}
	return errPermissionDenied
}

func (s *aclService) ListTokens(ctx context.Context) ([]ACLLink, error) {
	resp, err := s.admin.AdminRepo().List(ctx, ConseeInternalKeyPrefix+"acl-token/id-name/")
	if err != nil {
		slog.Error("failed to list tokens", "error", err)
		return nil, errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return nil, errAdminPermissionDenied
	}
	if resp.Err != nil {
		slog.Error("failed to parse token list response", "error", resp.Err)
		return nil, errFailedToParse
	}
	tokenList := make([]ACLLink, 0, len(resp.Body))
	for _, kvp := range resp.Body {
		tokenList = append(tokenList, ACLLink{
			ID:   kvp.Key[len(ConseeInternalKeyPrefix+"acl-token/id-name/"):],
			Name: string(kvp.Value),
		})
	}
	return tokenList, nil
}

func (a *aclService) listPolicyTokens(ctx context.Context, policyId string) ([]ACLLink, error) {
	resp, err := a.acl.ListTokensFiltered(ctx, consul.ACLTokenFilterOptions{Policy: policyId})
	if err != nil {
		slog.Error("failed to list policy tokens", "policyId", policyId, "error", err)
		return []ACLLink{}, errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return []ACLLink{}, errPermissionDenied
	}
	if resp.Err != nil {
		slog.Error("failed to parse policy tokens response", "policyId", policyId, "error", resp.Err)
		return []ACLLink{}, errFailedToParse
	}
	tokenList := make([]ACLLink, 0, len(resp.Body))
	for _, t := range resp.Body {
		name, err := a.admin.GetTokenName(ctx, t.AccessorID)
		if err != nil {
			slog.Error("failed to get token name during policy token listing", "tokenId", t.AccessorID, "policyId", policyId, "error", err)
			return []ACLLink{}, err
		}
		tokenList = append(tokenList, ACLLink{ID: t.AccessorID, Name: name})
	}
	return tokenList, nil
}

func (s *aclService) ReadToken(ctx context.Context, id string) (*ReadTokenResponse, error) {
	resp, err := s.acl.ReadToken(ctx, id)
	if err != nil {
		slog.Error("failed to read token", "tokenId", id, "error", err)
		return nil, errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return nil, errPermissionDenied
	}
	if resp.Status == http.StatusNotFound {
		return nil, &DomainError{Code: DomainErrorCodeNotFound, Message: "token not found"}
	}
	if resp.Err != nil {
		slog.Error("failed to parse token list response", "error", resp.Err)
		return nil, errFailedToParse
	}
	policies := make([]ACLLink, 0, len(resp.Body.Policies))
	for _, p := range resp.Body.Policies {
		policies = append(policies, ACLLink{ID: p.ID, Name: p.Name})
	}
	roles := make([]ACLLink, 0, len(resp.Body.Roles))
	for _, r := range resp.Body.Roles {
		roles = append(roles, ACLLink{ID: r.ID, Name: r.Name})
	}
	name, err := s.admin.GetTokenName(ctx, id)
	if err != nil {
		slog.Error("failed to get token name", "tokenId", id, "error", err)
		return nil, err
	}
	metadata, err := s.admin.GetTokenMetadata(ctx, id)
	if err != nil {
		slog.Error("failed to get token metadata", "tokenId", id, "error", err)
		return nil, err
	}

	return &ReadTokenResponse{
		AccessorID: resp.Body.AccessorID,
		SecretID:   resp.Body.SecretID,
		Policies:   policies,
		Roles:      roles,
		Name:       name,
		Metadata:   metadata,
	}, nil
}

func (s *aclService) CreateTokenApplicationRequest(ctx context.Context, req *TokenApplicationRequest) (*TokenApplicationResponse, error) {
	return nil, errNotImplemented
}

func (s *aclService) ReviewTokenApplicationRequest(ctx context.Context, id string, req *HandleTokenApplicationRequest) error {
	return errNotImplemented
}

func (s *aclService) CreateToken(ctx context.Context, req *CreateTokenRequest) (err error) {
	// Validations first
	if req.AccessorID != "" {
		name, _ := s.admin.GetTokenName(ctx, req.AccessorID)
		if name != "" {
			return &DomainError{Code: DomainErrorCodeAlreadyExists, Message: "token accessor id already exists"}
		}
	}
	if req.Name != "" {
		accessorId, _ := s.admin.GetTokenIdByName(ctx, req.Name)
		if accessorId != "" {
			return &DomainError{Code: DomainErrorCodeAlreadyExists, Message: "token name already exists"}
		}
	}
	if req.SecretID != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		resp, err := s.acl.ReadSelf(consul.ContextWithQueryOptions(ctx, &consul.QueryOptions{Token: req.SecretID}))
		if err != nil {
			return errFailedToConnectConsul
		}
		if resp.Status == http.StatusOK || resp.Body != nil {
			return &DomainError{Code: DomainErrorCodeAlreadyExists, Message: "token secret id already exists"}
		}
	}

	// validate policy mode
	switch req.PolicyMode {
	case "", "common", "exclusive":
	default:
		return &DomainError{Code: DomainErrorCodeInvalidInput, Message: "invalid policy mode"}
	}

	// prepare creator info
	resp1, err := s.acl.ReadSelf(ctx)
	if err != nil {
		return errFailedToConnectConsul
	}
	if resp1.Status != http.StatusOK {
		err = errUnknown
		return
	}
	creatorToken := resp1.Body
	creator, _ := s.admin.GetTokenName(ctx, creatorToken.AccessorID)
	if creator == "" {
		creator = "unknown"
	}
	creator = creatorToken.AccessorID + " (" + creator + ")"

	// for uuid.Must panic handling, but actually it rarely panics
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		if err1, ok := r.(error); ok {
			err = &DomainError{Code: DomainErrorCodeInternalError, Message: err1.Error()}
			return
		}
		err = errUnknown
	}()

	accessorId := req.AccessorID
	if accessorId == "" {
		accessorId = uuid.Must(uuid.NewV7()).String()
	}
	secretId := req.SecretID
	if secretId == "" {
		secretId = uuid.Must(uuid.NewV7()).String()
	}

	defer func() {
		if err != nil {
			return
		}
		tokenName := req.Name
		if tokenName == "" {
			tokenName = "consee-token-" + accessorId
		}
		now := time.Now().Format(time.DateTime)
		err = s.admin.WriteIdNameMapping(ctx, accessorId, tokenName)
		if err != nil {
			return
		}
		err = s.admin.WriteTokenMetadata(ctx, accessorId, &TokenMetadata{
			CreatedAt:     now,
			CreatedBy:     creator,
			LastUpdatedAt: now,
			LastUpdatedBy: creator,
			Version:       now,
		})
		if err != nil {
			return
		}
	}()

	if req.PolicyMode == "exclusive" {
		return s.createTokenWithExclusivePolicy(ctx, &CreateTokenRequest2{
			AccessorID: accessorId,
			SecretID:   secretId,
			Rules:      req.Rules,
		})
	}
	for _, policy := range req.Policies {
		p, err := s.ReadPolicy(ctx, policy)
		if err != nil {
			return err
		}
		if ConseeExclusivePolicyNameRegexp.MatchString(p.Name) {
			return &DomainError{Code: DomainErrorCodeInvalidInput, Message: "policy " + p.Name + " is exclusive"}
		}
	}
	return s.createCommonToken(ctx, &CreateTokenRequest1{
		AccessorID: accessorId,
		SecretID:   secretId,
		Policies:   req.Policies,
		Roles:      req.Roles,
	})

	// resp, err := s.acl.CreateToken(ctx, makeTokenFromRequest(req))
	// if err != nil {
	// 	return errFailedToConnectConsul
	// }
	// if resp.Status == http.StatusForbidden {
	// 	return errPermissionDenied
	// }
	// if resp.Status != http.StatusOK {
	// 	return errUnknown
	// }

	// token := resp.Body
	// tokenName := req.Name
	// if tokenName == "" {
	// 	tokenName = "consee-token-" + token.AccessorID
	// }

}

func (s *aclService) createCommonToken(ctx context.Context, req *CreateTokenRequest1) (err error) {
	policies := make([]*consul.ACLLink, 0, len(req.Policies))
	for _, policy := range req.Policies {
		policies = append(policies, &consul.ACLLink{Name: policy})
	}
	roles := make([]*consul.ACLLink, 0, len(req.Roles))
	for _, role := range req.Roles {
		roles = append(roles, &consul.ACLLink{ID: role})
	}
	resp, err := s.acl.CreateToken(ctx, &consul.ACLToken{
		AccessorID: req.AccessorID,
		SecretID:   req.SecretID,
		Policies:   policies,
		Roles:      roles,
	})
	if err != nil {
		return errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return errPermissionDenied
	}
	if resp.Status != http.StatusOK || resp.Body == nil {
		return errUnknown
	}
	return nil
}

func (s *aclService) createTokenWithExclusivePolicy(ctx context.Context, req *CreateTokenRequest2) (err error) {
	resp, err := s.acl.CreatePolicy(ctx, &consul.ACLPolicy{
		Name:        "--" + req.AccessorID,
		Rules:       req.Rules,
		Description: "exclusive policy of token" + req.AccessorID,
	})
	if err != nil {
		return errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return errPermissionDenied
	}
	if resp.Status != http.StatusOK || resp.Body == nil {
		return errUnknown
	}
	policy := resp.Body
	resp2, err := s.acl.CreateToken(ctx, &consul.ACLToken{
		AccessorID: req.AccessorID,
		SecretID:   req.SecretID,
		Policies: []*consul.ACLLink{
			{ID: policy.ID},
		},
	})
	if err != nil {
		return errFailedToConnectConsul
	}
	if resp2.Status == http.StatusForbidden {
		return errPermissionDenied
	}
	if resp2.Status != http.StatusOK || resp.Body == nil {
		return errUnknown
	}
	return nil
}

func (s *aclService) UpdateToken(ctx context.Context, id string, req *UpdateTokenRequest) error {
	resp, err := s.acl.ReadToken(ctx, id)
	if err != nil {
		return errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return errPermissionDenied
	}
	if resp.Status != http.StatusOK || resp.Body == nil {
		return errUnknown
	}
	token := resp.Body
	if len(token.Policies) == 1 && ConseeExclusivePolicyNameRegexp.MatchString(token.Policies[0].Name) {
		return &DomainError{Code: DomainErrorCodePermissionDenied, Message: "token has an exclusive policy"}
	}

	for _, policyName := range req.Policies {
		resp1, err := s.acl.ReadPolicyByName(ctx, policyName)
		if err != nil {
			return err
		}
		if ConseeExclusivePolicyNameRegexp.MatchString(resp1.Body.Name) {
			return &DomainError{Code: DomainErrorCodePermissionDenied, Message: "policy " + resp1.Body.Name + " is exclusive"}
		}
	}
	policies := make([]*consul.ACLLink, 0, len(req.Policies))
	for _, policyName := range req.Policies {
		policies = append(policies, &consul.ACLLink{Name: policyName})
	}
	roles := make([]*consul.ACLLink, 0, len(req.Roles))
	for _, roleId := range req.Roles {
		roles = append(roles, &consul.ACLLink{ID: roleId})
	}
	resp2, err := s.acl.UpdateToken(ctx, &consul.ACLToken{
		AccessorID: id,
		Policies:   policies,
		Roles:      roles,
	})
	if err != nil {
		return errFailedToConnectConsul
	}
	if resp2.Status == http.StatusForbidden {
		return errPermissionDenied
	}
	// TODO: make it as a conditional compilation function
	metadata, err := s.admin.GetTokenMetadata(ctx, id)
	if err != nil {
		return err
	}

	resp1, err := s.acl.ReadSelf(ctx)
	if err != nil {
		return errFailedToConnectConsul
	}
	if resp1.Status != http.StatusOK {
		return errUnknown
	}
	creatorToken := resp1.Body
	creator, _ := s.admin.GetTokenName(ctx, creatorToken.AccessorID)
	if creator == "" {
		creator = "unknown"
	}
	creator = creatorToken.AccessorID + " (" + creator + ")"

	now := time.Now().Format(time.DateTime)
	err = s.admin.WriteTokenMetadata(ctx, id, &TokenMetadata{
		CreatedAt:     metadata.CreatedAt,
		CreatedBy:     metadata.CreatedBy,
		From:          metadata.From,
		Version:       now,
		LastUpdatedAt: now,
		LastUpdatedBy: creator,
	})
	return err
}

func (s *aclService) DeleteToken(ctx context.Context, id string) error {
	readTokenResponse, err := s.ReadToken(ctx, id)
	if err != nil {
		return err
	}
	resp, err := s.acl.DeleteToken(ctx, id)
	if err != nil {
		return errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return errPermissionDenied
	}
	if resp.Status != http.StatusOK {
		return errUnknown
	}
	s.admin.DeleteTokenMetadata(ctx, readTokenResponse.AccessorID, readTokenResponse.Name)
	if len(readTokenResponse.Policies) == 1 && ConseeExclusivePolicyNameRegexp.MatchString(readTokenResponse.Policies[0].Name) {
		s.acl.DeletePolicy(ctx, readTokenResponse.Policies[0].ID)
	}
	return nil
}

func (s *aclService) ValidateHCLRules(rules string) *ValidateHCLRulesResponse {
	// hclRuleList, err := parseHCLRules(rules)
	// if err != nil {
	// 	return &ValidateHCLRulesResponse{Error: err.Error()}
	// }
	return &ValidateHCLRulesResponse{
		Valid: true,
		// ParsedRules:   hclRuleList.ToParsedRuleList(),
		// UnparsedRules: hclRuleList.Other,
	}
}

func aclLinkCompare(a, b ACLLink) int {
	if a.Name < b.Name {
		return -1
	}
	if a.Name > b.Name {
		return 1
	}
	return 0
}

func (s *aclService) ListPolicies(ctx context.Context, options ListPoliciesOptions) ([]ACLLink, error) {
	resp, err := s.acl.ListPolicies(ctx)
	if err != nil {
		return nil, errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return nil, errPermissionDenied
	}
	if resp.Err != nil {
		slog.Error("failed to parse token list response", "error", resp.Err)
		return nil, errFailedToParse
	}
	policyList := make([]ACLLink, 0, len(resp.Body))
	for _, p := range resp.Body {
		isExclusive := ConseeExclusivePolicyNameRegexp.MatchString(p.Name)
		if options.Exclusive == "0" && isExclusive {
			continue
		} else if options.Exclusive == "1" && !isExclusive {
			continue
		}
		policyList = append(policyList, ACLLink{
			ID:   p.ID,
			Name: p.Name,
		})
	}
	slices.SortStableFunc(policyList, aclLinkCompare)
	return policyList, nil
}

func (s *aclService) ReadPolicy(ctx context.Context, name string) (*ReadPolicyResponse, error) {
	resp, err := s.acl.ReadPolicyByName(ctx, name)
	if err != nil {
		return nil, errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return nil, errPermissionDenied
	}
	if resp.Status == http.StatusNotFound {
		return nil, &DomainError{Code: DomainErrorCodeNotFound, Message: "policy not found"}
	}
	if resp.Err != nil || resp.Body == nil {
		return nil, errFailedToParse
	}
	parsedRules := []ParsedRule{}
	ruleList := HCLRuleList{}
	err = ParseHCLRules(resp.Body.Rules, &ruleList)
	if err != nil {
		slog.Error("Error parsing HCL rules", "err", err)
	}
	if err == nil {
		parsedRules = ruleList.ToParsedRuleList()
	}
	tokens, err := s.listPolicyTokens(ctx, resp.Body.ID)
	if err != nil {
		return nil, err
	}

	return &ReadPolicyResponse{
		ID:          resp.Body.ID,
		Name:        resp.Body.Name,
		Description: resp.Body.Description,
		ParsedRules: parsedRules,
		Rules:       resp.Body.Rules,
		Tokens:      tokens,
	}, nil
}

func (s *aclService) CreatePolicy(ctx context.Context, req *CreatePolicyRequest) (err error) {
	if ConseeExclusivePolicyNameRegexp.MatchString(req.Name) {
		return &DomainError{Code: DomainErrorCodeInvalidInput, Message: "exclusive policy can not be created separately"}
	}
	if req.Name == "" {
		return &DomainError{Code: DomainErrorCodeInvalidInput, Message: "policy name is required"}
	}
	resp1, err := s.acl.ReadPolicyByName(ctx, req.Name)
	if err != nil {
		return errFailedToConnectConsul
	}
	if resp1.Status == http.StatusForbidden {
		return errPermissionDenied
	}
	if resp1.Body != nil {
		return &DomainError{Code: DomainErrorCodeAlreadyExists, Message: "policy name already exists"}
	}

	resp, err := s.acl.CreatePolicy(ctx, &consul.ACLPolicy{
		Name:        req.Name,
		Description: req.Description,
		Rules:       req.Rules,
	})
	if err != nil {
		return errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return errPermissionDenied
	}
	// if resp.Status == http.StatusInternalServerError {
	// 	errMsg := string(resp.RawBody)
	// 	if strings.Contains(errMsg, "already exists") {
	// 		return &DomainError{Code: DomainErrorCodeAlreadyExists, Message: "policy name already exists"}
	// 	}
	// 	return &DomainError{Code: DomainErrorCodeInternalError, Message: errMsg}
	// }
	if resp.Status != http.StatusOK || resp.Body == nil {
		return errUnknown
	}
	return nil
}

func (s *aclService) DeletePolicy(ctx context.Context, name string) error {
	resp, err := s.acl.ReadPolicyByName(ctx, name)
	if err != nil {
		return errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return errPermissionDenied
	}
	if resp.Status == http.StatusNotFound {
		return &DomainError{Code: DomainErrorCodeNotFound, Message: "policy not found"}
	}
	if resp.Err != nil || resp.Body == nil {
		return errFailedToParse
	}
	if ConseeExclusivePolicyNameRegexp.MatchString(resp.Body.Name) {
		return &DomainError{Code: DomainErrorCodePermissionDenied, Message: "exclusive policy can not be deleted separately"}
	}

	return s.deletePolicy(ctx, resp.Body.ID)
}

func (s *aclService) deletePolicy(ctx context.Context, id string) error {
	resp, err := s.acl.DeletePolicy(ctx, id)
	if err != nil {
		return errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return errPermissionDenied
	}
	if resp.Status == http.StatusNotFound {
		return &DomainError{Code: DomainErrorCodeInternalError, Message: "policy not found"}
	}
	return nil
}

func (s *aclService) UpdatePolicyRule(ctx context.Context, name, rules string) error {
	// 1. 调用ReadPolicyByName检查policy是否存在
	resp, err := s.acl.ReadPolicyByName(ctx, name)
	if err != nil {
		return errFailedToConnectConsul
	}
	if resp.Status == http.StatusForbidden {
		return errPermissionDenied
	}
	if resp.Status == http.StatusNotFound {
		return &DomainError{Code: DomainErrorCodeNotFound, Message: "policy not found"}
	}
	if resp.Err != nil || resp.Body == nil {
		return errFailedToParse
	}

	// 2. ParseHCLRules检查rule合法性
	var ruleList HCLRuleList
	err = ParseHCLRules(rules, &ruleList)
	if err != nil {
		return &DomainError{Code: DomainErrorCodeInvalidInput, Message: "invalid HCL rules: " + err.Error()}
	}

	// 3. 修改repo相应方法实现更新规则
	policy := resp.Body
	policy.Rules = rules

	updateResp, err := s.acl.UpdatePolicy(ctx, policy)
	if err != nil {
		return errFailedToConnectConsul
	}
	if updateResp.Status == http.StatusForbidden {
		return errPermissionDenied
	}
	if updateResp.Status == http.StatusNotFound {
		return &DomainError{Code: DomainErrorCodeNotFound, Message: "policy not found during update"}
	}
	if updateResp.Err != nil {
		return errFailedToParse
	}

	return nil
}
