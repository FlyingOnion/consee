// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package repo

// AdminRepo interface. It's a combination of KVRepo and ACLRepo.
// It's useful for testing and admin operations.
// In most cases, AdminService should be used instead.
type AdminRepo interface {
	KVRepo
	ACLRepo
}
