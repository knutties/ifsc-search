// smithy-typescript generated code
import { createAggregatedClient } from "@smithy/smithy-client";
import type { HttpHandlerOptions as __HttpHandlerOptions } from "@smithy/types";

import { BankSearchClient } from "./BankSearchClient";
import { type GetBranchCommandInput, type GetBranchCommandOutput, GetBranchCommand } from "./commands/GetBranchCommand";
import { type HealthzCommandInput, type HealthzCommandOutput, HealthzCommand } from "./commands/HealthzCommand";
import { type ListBanksCommandInput, type ListBanksCommandOutput, ListBanksCommand } from "./commands/ListBanksCommand";
import { type SearchCommandInput, type SearchCommandOutput, SearchCommand } from "./commands/SearchCommand";
import { type StatusCommandInput, type StatusCommandOutput, StatusCommand } from "./commands/StatusCommand";

const commands = {
  GetBranchCommand,
  HealthzCommand,
  ListBanksCommand,
  SearchCommand,
  StatusCommand,
};

export interface BankSearch {
  /**
   * @see {@link GetBranchCommand}
   */
  getBranch(
    args: GetBranchCommandInput,
    options?: __HttpHandlerOptions
  ): Promise<GetBranchCommandOutput>;
  getBranch(
    args: GetBranchCommandInput,
    cb: (err: any, data?: GetBranchCommandOutput) => void
  ): void;
  getBranch(
    args: GetBranchCommandInput,
    options: __HttpHandlerOptions,
    cb: (err: any, data?: GetBranchCommandOutput) => void
  ): void;

  /**
   * @see {@link HealthzCommand}
   */
  healthz(): Promise<HealthzCommandOutput>;
  healthz(
    args: HealthzCommandInput,
    options?: __HttpHandlerOptions
  ): Promise<HealthzCommandOutput>;
  healthz(
    args: HealthzCommandInput,
    cb: (err: any, data?: HealthzCommandOutput) => void
  ): void;
  healthz(
    args: HealthzCommandInput,
    options: __HttpHandlerOptions,
    cb: (err: any, data?: HealthzCommandOutput) => void
  ): void;

  /**
   * @see {@link ListBanksCommand}
   */
  listBanks(): Promise<ListBanksCommandOutput>;
  listBanks(
    args: ListBanksCommandInput,
    options?: __HttpHandlerOptions
  ): Promise<ListBanksCommandOutput>;
  listBanks(
    args: ListBanksCommandInput,
    cb: (err: any, data?: ListBanksCommandOutput) => void
  ): void;
  listBanks(
    args: ListBanksCommandInput,
    options: __HttpHandlerOptions,
    cb: (err: any, data?: ListBanksCommandOutput) => void
  ): void;

  /**
   * @see {@link SearchCommand}
   */
  search(): Promise<SearchCommandOutput>;
  search(
    args: SearchCommandInput,
    options?: __HttpHandlerOptions
  ): Promise<SearchCommandOutput>;
  search(
    args: SearchCommandInput,
    cb: (err: any, data?: SearchCommandOutput) => void
  ): void;
  search(
    args: SearchCommandInput,
    options: __HttpHandlerOptions,
    cb: (err: any, data?: SearchCommandOutput) => void
  ): void;

  /**
   * @see {@link StatusCommand}
   */
  status(): Promise<StatusCommandOutput>;
  status(
    args: StatusCommandInput,
    options?: __HttpHandlerOptions
  ): Promise<StatusCommandOutput>;
  status(
    args: StatusCommandInput,
    cb: (err: any, data?: StatusCommandOutput) => void
  ): void;
  status(
    args: StatusCommandInput,
    options: __HttpHandlerOptions,
    cb: (err: any, data?: StatusCommandOutput) => void
  ): void;
}

/**
 * HTTP search service for Indian bank branches.
 * @public
 */
export class BankSearch extends BankSearchClient implements BankSearch {}
createAggregatedClient(commands, BankSearch);
