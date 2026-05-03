// smithy-typescript generated code
import { getEndpointPlugin } from "@smithy/middleware-endpoint";
import { Command as $Command } from "@smithy/smithy-client";
import type { MetadataBearer as __MetadataBearer } from "@smithy/types";

import type { BankSearchClientResolvedConfig, ServiceInputTypes, ServiceOutputTypes } from "../BankSearchClient";
import { commonParams } from "../endpoint/EndpointParameters";
import type { StatusOutput } from "../models/models_0";
import { Status$ } from "../schemas/schemas_0";

/**
 * @public
 */
export type { __MetadataBearer };
export { $Command };
/**
 * @public
 *
 * The input for {@link StatusCommand}.
 */
export interface StatusCommandInput {}
/**
 * @public
 *
 * The output of {@link StatusCommand}.
 */
export interface StatusCommandOutput extends StatusOutput, __MetadataBearer {}

/**
 * Index version metadata and document count.
 * @example
 * Use a bare-bones client and the command you need to make an API call.
 * ```javascript
 * import { BankSearchClient, StatusCommand } from "@knutties/bank-search-client"; // ES Modules import
 * // const { BankSearchClient, StatusCommand } = require("@knutties/bank-search-client"); // CommonJS import
 * // import type { BankSearchClientConfig } from "@knutties/bank-search-client";
 * const config = {}; // type is BankSearchClientConfig
 * const client = new BankSearchClient(config);
 * const input = {};
 * const command = new StatusCommand(input);
 * const response = await client.send(command);
 * // { // StatusOutput
 * //   status: "STRING_VALUE", // required
 * //   indexed_docs: Number("long"),
 * //   release_tag: "STRING_VALUE",
 * //   rbi_update_date: "STRING_VALUE",
 * //   built_at: "STRING_VALUE",
 * // };
 *
 * ```
 *
 * @param StatusCommandInput - {@link StatusCommandInput}
 * @returns {@link StatusCommandOutput}
 * @see {@link StatusCommandInput} for command's `input` shape.
 * @see {@link StatusCommandOutput} for command's `response` shape.
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
export class StatusCommand extends $Command
  .classBuilder<
    StatusCommandInput,
    StatusCommandOutput,
    BankSearchClientResolvedConfig,
    ServiceInputTypes,
    ServiceOutputTypes
  >()
  .ep(commonParams)
  .m(function (this: any, Command: any, cs: any, config: BankSearchClientResolvedConfig, o: any) {
    return [getEndpointPlugin(config, Command.getEndpointParameterInstructions())];
  })
  .s("BankSearch", "Status", {})
  .n("BankSearchClient", "StatusCommand")
  .sc(Status$)
  .build() {
  /** @internal type navigation helper, not in runtime. */
  protected declare static __types: {
    api: {
      input: {};
      output: StatusOutput;
    };
    sdk: {
      input: StatusCommandInput;
      output: StatusCommandOutput;
    };
  };
}
