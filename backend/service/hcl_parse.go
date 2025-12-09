// Copyright (c) 2025 The Consee Authors. All rights reserved.
// SPDX-License-Identifier: MulanPSL-2.0

package service

import (
	"slices"

	. "github.com/FlyingOnion/consee/backend/common"
	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type HCLRuleList struct {
	ACL            *string           `hcl:"acl"`
	Agent          []PolicyRuleBlock `hcl:"agent,block"`
	AgentPrefix    []PolicyRuleBlock `hcl:"agent_prefix,block"`
	Event          []PolicyRuleBlock `hcl:"event,block"`
	EventPrefix    []PolicyRuleBlock `hcl:"event_prefix,block"`
	Identity       []PolicyRuleBlock `hcl:"identity,block"`
	IdentityPrefix []PolicyRuleBlock `hcl:"identity_prefix,block"`
	Key            []PolicyRuleBlock `hcl:"key,block"`
	KeyPrefix      []PolicyRuleBlock `hcl:"key_prefix,block"`
	KeyRing        *string           `hcl:"keyring"`
	Mesh           *string           `hcl:"mesh"`
	Node           []PolicyRuleBlock `hcl:"node,block"`
	NodePrefix     []PolicyRuleBlock `hcl:"node_prefix,block"`
	Operator       *string           `hcl:"operator"`
	// Partition       []PolicyRuleBlock `hcl:"partition,block"`
	// PartitionPrefix []PolicyRuleBlock `hcl:"partition_prefix,block"`
	Peering       *string           `hcl:"peering"`
	Query         []PolicyRuleBlock `hcl:"query,block"`
	QueryPrefix   []PolicyRuleBlock `hcl:"query_prefix,block"`
	Service       []ServiceBlock    `hcl:"service,block"`
	ServicePrefix []ServiceBlock    `hcl:"service_prefix,block"`
	Session       []PolicyRuleBlock `hcl:"session,block"`
	SessionPrefix []PolicyRuleBlock `hcl:"session_prefix,block"`

	B     hcl.Body `hcl:",remain"`
	Other string
}

type PolicyRuleBlock struct {
	Label  string `hcl:",label"`
	Access string `hcl:"policy"`
}

type ServiceBlock struct {
	Label      string  `hcl:",label"`
	Access     string  `hcl:"policy"`
	Intentions *string `hcl:"intentions"`
}

func ParseHCLRules(rules string, v any) error {
	if rules == "" {
		return nil
	}
	f, d := hclsyntax.ParseConfig([]byte(rules), "", hcl.Pos{Line: 1, Column: 1})
	if d.HasErrors() {
		return d
	}
	d = gohcl.DecodeBody(f.Body, nil, v)
	if d.HasErrors() {
		return d
	}

	// wf, _ := hclwrite.ParseConfig([]byte(rules), "", hcl.Pos{Line: 1, Column: 1})
	// blocks := wf.Body().Blocks()
	// kBlocks := make([]*hclwrite.Block, 0, len(blocks))
	// for _, block := range blocks {
	// 	switch block.Type() {
	// 	case "acl", "agent", "agent_prefix", "key", "key_prefix", "service", "service_prefix":
	// 		kBlocks = append(kBlocks, block)
	// 	}
	// }
	// for _, block := range kBlocks {
	// 	wf.Body().RemoveBlock(block)
	// }
	// rule.Other = string(bytes.TrimSpace(wf.Bytes()))
	return nil
}

func parsedRuleCompare(a, b ParsedRule) int {
	if a.Type < b.Type {
		return -1
	}
	if a.Type > b.Type {
		return 1
	}
	if a.Match == b.Match {
		return 0
	}
	if a.Match == "all" {
		return 1
	}
	if b.Match == "all" {
		return -1
	}
	if a.Match == "prefix" {
		return 1
	}
	if a.Param < b.Param {
		return -1
	}
	if a.Param > b.Param {
		return 1
	}
	if a.Access == b.Access {
		return 0
	}
	if a.Access == "deny" {
		return 1
	}
	if b.Access == "deny" {
		return -1
	}
	if a.Access == "write" {
		return 1
	}
	return -1
}

func (r HCLRuleList) ToParsedRuleList() []ParsedRule {
	rules := make([]ParsedRule, 0, 5+
		len(r.Agent)+len(r.AgentPrefix)+
		len(r.Event)+len(r.EventPrefix)+
		len(r.Identity)+len(r.IdentityPrefix)+
		len(r.Key)+len(r.KeyPrefix)+
		len(r.Node)+len(r.NodePrefix)+
		// len(r.Partition)+len(r.PartitionPrefix)+
		len(r.Query)+len(r.QueryPrefix)+
		len(r.Service)+len(r.ServicePrefix)+
		len(r.Session)+len(r.SessionPrefix))
	if r.ACL != nil {
		rules = append(rules, ParsedRule{Type: "acl", Access: *r.ACL})
	}
	for _, a := range r.Agent {
		rules = append(rules, ParsedRule{Type: "agent", Match: "exact", Param: a.Label, Access: a.Access})
	}
	for _, a := range r.AgentPrefix {
		rules = append(rules, ParsedRule{Type: "agent", Match: "prefix", Param: a.Label, Access: a.Access})
	}
	for _, e := range r.Event {
		rules = append(rules, ParsedRule{Type: "event", Match: "exact", Param: e.Label, Access: e.Access})
	}
	for _, e := range r.EventPrefix {
		rules = append(rules, ParsedRule{Type: "event", Match: "prefix", Param: e.Label, Access: e.Access})
	}
	for _, i := range r.Identity {
		rules = append(rules, ParsedRule{Type: "identity", Match: "exact", Param: i.Label, Access: i.Access})
	}
	for _, i := range r.IdentityPrefix {
		rules = append(rules, ParsedRule{Type: "identity", Match: "prefix", Param: i.Label, Access: i.Access})
	}
	for _, k := range r.Key {
		rules = append(rules, ParsedRule{Type: "key", Match: "exact", Param: k.Label, Access: k.Access})
	}
	for _, k := range r.KeyPrefix {
		rules = append(rules, ParsedRule{Type: "key", Match: "prefix", Param: k.Label, Access: k.Access})
	}
	if r.KeyRing != nil {
		rules = append(rules, ParsedRule{Type: "keyring", Access: *r.KeyRing})
	}
	if r.Mesh != nil {
		rules = append(rules, ParsedRule{Type: "mesh", Access: *r.Mesh})
	}
	for _, n := range r.Node {
		rules = append(rules, ParsedRule{Type: "node", Match: "exact", Param: n.Label, Access: n.Access})
	}
	for _, n := range r.NodePrefix {
		rules = append(rules, ParsedRule{Type: "node", Match: "prefix", Param: n.Label, Access: n.Access})
	}
	if r.Operator != nil {
		rules = append(rules, ParsedRule{Type: "operator", Access: *r.Operator})
	}
	// for _, p := range r.Partition {
	// 	rules = append(rules, ParsedRule{Type: "partition", Match: "exact", Param: p.Label, Access: p.Access})
	// }
	// for _, p := range r.PartitionPrefix {
	// 	rules = append(rules, ParsedRule{Type: "partition", Match: "prefix", Param: p.Label, Access: p.Access})
	// }
	if r.Peering != nil {
		rules = append(rules, ParsedRule{Type: "peering", Access: *r.Peering})
	}
	for _, q := range r.Query {
		rules = append(rules, ParsedRule{Type: "query", Match: "exact", Param: q.Label, Access: q.Access})
	}
	for _, q := range r.QueryPrefix {
		rules = append(rules, ParsedRule{Type: "query", Match: "prefix", Param: q.Label, Access: q.Access})
	}
	for _, s := range r.Service {
		rules = append(rules, ParsedRule{Type: "service", Match: "exact", Param: s.Label, Access: s.Access})
	}
	for _, s := range r.ServicePrefix {
		rules = append(rules, ParsedRule{Type: "service", Match: "prefix", Param: s.Label, Access: s.Access})
	}
	for _, s := range r.Session {
		rules = append(rules, ParsedRule{Type: "session", Match: "exact", Param: s.Label, Access: s.Access})
	}
	for _, s := range r.SessionPrefix {
		rules = append(rules, ParsedRule{Type: "session", Match: "prefix", Param: s.Label, Access: s.Access})
	}

	slices.SortStableFunc(rules, parsedRuleCompare)
	return rules
}
