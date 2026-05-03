// smithy-typescript generated code
/**
 * Wire shape for a bank in the listing.
 * @public
 */
export interface BankSummary {
  bank_code: string | undefined;
  bank_name: string | undefined;
}

/**
 * @public
 */
export interface ListBanksOutput {
  total: number | undefined;
  banks: BankSummary[] | undefined;
}

/**
 * @public
 */
export interface GetBranchInput {
  ifsc: string | undefined;
}

/**
 * @public
 */
export interface GetBranchOutput {
  ifsc: string | undefined;
  bank_code: string | undefined;
  bank_name: string | undefined;
  branch: string | undefined;
  centre?: string | undefined;
  district?: string | undefined;
  state?: string | undefined;
  address?: string | undefined;
  city?: string | undefined;
  contact?: string | undefined;
  micr?: string | undefined;
  swift?: string | undefined;
  upi?: boolean | undefined;
  neft?: boolean | undefined;
  rtgs?: boolean | undefined;
  imps?: boolean | undefined;
}

/**
 * @public
 */
export interface SearchInput {
  /**
   * 4-char IFSC bank code or fuzzy bank name.
   * @public
   */
  bank?: string | undefined;

  /**
   * Free-text query over branch, city, address and IFSC code prefix.
   * @public
   */
  q?: string | undefined;

  /**
   * Case-insensitive IFSC prefix, e.g. "HDFC0CAG".
   * @public
   */
  ifsc?: string | undefined;

  /**
   * Strict, case-insensitive state filter.
   * @public
   */
  state?: string | undefined;

  /**
   * Strict, case-insensitive district filter.
   * @public
   */
  district?: string | undefined;

  /**
   * Strict, case-insensitive city filter.
   * @public
   */
  city?: string | undefined;

  /**
   * Page size; defaults to 20 and is clamped to 100.
   * @public
   */
  limit?: number | undefined;

  /**
   * Result offset for pagination; defaults to 0.
   * @public
   */
  offset?: number | undefined;
}

/**
 * A search hit — a Branch's properties plus its relevance score.
 * @public
 */
export interface ResultItem {
  ifsc: string | undefined;
  bank_code: string | undefined;
  bank_name: string | undefined;
  branch: string | undefined;
  centre?: string | undefined;
  district?: string | undefined;
  state?: string | undefined;
  address?: string | undefined;
  city?: string | undefined;
  contact?: string | undefined;
  micr?: string | undefined;
  swift?: string | undefined;
  upi?: boolean | undefined;
  neft?: boolean | undefined;
  rtgs?: boolean | undefined;
  imps?: boolean | undefined;
  score: number | undefined;
}

/**
 * @public
 */
export interface SearchOutput {
  total: number | undefined;
  limit: number | undefined;
  offset: number | undefined;
  results: ResultItem[] | undefined;
}

/**
 * @public
 */
export interface HealthzOutput {
  status: string | undefined;
}

/**
 * @public
 */
export interface StatusOutput {
  status: string | undefined;
  indexed_docs?: number | undefined;
  release_tag?: string | undefined;
  rbi_update_date?: string | undefined;
  built_at?: string | undefined;
}
