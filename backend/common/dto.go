// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package common

import (
	"github.com/FlyingOnion/consee/backend/buffer"
)

type KeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type GetValueResponse KeyValue

type CreateKeyValueRequest struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	ValueType string `json:"value_type"`
}

type UpdateValueRequest struct {
	Value string `json:"value"`
}

type BatchUpdateRequest struct {
	KeyValues []*KeyValue `json:"kvs"`
}

// import

type ImportResponseItem struct {
	// Kind specifies the error occurs on which kind when importing
	Kind  string `json:"kind"`
	Param string `json:"param"`
	Cause string `json:"cause,omitempty"`
}

type ImportResponse struct {
	Successes []ImportResponseItem `json:"successes"`
	Conflicts []ImportResponseItem `json:"conflicts"`
	Errors    []ImportResponseItem `json:"errors"`
}

type OnConflictPolicy string

const (
	OnConflictPolicySkip    = "skip"
	OnConflictPolicyReplace = "replace"
)

type ImportRequest struct {
	Format      string
	Dryrun      bool
	OnConflict  OnConflictPolicy
	FileContent []byte
}

// export

type DryrunMetadata struct {
	Keys     []string  `json:"keys"`
	Tokens   []ACLLink `json:"tokens"`
	Policies []string  `json:"policies"`
}

type CompatibleKVMetaList []*CompatibleKVMeta

func (l CompatibleKVMetaList) DryrunMetadata() *DryrunMetadata {
	keys := make([]string, len(l))
	for i, kv := range l {
		keys[i] = kv.Key
	}
	return &DryrunMetadata{
		Keys: keys,
	}
}

type CompatibleKVMeta struct {
	Key   string `json:"key"`
	Flags uint64 `json:"flags"`
	Value []byte `json:"value"`
}

// zip export
type ExportedKVMeta struct {
	Name            string   `json:"name"`
	ValueType       string   `json:"value_type"`
	HistoryVersions []string `json:"history_versions"`
}

type ExportMetadata struct {
	Keys     []ExportedKVMeta `json:"keys" yaml:"keys"`
	Tokens   []ACLLink        `json:"tokens" yaml:"tokens"`
	Policies []string         `json:"policies" yaml:"policies"`
}

func (m *ExportMetadata) DryrunMetadata() *DryrunMetadata {
	keys := make([]string, len(m.Keys))
	for i, kv := range m.Keys {
		keys[i] = kv.Name
	}
	return &DryrunMetadata{
		Keys:     keys,
		Tokens:   m.Tokens,
		Policies: m.Policies,
	}
}

type ExportRequest struct {
	Keys   []string `json:"keys"`
	Format string   `json:"format"`
	ACL    bool     `json:"acl"`
}

// type KVMeta struct {
// 	ValueType string `json:"value_type"`
// }

type ListNotificationsResponse struct {
	Open     []Notification         `json:"open"`
	Archived []ArchivedNotification `json:"archived"`
}

type NotificationOp string

const (
	NotificationOpOK           NotificationOp = "ok"
	NotificationOpAcceptReject NotificationOp = "accept_reject"
)

type NotificationType string

const (
	NotificationTypeTokenApplication NotificationType = "token_application"
	NotificationTypeOther            NotificationType = "other"
)

type Notification struct {
	ID            string           `json:"id"`
	Type          NotificationType `json:"type"`
	Data          []byte           `json:"data"`
	Operation     NotificationOp   `json:"operation"`
	OperationArgs map[string]any   `json:"operation_args,omitempty"`
	CreatedAt     string           `json:"created_at"`
	CreatedBy     string           `json:"created_by"`
}

func CompareNotifications(a, b Notification) int {
	if a.CreatedAt < b.CreatedAt {
		return 1
	}
	if a.CreatedAt > b.CreatedAt {
		return -1
	}
	if a.ID < b.ID {
		return -1
	}
	if a.ID > b.ID {
		return 1
	}
	if a.Type < b.Type {
		return -1
	}
	if a.Type > b.Type {
		return 1
	}
	return 0
}

type ArchivedNotification struct {
	ID         string           `json:"id"`
	Type       NotificationType `json:"type"`
	OriginData any              `json:"origin_data"`
	// could be empty if Operation is ok,
	//
	// "accepted" or "rejected" if Operation is accept_deny
	Reason     string `json:"reason"`
	CreatedAt  string `json:"created_at"`
	CreatedBy  string `json:"created_by"`
	ArchivedAt string `json:"archived_at"`
	ArchivedBy string `json:"archived_by"`
}

func CompareArchivedNotifications(a, b ArchivedNotification) int {
	if a.ArchivedAt < b.ArchivedAt {
		return 1
	}
	if a.ArchivedAt > b.ArchivedAt {
		return -1
	}
	if a.CreatedAt < b.CreatedAt {
		return 1
	}
	if a.CreatedAt > b.CreatedAt {
		return -1
	}
	if a.ID < b.ID {
		return -1
	}
	if a.ID > b.ID {
		return 1
	}
	if a.Type < b.Type {
		return -1
	}
	if a.Type > b.Type {
		return 1
	}
	return 0
}

// type NotificationStatus struct {
// 	Open     int `json:"open"`
// 	Archived int `json:"archived"`
// }

type TokenBasicInfo struct {
	Name       string `json:"name"`        // token name
	AccessorID string `json:"accessor_id"` // token accessor id
}

type ACLLink struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ReadTokenResponse struct {
	AccessorID string         `json:"accessor_id"`
	SecretID   string         `json:"secret_id"`
	Policies   []ACLLink      `json:"policies"`
	Roles      []ACLLink      `json:"roles"`
	Name       string         `json:"name"`
	Metadata   *TokenMetadata `json:"metadata"`
}

type CreateTokenRequest struct {
	AccessorID string `json:"accessor_id"`
	SecretID   string `json:"secret_id"`
	Name       string `json:"name"`
	// PolicyMode specifies whether the token is common or with an exclusive policy
	//  "", "common" // token applying common policies
	//  "exclusive"  // token with an exclusive policy
	PolicyMode string   `json:"policy_mode"`
	Rules      string   `json:"rules"`
	Policies   []string `json:"policies"`
	Roles      []string `json:"roles"`
}

type UpdateTokenRequest struct {
	Policies []string `json:"policies"`
	Roles    []string `json:"roles"`
}

type TokenMetadata struct {
	// metadata from consul kv
	CreatedAt     string `json:"created_at"`
	CreatedBy     string `json:"created_by"`
	LastUpdatedAt string `json:"last_updated_at"`
	LastUpdatedBy string `json:"last_updated_by"`
	// Version is a "yyyy-MM-dd hh:mm:ss" timestamp
	Version string `json:"version"`
	// From is the version that the token was copied from
	From string `json:"from"`
}

func (t TokenMetadata) MarshalJSON() ([]byte, error) {
	var b buffer.Buffer
	b.WriteString(`{"created_at":"`).WriteJsonSafeString(t.CreatedAt).
		WriteString(`","created_by":"`).WriteJsonSafeString(t.CreatedBy).
		WriteString(`","last_updated_at":"`).WriteJsonSafeString(t.LastUpdatedAt).
		WriteString(`","last_updated_by":"`).WriteJsonSafeString(t.LastUpdatedBy).
		WriteString(`","version":"`).WriteJsonSafeString(t.Version).
		WriteString(`","from":"`).WriteJsonSafeString(t.From).
		WriteString(`"}`)
	return b.Bytes(), nil
}

type ValidateHCLRulesResponse struct {
	Valid         bool         `json:"valid"`
	ParsedRules   []ParsedRule `json:"parsed"`
	UnparsedRules string       `json:"unparsed"`
	Error         string       `json:"error"`
}

type ListPoliciesOptions struct {
	// Exclusive filter policies by their exclusiveness
	//  "1": exclusive policies only
	//  "0": non-exclusive policies only
	//  "" : all policies
	Exclusive string
}

type CreatePolicyRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Rules       string `json:"rules"`
}

type ParsedRule struct {
	Type   string `json:"rtype"`
	Match  string `json:"match"`
	Param  string `json:"param"`
	Access string `json:"access"`
}

func (r ParsedRule) MarshalJSON() ([]byte, error) {
	var b buffer.Buffer
	b.WriteString(`{"rtype":"`).WriteJsonSafeString(r.Type).
		WriteString(`","access":"`).WriteJsonSafeString(r.Access)
	switch r.Type {
	case "acl", "keyring", "mesh", "operator", "peering":
	default:
		b.WriteString(`","match":"`).WriteJsonSafeString(r.Match).
			WriteString(`","param":"`).WriteJsonSafeString(r.Param)
	}
	b.WriteString(`"}`)
	return b.Bytes(), nil
}

type ReadPolicyResponse struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	ParsedRules []ParsedRule `json:"parsed_rules"`
	Rules       string       `json:"rules"`
	Tokens      []ACLLink    `json:"tokens"`
}

type ListRolesOptions struct {
}

type CreateRoleRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Policies    []string `json:"policies"`
}

type ReadRoleResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Policies    []ACLLink `json:"policies"`
}

type UpdateRoleRequest struct {
	Policies []string `json:"policies"`
}

type TokenApplicationRequest struct {
	AccessorID string `json:"accessor_id"`
	SecretID   string `json:"secret_id"`
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
	Rules      string `json:"rules"`
}

// TokenApplicationResponse does not contain review result.
// It only contains request info because applier could send request without token
type TokenApplicationResponse struct {
	AccessorID string `json:"accessor_id"`
	SecretID   string `json:"secret_id"`
	Name       string `json:"name"`
}

type HandleTokenApplicationRequest struct {
	Result string `json:"result"` // "accept" or "reject"
	Reason string `json:"reason"` // could be empty if accepted
}

type TokenApplicationReviewResult struct {
	Action       string `json:"action"`
	RejectReason string `json:"reject_reason"`

	AccessorID string    `json:"accessor_id"`
	Token      string    `json:"token"`
	Name       string    `json:"name"`
	Policies   []ACLLink `json:"policies"`
	Roles      []ACLLink `json:"roles"`

	ReviewedAt string `json:"reviewed_at"`
	Reviewer   string `json:"reviewer"`
}

type AuthenticateResult struct {
	IsValid                int `json:"valid"`
	IsAdmin                int `json:"admin"`
	OpenNotificationsCount int `json:"n,omitempty"`
}
