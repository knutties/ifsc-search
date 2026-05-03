// smithy-typescript generated code
import { getEndpointPlugin } from "@smithy/middleware-endpoint";
import { Command as $Command } from "@smithy/smithy-client";
import type { MetadataBearer as __MetadataBearer } from "@smithy/types";

import type { BankSearchClientResolvedConfig, ServiceInputTypes, ServiceOutputTypes } from "../BankSearchClient";
import { commonParams } from "../endpoint/EndpointParameters";
import type { GetBranchInput, GetBranchOutput } from "../models/models_0";
import { GetBranch$ } from "../schemas/schemas_0";

/**
 * @public
 */
export type { __MetadataBearer };
export { $Command };
/**
 * @public
 *
 * The input for {@link GetBranchCommand}.
 */
export interface GetBranchCommandInput extends GetBranchInput {}
/**
 * @public
 *
 * The output of {@link GetBranchCommand}.
 */
export interface GetBranchCommandOutput extends GetBranchOutput, __MetadataBearer {}

/**
 * Look up a single branch by IFSC code.
 * @example
 * Use a bare-bones client and the command you need to make an API call.
 * ```javascript
 * import { BankSearchClient, GetBranchCommand } from "@knutties/bank-search-client"; // ES Modules import
 * // const { BankSearchClient, GetBranchCommand } = require("@knutties/bank-search-client"); // CommonJS import
 * // import type { BankSearchClientConfig } from "@knutties/bank-search-client";
 * const config = {}; // type is BankSearchClientConfig
 * const client = new BankSearchClient(config);
 * const input = { // GetBranchInput
 *   ifsc: "STRING_VALUE", // required
 * };
 * const command = new GetBranchCommand(input);
 * const response = await client.send(command);
 * // { // GetBranchOutput
 * //   ifsc: "STRING_VALUE", // required
 * //   bank_code: "STRING_VALUE", // required
 * //   bank_name: "STRING_VALUE", // required
 * //   branch: "STRING_VALUE", // required
 * //   centre: "STRING_VALUE",
 * //   district: "STRING_VALUE",
 * //   state: "STRING_VALUE",
 * //   address: "STRING_VALUE",
 * //   city: "STRING_VALUE",
 * //   contact: "STRING_VALUE",
 * //   micr: "STRING_VALUE",
 * //   swift: "STRING_VALUE",
 * //   upi: true || false,
 * //   neft: true || false,
 * //   rtgs: true || false,
 * //   imps: true || false,
 * // };
 *
 * ```
 *
 * @param GetBranchCommandInput - {@link GetBranchCommandInput}
 * @returns {@link GetBranchCommandOutput}
 * @see {@link GetBranchCommandInput} for command's `input` shape.
 * @see {@link GetBranchCommandOutput} for command's `response` shape.
 * @see {@link BankSearchClientResolvedConfig | config} for BankSearchClient's `config` shape.
 *
 * @throws {@link BranchNotFound} (client fault)
 *
 * @throws {@link BadRequest} (client fault)
 *
 * @throws {@link BankSearchServiceException}
 * <p>Base exception class for all service exceptions from BankSearch service.</p>
 *
 *
 * @public
 */
export class GetBranchCommand extends $Command
  .classBuilder<
    GetBranchCommandInput,
    GetBranchCommandOutput,
    BankSearchClientResolvedConfig,
    ServiceInputTypes,
    ServiceOutputTypes
  >()
  .ep(commonParams)
  .m(function (this: any, Command: any, cs: any, config: BankSearchClientResolvedConfig, o: any) {
    return [getEndpointPlugin(config, Command.getEndpointParameterInstructions())];
  })
  .s("BankSearch", "GetBranch", {})
  .n("BankSearchClient", "GetBranchCommand")
  .sc(GetBranch$)
  .build() {
  /** @internal type navigation helper, not in runtime. */
  protected declare static __types: {
    api: {
      input: GetBranchInput;
      output: GetBranchOutput;
    };
    sdk: {
      input: GetBranchCommandInput;
      output: GetBranchCommandOutput;
    };
  };
}
