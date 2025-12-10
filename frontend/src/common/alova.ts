import type { Method, RequestBody } from "alova";
import {
  alova,
  b64Encode,
  type ACLLink,
  type CreateTokenRequest,
  type KeyValue,
  type PolicyDetailInfo,
  type PolicyFormRule,
  type TokenDetailInfo,
} from "./kz";
import { conseeErrorKey, conseeTokenKey } from "./const";

// export function consulGetTokens(): Method {
//   return alovaGet<ACLLink[]>("/acl/tokens", "Failed to fetch tokens");
// }

// export function consulGetTokenDetail(tokenId: string): Method {
//   return alovaGet<TokenDetailInfo>(`/acl/token/${tokenId}`, "Failed to fetch token detail");
// }

// export function consulGetPolicies(): Method {
//   return alovaGet<ACLLink[]>("/acl/policies", "Failed to fetch policies");
// }

// export function consulGetPolicyDetail(policyId: string): Method {
//   return alovaGet<PolicyDetailInfo>(`/acl/policy/${policyId}`, "Failed to fetch policy detail");
// }

// export function consulGetKeys(): Method {
//   return alovaGet<string[]>("/kv/keys", "Failed to fetch keys");
// }

// export function consulGetValue(b64key: string): Method {
//   return alovaGet(`/kv/value/${b64key}`, "Failed to fetch value");
// }

// export function consulGetValueType(b64key: string): Method {
//   return alova.Get(`/kv/valuetype/${b64key}`, {
//     headers: {
//       [conseeTokenKey]: localStorage.getItem(conseeTokenKey) || "",
//     },
//     cacheFor: null,
//   });
// }

// export function consulGetValueHistory(b64key: string): Method {
//   return alovaGet(`/kv/history/${b64key}`, "Failed to get history versions");
// }

export function parseHCL(hcl: string): Method {
  return alova.Post("/acl/hcl-rule", hcl, {
    transform: (resp: Response) => {
      return resp.status === 200
        ? (resp.json() as Promise<PolicyFormRule[]>)
        : Promise.reject(new Error(resp.headers.get(conseeErrorKey) || "Failed to parse HCL"));
    },
  });
}

export function notificationStatus(): Method {
  return alovaGet("/notifications/status", "Failed to fetch notification status");
}

export function notificationList(): Method {
  return alovaGet("/notifications", "Failed to fetch notification list");
}

export function alovaGet<T>(
  url: string,
  defaultErrorMsg: string,
  params?: { [x: string]: any },
): Method {
  return alova.Get(url, {
    headers: {
      [conseeTokenKey]: localStorage.getItem(conseeTokenKey) || "",
    },
    params,
    transform: (resp: Response) => {
      return resp.status === 200
        ? (resp.json() as Promise<T>)
        : Promise.reject(new Error(resp.headers.get(conseeErrorKey) || defaultErrorMsg));
    },
  });
}

export async function respToJson<T>(resp: Response): Promise<T> {
  return resp.clone().json() as Promise<T>;
}

export async function respToText(resp: Response): Promise<string> {
  return resp.clone().text();
}

export async function respToBlob(resp: Response): Promise<Blob> {
  return resp.clone().blob();
}

export async function doNothing(_: Response): Promise<void> {
  return;
}

export function alovaMethod<T>(
  url: string,
  options?: {
    name?: string;
    method?: string;
    withToken?: boolean;
    query?: { [key: string]: any };
    expectedStatus?: number;
    defaultErrorMsg?: string;
    transform?: (resp: Response) => Promise<T>;
  },
): Method {
  return alova.Request({
    url: `/api/v0${url}`,
    method: options?.method || "GET",
    headers: options?.withToken
      ? {
          [conseeTokenKey]: localStorage.getItem(conseeTokenKey) || "",
        }
      : undefined,
    transform: (resp: Response) => {
      if (resp.status === (options?.expectedStatus || 200)) {
        return options?.transform?.(resp);
      }
      throw new Error(
        resp.headers.get(conseeErrorKey) ||
          options?.defaultErrorMsg ||
          "Oops! Something went wrong.",
      );
    },
  });
}

type HitSource = string | RegExp | Method | (string | RegExp | Method)[];

export async function alovaCall<T = Response>(
  url: string,
  options?: {
    name?: string;
    method?: string;
    query?: { [key: string]: any };
    body?: RequestBody;
    withToken?: boolean;
    expectedStatus?: number;
    defaultErrorMsg?: string;
    transform?: (resp: Response) => Promise<T>;
    hitSource?: HitSource
  },
): Promise<T>;

export async function alovaCall(
  url: string,
  options?: {
    name?: string;
    method?: string;
    query?: { [key: string]: any };
    body?: RequestBody;
    withToken?: boolean;
    expectedStatus?: number;
    defaultErrorMsg?: string;
    hitSource?: HitSource;
  },
): Promise<Response>;

export async function alovaCall<T>(
  url: string,
  options?: {
    name?: string;
    method?: string;
    query?: { [key: string]: any };
    body?: RequestBody;
    withToken?: boolean;
    expectedStatus?: number;
    defaultErrorMsg?: string;
    transform?: (resp: Response) => Promise<T>;
    hitSource?: HitSource;
  },
): Promise<Response | T> {
  const resp = await alova.Request<Response>({
    url: `/api/v0${url}`,
    name: options?.name,
    data: options?.body,
    params: options?.query,
    method: options?.method || "GET",
    headers: options?.withToken
      ? {
          [conseeTokenKey]: localStorage.getItem(conseeTokenKey) || "",
        }
      : undefined,
    hitSource: options?.hitSource,
  });

  if (resp.status === (options?.expectedStatus || 200)) {
    return options?.transform ? await options.transform(resp) : resp;
  }
  throw new Error(
    resp.headers.get(conseeErrorKey) || options?.defaultErrorMsg || "Oops! Something went wrong.",
  );
}

export interface AuthResult {
  valid: 0 | 1;
  admin: 0 | 1;
  n?: number;
}

export function authenticate(): Promise<AuthResult> {
  return alovaCall("/authenticate", {
    method: "POST",
    name: "authenticate",
    withToken: true,
    expectedStatus: 200,
    defaultErrorMsg: "Failed to authenticate",
    transform: respToJson<AuthResult>,
  });
}

/* 一堆封装的方法 */

export function kvList(): Promise<string[]> {
  return alovaCall(`/kv/keys`, {
    name: "kvList",
    withToken: true,
    defaultErrorMsg: "Failed to get key list",
    transform: respToJson<string[]>,
    hitSource: ["kvCreate", "kvDelete"],
  });
}

export function kvCreate(data: { key: string; value: string; value_type: string }): Promise<void> {
  return alovaCall(`/kv/value`, {
    name: "kvCreate",
    method: "POST",
    withToken: true,
    expectedStatus: 201,
    body: data,
    defaultErrorMsg: "Failed to create key/value",
  });
}

export function kvUpdate(b64key: string, value: string): Promise<void> {
  return alovaCall(`/kv/value/${b64key}`, {
    name: "kvUpdate",
    method: "PUT",
    withToken: true,
    expectedStatus: 204,
    body: { value },
    defaultErrorMsg: "Failed to save key/value",
  });
}

export function kvUpdateValueType(b64key: string, vt: string): Promise<void> {
  return alovaCall(`/kv/valuetype/${b64key}`, {
    name: "kvUpdateValueType",
    method: "PUT",
    withToken: true,
    expectedStatus: 204,
    body: vt,
    defaultErrorMsg: "Failed to update value type",
  });
}

export function kvDelete(b64key: string): Promise<void> {
  return alovaCall(`/kv/value/${b64key}`, {
    name: "kvDelete",
    method: "DELETE",
    withToken: true,
    expectedStatus: 204,
    defaultErrorMsg: "Failed to delete key/value",
  });
}

export const kvDeleteHints = [
  "If keys/folders are created directly without creating their parent folders, deleting these keys/folders may also cause deletion of their parents.",
  "The following keys and/or folders will be deleted:",
];

export function kvGetValue(b64key: string, version?: string): Promise<KeyValue> {
  return alovaCall(`/kv/value/${b64key}`, {
    name: "kvGetValue",
    query: version ? { v: version } : undefined,
    withToken: true,
    defaultErrorMsg: "Failed to fetch value",
    transform: respToJson<KeyValue>,
    hitSource: ["kvCreate", "kvUpdate", "kvDelete"],
  });
}

export function kvGetValueType(b64key: string): Promise<string> {
  return alovaCall(`/kv/valuetype/${b64key}`, {
    name: "kvGetValueType",
    withToken: true,
    transform: respToText,
    defaultErrorMsg: "Failed to get value type",
    hitSource: ["kvUpdateValueType"],
  });
}

export function kvGetHistory(b64key: string): Promise<string[]> {
  return alovaCall<string[]>(`/kv/history/${b64key}`, {
    name: "kvGetHistory",
    withToken: true,
    transform: respToJson,
    defaultErrorMsg: "Failed to get value history",
    hitSource: ["kvUpdate", "kvDelete"],
  });
}

/* 一堆封装的方法 */

export function aclTokenCreate(req: CreateTokenRequest): Promise<void> {
  return alovaCall(`/acl/token`, {
    name: "aclTokenCreate",
    method: "POST",
    withToken: true,
    body: req,
    expectedStatus: 201,
    defaultErrorMsg: "Failed to create token",
  });
}

export function aclTokenList(): Promise<ACLLink[]> {
  return alovaCall("/acl/tokens", {
    name: "aclTokenList",
    withToken: true,
    transform: respToJson<ACLLink[]>,
    defaultErrorMsg: "Failed to get token list",
  });
}

export function aclTokenGetDetail(tokenId: string): Promise<TokenDetailInfo> {
  return alovaCall(`/acl/token/${tokenId}`, {
    name: "aclTokenGetDetail",
    withToken: true,
    transform: respToJson<TokenDetailInfo>,
    defaultErrorMsg: "Failed to get token detail",
  });
}

export function aclTokenUpdate(tokenId: string, policies?: string[]): Promise<void> {
  return alovaCall(`/acl/token/${tokenId}`, {
    name: "aclTokenUpdate",
    method: "PUT",
    withToken: true,
    expectedStatus: 204,
    body: { policies }, // TODO: add roles
    defaultErrorMsg: "Failed to update token",
  });
}

export function aclTokenDelete(tokenId: string): Promise<void> {
  return alovaCall(`/acl/token/${tokenId}`, {
    name: "aclTokenDelete",
    method: "DELETE",
    withToken: true,
    expectedStatus: 204,
    defaultErrorMsg: "Failed to delete token",
  });
}

export function aclPolicyList(query?: { exclusive: string }): Promise<ACLLink[]> {
  return alovaCall("/acl/policies", {
    name: "aclPolicyList",
    query,
    withToken: true,
    transform: respToJson<ACLLink[]>,
    defaultErrorMsg: "Failed to get policy list",
  });
}

export function aclPolicyGetDetail(b64policyName: string): Promise<PolicyDetailInfo> {
  return alovaCall(`/acl/policy/${b64policyName}`, {
    name: "aclPolicyGetDetail",
    withToken: true,
    transform: respToJson<PolicyDetailInfo>,
    defaultErrorMsg: "Failed to get policy detail",
  });
}

export function aclPolicyDelete(b64policyName: string): Promise<void> {
  return alovaCall(`/acl/policy/${b64policyName}`, {
    name: "aclPolicyDelete",
    method: "DELETE",
    withToken: true,
    expectedStatus: 204,
    defaultErrorMsg: "Failed to delete policy",
  });
}

/* 一堆封装的方法：导入导出 */

export interface ExportReq {
  keys?: string[];
  include_acl?: boolean;
  format?: string;
}

function exportReqToQuery(req?: ExportReq): { [key: string]: any } | undefined {
  return req
    ? {
        keys: (req.keys || []).map((k) => b64Encode(k)),
        include_acl: req.include_acl ? "1" : "0",
        format: req.format,
      }
    : undefined;
}

export function wExport(req?: ExportReq): Promise<Blob> {
  return alovaCall(`/export`, {
    name: "wExport",
    query: exportReqToQuery(req),
    withToken: true,
    transform: respToBlob,
    defaultErrorMsg: "Failed to export resources",
  });
}
