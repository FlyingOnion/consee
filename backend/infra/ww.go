// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package infra

import (
	"github.com/FlyingOnion/consee/backend/consul"
	"github.com/FlyingOnion/consee/backend/repo"
)

var (
	_ repo.KVRepo  = &kv{}
	_ repo.KVRepo  = &admin{}
	_ repo.ACLRepo = &acl{}
	_ repo.ACLRepo = &admin{}
)

func NewKV(client *consul.Client) repo.KVRepo {
	return &kv{client: client}
}

func NewAdmin(client *consul.Client, qAdmin *consul.QueryOptions, wAdmin *consul.WriteOptions) repo.AdminRepo {
	return &admin{qAdmin, wAdmin, client}
}

func NewACL(client *consul.Client) repo.ACLRepo {
	return &acl{client: client}
}
