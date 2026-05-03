// smithy-typescript generated code
import type { ExceptionOptionType as __ExceptionOptionType } from "@smithy/smithy-client";

import { BankSearchServiceException as __BaseException } from "./BankSearchServiceException";

/**
 * @public
 */
export class BadRequest extends __BaseException {
  readonly name = "BadRequest" as const;
  readonly $fault = "client" as const;
  error: string | undefined;
  /**
   * @internal
   */
  constructor(opts: __ExceptionOptionType<BadRequest, __BaseException>) {
    super({
      name: "BadRequest",
      $fault: "client",
      ...opts,
    });
    Object.setPrototypeOf(this, BadRequest.prototype);
    this.error = opts.error;
  }
}

/**
 * @public
 */
export class BranchNotFound extends __BaseException {
  readonly name = "BranchNotFound" as const;
  readonly $fault = "client" as const;
  error: string | undefined;
  /**
   * @internal
   */
  constructor(opts: __ExceptionOptionType<BranchNotFound, __BaseException>) {
    super({
      name: "BranchNotFound",
      $fault: "client",
      ...opts,
    });
    Object.setPrototypeOf(this, BranchNotFound.prototype);
    this.error = opts.error;
  }
}
