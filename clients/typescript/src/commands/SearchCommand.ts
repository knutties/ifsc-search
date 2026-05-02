// smithy-typescript generated code
import { getEndpointPlugin } from "@smithy/middleware-endpoint";
import { Command as $Command } from "@smithy/smithy-client";
import type { MetadataBearer as __MetadataBearer } from "@smithy/types";

import type { BankSearchClientResolvedConfig, ServiceInputTypes, ServiceOutputTypes } from "../BankSearchClient";
import { commonParams } from "../endpoint/EndpointParameters";
import type { SearchInput, SearchOutput } from "../models/models_0";
import { Search$ } from "../schemas/schemas_0";

/**
 * @public
 */
export type { __MetadataBearer };
export { $Command };
/**
 * @public
 *
 * The input for {@link SearchCommand}.
 */
export interface SearchCommandInput extends SearchInput {}
/**
 * @public
 *
 * The output of {@link SearchCommand}.
 */
export interface SearchCommandOutput extends SearchOutput, __MetadataBearer {}

/**
 * Free-text plus structured search across the indexed branches.
 * At least one of \{bank, q, ifsc, state, district, city\} is required.
 * @example
 * Use a bare-bones client and the command you need to make an API call.
 * ```javascript
 * import { BankSearchClient, SearchCommand } from "@knutties/bank-search-client"; // ES Modules import
 * // const { BankSearchClient, SearchCommand } = require("@knutties/bank-search-client"); // CommonJS import
 * // import type { BankSearchClientConfig } from "@knutties/bank-search-client";
 * const config = {}; // type is BankSearchClientConfig
 * const client = new BankSearchClient(config);
 * const input = { // SearchInput
 *   bank: "STRING_VALUE",
 *   q: "STRING_VALUE",
 *   ifsc: "STRING_VALUE",
 *   state: "STRING_VALUE",
 *   district: "STRING_VALUE",
 *   city: "STRING_VALUE",
 *   limit: Number("int"),
 *   offset: Number("int"),
 * };
 * const command = new SearchCommand(input);
 * const response = await client.send(command);
 * // { // SearchOutput
 * //   total: Number("int"), // required
 * //   limit: Number("int"), // required
 * //   offset: Number("int"), // required
 * //   results: [ // ResultItemList // required
 * //     { // ResultItem
 * //       ifsc: "STRING_VALUE", // required
 * //       bank_code: "STRING_VALUE", // required
 * //       bank_name: "STRING_VALUE", // required
 * //       branch: "STRING_VALUE", // required
 * //       centre: "STRING_VALUE",
 * //       district: "STRING_VALUE",
 * //       state: "STRING_VALUE",
 * //       address: "STRING_VALUE",
 * //       city: "STRING_VALUE",
 * //       contact: "STRING_VALUE",
 * //       micr: "STRING_VALUE",
 * //       swift: "STRING_VALUE",
 * //       upi: true || false,
 * //       neft: true || false,
 * //       rtgs: true || false,
 * //       imps: true || false,
 * //       score: Number("double"), // required
 * //     },
 * //   ],
 * // };
 *
 * ```
 *
 * @param SearchCommandInput - {@link SearchCommandInput}
 * @returns {@link SearchCommandOutput}
 * @see {@link SearchCommandInput} for command's `input` shape.
 * @see {@link SearchCommandOutput} for command's `response` shape.
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
export class SearchCommand extends $Command
  .classBuilder<
    SearchCommandInput,
    SearchCommandOutput,
    BankSearchClientResolvedConfig,
    ServiceInputTypes,
    ServiceOutputTypes
  >()
  .ep(commonParams)
  .m(function (this: any, Command: any, cs: any, config: BankSearchClientResolvedConfig, o: any) {
    return [getEndpointPlugin(config, Command.getEndpointParameterInstructions())];
  })
  .s("BankSearch", "Search", {})
  .n("BankSearchClient", "SearchCommand")
  .sc(Search$)
  .build() {
  /** @internal type navigation helper, not in runtime. */
  protected declare static __types: {
    api: {
      input: SearchInput;
      output: SearchOutput;
    };
    sdk: {
      input: SearchCommandInput;
      output: SearchCommandOutput;
    };
  };
}
