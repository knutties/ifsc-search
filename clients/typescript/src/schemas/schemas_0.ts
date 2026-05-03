const _BL = "BankList";
const _BNF = "BranchNotFound";
const _BR = "BadRequest";
const _BS = "BankSummary";
const _GB = "GetBranch";
const _GBI = "GetBranchInput";
const _GBO = "GetBranchOutput";
const _H = "Healthz";
const _HO = "HealthzOutput";
const _LB = "ListBanks";
const _LBO = "ListBanksOutput";
const _RI = "ResultItem";
const _RIL = "ResultItemList";
const _S = "Search";
const _SI = "SearchInput";
const _SO = "SearchOutput";
const _SOt = "StatusOutput";
const _St = "Status";
const _a = "address";
const _b = "branch";
const _ba = "banks";
const _ba_ = "built_at";
const _ban = "bank";
const _bc = "bank_code";
const _bn = "bank_name";
const _c = "client";
const _ce = "centre";
const _ci = "city";
const _co = "contact";
const _d = "district";
const _e = "error";
const _h = "http";
const _hE = "httpError";
const _hQ = "httpQuery";
const _i = "ifsc";
const _id = "indexed_docs";
const _im = "imps";
const _l = "limit";
const _m = "micr";
const _n = "neft";
const _o = "offset";
const _q = "q";
const _r = "rtgs";
const _re = "results";
const _rt = "release_tag";
const _rud = "rbi_update_date";
const _s = "smithy.ts.sdk.synthetic.io.knutties.banksearch";
const _sc = "score";
const _st = "state";
const _sta = "status";
const _sw = "swift";
const _t = "total";
const _u = "upi";
const n0 = "io.knutties.banksearch";

// smithy-typescript generated code
import { TypeRegistry } from "@smithy/core/schema";
import type { StaticErrorSchema, StaticListSchema, StaticOperationSchema, StaticStructureSchema } from "@smithy/types";

import { BankSearchServiceException } from "../models/BankSearchServiceException";
import { BadRequest, BranchNotFound } from "../models/errors";

/* eslint no-var: 0 */
const _s_registry = TypeRegistry.for(_s);
export var BankSearchServiceException$: StaticErrorSchema = [-3, _s, "BankSearchServiceException", 0, [], []];
_s_registry.registerError(BankSearchServiceException$, BankSearchServiceException);
const n0_registry = TypeRegistry.for(n0);
export var BadRequest$: StaticErrorSchema = [-3, n0, _BR,
  { [_e]: _c, [_hE]: 400 },
  [_e],
  [0], 1
];
n0_registry.registerError(BadRequest$, BadRequest);
export var BranchNotFound$: StaticErrorSchema = [-3, n0, _BNF,
  { [_e]: _c, [_hE]: 404 },
  [_e],
  [0], 1
];
n0_registry.registerError(BranchNotFound$, BranchNotFound);
/**
 * TypeRegistry instances containing modeled errors.
 * @internal
 *
 */
export const errorTypeRegistries = [
  _s_registry,
  n0_registry,
]
export var BankSummary$: StaticStructureSchema = [3, n0, _BS,
  0,
  [_bc, _bn],
  [0, 0], 2
];
export var GetBranchInput$: StaticStructureSchema = [3, n0, _GBI,
  0,
  [_i],
  [[0, 1]], 1
];
export var GetBranchOutput$: StaticStructureSchema = [3, n0, _GBO,
  0,
  [_i, _bc, _bn, _b, _ce, _d, _st, _a, _ci, _co, _m, _sw, _u, _n, _r, _im],
  [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2, 2], 4
];
export var HealthzOutput$: StaticStructureSchema = [3, n0, _HO,
  0,
  [_sta],
  [0], 1
];
export var ListBanksOutput$: StaticStructureSchema = [3, n0, _LBO,
  0,
  [_t, _ba],
  [1, () => BankList], 2
];
export var ResultItem$: StaticStructureSchema = [3, n0, _RI,
  0,
  [_i, _bc, _bn, _b, _sc, _ce, _d, _st, _a, _ci, _co, _m, _sw, _u, _n, _r, _im],
  [0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 2, 2, 2, 2], 5
];
export var SearchInput$: StaticStructureSchema = [3, n0, _SI,
  0,
  [_ban, _q, _i, _st, _d, _ci, _l, _o],
  [[0, { [_hQ]: _ban }], [0, { [_hQ]: _q }], [0, { [_hQ]: _i }], [0, { [_hQ]: _st }], [0, { [_hQ]: _d }], [0, { [_hQ]: _ci }], [1, { [_hQ]: _l }], [1, { [_hQ]: _o }]]
];
export var SearchOutput$: StaticStructureSchema = [3, n0, _SO,
  0,
  [_t, _l, _o, _re],
  [1, 1, 1, () => ResultItemList], 4
];
export var StatusOutput$: StaticStructureSchema = [3, n0, _SOt,
  0,
  [_sta, _id, _rt, _rud, _ba_],
  [0, 1, 0, 0, 0], 1
];
var __Unit = "unit" as const;
var BankList: StaticListSchema = [1, n0, _BL,
  0, () => BankSummary$
];
var ResultItemList: StaticListSchema = [1, n0, _RIL,
  0, () => ResultItem$
];
export var GetBranch$: StaticOperationSchema = [9, n0, _GB,
  { [_h]: ["GET", "/ifsc/{ifsc}", 200] }, () => GetBranchInput$, () => GetBranchOutput$
];
export var Healthz$: StaticOperationSchema = [9, n0, _H,
  { [_h]: ["GET", "/healthz", 200] }, () => __Unit, () => HealthzOutput$
];
export var ListBanks$: StaticOperationSchema = [9, n0, _LB,
  { [_h]: ["GET", "/list", 200] }, () => __Unit, () => ListBanksOutput$
];
export var Search$: StaticOperationSchema = [9, n0, _S,
  { [_h]: ["GET", "/search", 200] }, () => SearchInput$, () => SearchOutput$
];
export var Status$: StaticOperationSchema = [9, n0, _St,
  { [_h]: ["GET", "/status", 200] }, () => __Unit, () => StatusOutput$
];
