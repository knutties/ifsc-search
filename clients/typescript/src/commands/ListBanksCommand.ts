// smithy-typescript generated code
import { getEndpointPlugin } from "@smithy/middleware-endpoint";
import { Command as $Command } from "@smithy/smithy-client";
import type { MetadataBearer as __MetadataBearer } from "@smithy/types";

import type { BankSearchClientResolvedConfig, ServiceInputTypes, ServiceOutputTypes } from "../BankSearchClient";
import { commonParams } from "../endpoint/EndpointParameters";
import type { ListBanksOutput } from "../models/models_0";
import { ListBanks$ } from "../schemas/schemas_0";

/**
 * @public
 */
export type { __MetadataBearer };
export { $Command };
/**
 * @public
 *
 * The input for {@link ListBanksCommand}.
 */
export interface ListBanksCommandInput {}
/**
 * @public
 *
 * The output of {@link ListBanksCommand}.
 */
export interface ListBanksCommandOutput extends ListBanksOutput, __MetadataBearer {}

/**
 * List the distinct banks present in the index, sorted by bank_code.
 * @example
 * Use a bare-bones client and the command you need to make an API call.
 * ```javascript
 * import { BankSearchClient, ListBanksCommand } from "@knutties/bank-search-client"; // ES Modules import
 * // const { BankSearchClient, ListBanksCommand } = require("@knutties/bank-search-client"); // CommonJS import
 * // import type { BankSearchClientConfig } from "@knutties/bank-search-client";
 * const config = {}; // type is BankSearchClientConfig
 * const client = new BankSearchClient(config);
 * const input = {};
 * const command = new ListBanksCommand(input);
 * const response = await client.send(command);
 * // { // ListBanksOutput
 * //   total: Number("int"), // required
 * //   banks: [ // BankList // required
 * //     { // BankSummary
 * //       bank_code: "STRING_VALUE", // required
 * //       bank_name: "STRING_VALUE", // required
 * //     },
 * //   ],
 * // };
 *
 * ```
 *
 * @param ListBanksCommandInput - {@link ListBanksCommandInput}
 * @returns {@link ListBanksCommandOutput}
 * @see {@link ListBanksCommandInput} for command's `input` shape.
 * @see {@link ListBanksCommandOutput} for command's `response` shape.
 * @see {@link BankSearchClientResolvedConfig | config} for BankSearchClient's `config` shape.
 *
 * @throws {@link BadRequest} (client fault)
 *
 * @throws {@link BankSearchServiceException}
 * <p>Base exception class for all service exceptions from BankSearch service.</p>
 *
 *
 * @public
 */
export class ListBanksCommand extends $Command
  .classBuilder<
    ListBanksCommandInput,
    ListBanksCommandOutput,
    BankSearchClientResolvedConfig,
    ServiceInputTypes,
    ServiceOutputTypes
  >()
  .ep(commonParams)
  .m(function (this: any, Command: any, cs: any, config: BankSearchClientResolvedConfig, o: any) {
    return [getEndpointPlugin(config, Command.getEndpointParameterInstructions())];
  })
  .s("BankSearch", "ListBanks", {})
  .n("BankSearchClient", "ListBanksCommand")
  .sc(ListBanks$)
  .build() {
  /** @internal type navigation helper, not in runtime. */
  protected declare static __types: {
    api: {
      input: {};
      output: ListBanksOutput;
    };
    sdk: {
      input: ListBanksCommandInput;
      output: ListBanksCommandOutput;
    };
  };
}
