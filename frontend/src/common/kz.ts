export interface TreeData {
  key: string; // key name
  path: string; // full path key name
  isLeaf: boolean;
  depth: number;
}

export type Access = "read" | "write" | "deny" | "";

export interface TreeItem extends TreeData {
  children: TreeItem[];
  access: Access;
}

const rootItem: TreeItem = {
  key: "/",
  path: "/",
  isLeaf: false,
  depth: 0,
  children: [],
  access: "",
};

export function parseKey(key: string, sep: string): TreeData[] {
  if (key.length === 0) {
    return [];
  }
  const result: TreeData[] = [];
  for (let i = 0, lastSep = -1, depth = 1; ; ) {
    if (i === key.length - 1) {
      const isLeaf = key[i] !== sep;
      result.push({
        key: key.substring(lastSep + 1, isLeaf ? i + 1 : i),
        path: key.substring(0, i + 1),
        depth: depth,
        isLeaf: isLeaf,
      });
      break;
    }

    if (key[i] !== sep) {
      i++;
      continue;
    }

    result.push({
      key: key.substring(lastSep + 1, i),
      path: key.substring(0, i + 1),
      depth: depth,
      isLeaf: false,
    });
    lastSep = i;
    i++;
    depth++;
  }
  return result;
}

export function parseTree(
  keys: string[],
  options?: {
    ruleMap?: Map<string, Access>;
    withRoot?: boolean;
  }
): TreeItem[] {
  const tree: TreeItem[] = options?.withRoot
    ? [{ ...rootItem, access: options?.ruleMap?.get("") || "" }]
    : [];
  for (const key of keys) {
    const dataList = parseKey(key, "/");
    let current = tree;
    for (const treeData of dataList) {
      let target = current.find((item) => item.path === treeData.path);
      if (target === undefined) {
        target = {
          ...treeData,
          children: [],
          access: options?.ruleMap?.get(treeData.path) || "",
        };
        current.push(target);
      }
      current = target.children;
    }
  }
  return tree;
}

// use parseTree for most cases.
// we export this function only for tests.
export function addNewKey(tree: TreeItem[], key: string) {
  if (key.length === 0) {
    return;
  }
  const dataList = parseKey(key, "/");
  let current = tree;
  for (let i = 0; i < dataList.length; i++) {
    const treeData = dataList[i];
    let target = current.find((item) => item.path === treeData.path);
    if (target === undefined) {
      target = {
        ...treeData,
        children: [],
        access: "",
      };
      current.push(target);
    }
    current = target.children;
  }
}

import { createAlova } from "alova";
import VueHook from "alova/vue";
import adapterFetch from "alova/fetch";

export const alova = createAlova({
  requestAdapter: adapterFetch(),
  statesHook: VueHook,
  timeout: 3000,
});

const decoder = new TextDecoder();
export function b64Decode(b64: string): string {
  return decoder.decode(new Uint8Array([...atob(b64)].map((c) => c.charCodeAt(0))));
}

const encoder = new TextEncoder();
export function b64Encode(str: string): string {
  return btoa(String.fromCharCode(...encoder.encode(str)));
}

export interface CreateTokenRequest {
  accessor_id: string;
  secret_id: string;
  name?: string;
  policy_mode?: "common" | "exclusive";
  rules?: string;
  policies?: string[];
}

export interface ACLLink {
  id: string;
  name: string;
}

export interface TokenMetadata {
  created_at: string;
  created_by: string;
  last_updated_at: string;
  last_updated_by: string;
  version: string;
}

export interface TokenDetailInfo {
  accessor_id: string;
  secret_id: string;
  name: string;
  policies: ACLLink[];
  roles: ACLLink[];
  metadata: TokenMetadata;
}

export const valueTypeOptions = [
  "plaintext",
  "cmake",
  "hcl",
  "ini",
  "json",
  "json5",
  "jsonp",
  "lua",
  "makefile",
  "plsql",
  "properties",
  "qml",
  "sql",
  "toml",
  "xml",
  "yaml",
];

export interface PolicyDetailInfo {
  id: string;
  name: string;
  description: string;
  rules: string;
  parsed_rules: PolicyFormRule[];
  tokens: ACLLink[];
}

export interface RoleDetailInfo {
  id: string;
  name: string;
  description: string;
  policies: ACLLink[];
}

type PolicyFormRuleType =
  | "acl"
  | "agent"
  | "event"
  | "identity"
  | "key"
  | "keyring"
  | "mesh"
  | "node"
  | "operator"
  | "partition"
  | "peering"
  | "query"
  | "service"
  | "session"
  | "";

export type PolicyFormRuleAccess = "read" | "write" | "deny";

export interface PolicyRulePreset {
  rtype: PolicyFormRuleType;
  param: boolean;
}

export const policyRuleTypeList: PolicyFormRuleType[] = [
  "acl",
  "agent",
  "event",
  "identity",
  "key",
  "keyring",
  "mesh",
  "node",
  "operator",
  "partition",
  "peering",
  "query",
  "service",
  "session",
];

export const policyRuleWithParamTypeSet = new Set<PolicyFormRuleType>([
  "agent",
  "event",
  "identity",
  "key",
  "node",
  "partition",
  "query",
  "service",
  "session",
]);

export type PolicyFormRuleMatchType = "prefix" | "exact";
export interface PolicyFormRule {
  rtype: PolicyFormRuleType;
  match?: PolicyFormRuleMatchType | "all";
  param?: string;
  access: PolicyFormRuleAccess;
}

export interface PolicyFormRuleListElement {
  id: string;
  rule: PolicyFormRule;
}

export const uuidRegexp = /^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/;
export const exclusivePolicyRegexp =
  /^--[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/;

export function parseTimestampFromUUID(uuid: string): string {
  if (!uuidRegexp.test(uuid)) {
    return "unknown time";
  }

  // 提取时间戳部分（前 12 个字符，48 位）
  const timestampHex = `${uuid.substring(0, 8)}${uuid.substring(9, 13)}`;
  const timestampMs = parseInt(timestampHex, 16);

  // 处理无效时间戳
  if (isNaN(timestampMs)) {
    return "unknown time";
  }

  // 创建 Date 对象并格式化为本地时间
  const date = new Date(timestampMs);

  // 提取日期时间组件
  const year = date.getFullYear();
  const month = (date.getMonth() + 1).toString().padStart(2, "0");
  const day = date.getDate().toString().padStart(2, "0");
  const hours = date.getHours().toString().padStart(2, "0");
  const minutes = date.getMinutes().toString().padStart(2, "0");
  const seconds = date.getSeconds().toString().padStart(2, "0");
  const milliseconds = date.getMilliseconds().toString().padStart(3, "0");

  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}.${milliseconds}`;
}

export interface KeyValue {
  key: string;
  value: string;
}

// Debounce 函数
export function debounce<T extends (...args: any[]) => any>(
  func: T,
  delay: number
): (...args: Parameters<T>) => void {
  let timeoutId: number | null = null;

  return (...args: Parameters<T>) => {
    if (timeoutId) {
      clearTimeout(timeoutId);
    }

    timeoutId = setTimeout(() => {
      func(...args);
    }, delay);
  };
}
