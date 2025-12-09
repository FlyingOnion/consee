// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package common

type CreateTokenRequest1 struct {
	AccessorID string
	SecretID   string
	Policies   []string
	Roles      []string
}

type CreateTokenRequest2 struct {
	AccessorID string
	Rules      string
	SecretID   string
}
