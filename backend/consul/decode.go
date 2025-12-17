// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package consul

import "encoding/json"

func decodeTrue(b []byte) (bool, error) { return string(b) == "true", nil }

func decodeStringSlice(b []byte) ([]string, error) {
	ss := []string{}
	err := json.Unmarshal(b, &ss)
	return ss, err
}

func decodeKVPair(b []byte) (*KVPair, error) {
	kvPairs := []*KVPair{}
	e := json.Unmarshal(b, &kvPairs)
	if e != nil {
		return nil, e
	}
	return kvPairs[0], nil
}

func decodeKVPairs(b []byte) ([]*KVPair, error) {
	kvPairs := []*KVPair{}
	e := json.Unmarshal(b, &kvPairs)
	return kvPairs, e
}

func decodeACLToken(b []byte) (*ACLToken, error) {
	token := &ACLToken{}
	e := json.Unmarshal(b, token)
	if e != nil {
		return nil, e
	}
	return token, nil
}

func decodeACLTokenList(b []byte) ([]*ACLToken, error) {
	tokens := []*ACLToken{}
	e := json.Unmarshal(b, &tokens)
	return tokens, e
}

func decodeACLPolicy(b []byte) (*ACLPolicy, error) {
	policy := &ACLPolicy{}
	e := json.Unmarshal(b, policy)
	if e != nil {
		return nil, e
	}
	return policy, nil
}

func decodeACLPolicyList(b []byte) ([]*ACLPolicy, error) {
	policies := []*ACLPolicy{}
	e := json.Unmarshal(b, &policies)
	return policies, e
}

func decodeACLRole(b []byte) (*ACLRole, error) {
	role := &ACLRole{}
	e := json.Unmarshal(b, role)
	if e != nil {
		return nil, e
	}
	return role, nil
}

func decodeACLRoleList(b []byte) ([]*ACLRole, error) {
	roles := []*ACLRole{}
	e := json.Unmarshal(b, &roles)
	return roles, e
}
