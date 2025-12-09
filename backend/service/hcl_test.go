// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package service

import (
	"testing"

	"github.com/hashicorp/hcl/v2"
)

type HCLRuleList2 struct {
	ACL []string `hcl:"acl"`
	B   hcl.Body `hcl:",remain"`
}

func TestHCLParse(t *testing.T) {
	hclString := `
key "foo" {
  policy = "read"
}

key_prefix "foo/" {
  policy = "write"
}

agent "foo" {
  policy = "write"
}
agent_prefix "" {
  policy = "read"
}
agent_prefix "bar" {
  policy = "deny"
}

service "foo" {
  policy = "read"
}

fds "other" "fds" {
  policy = "read"
}
`
	var v HCLRuleList2
	err := ParseHCLRules(hclString, &v)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	t.Log(v.ACL)
}
