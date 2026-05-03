// smithy-typescript generated code
import { getEndpointPlugin } from "@smithy/middleware-endpoint";
import { Command as $Command } from "@smithy/smithy-client";
import type { MetadataBearer as __MetadataBearer } from "@smithy/types";

import type { BankSearchClientResolvedConfig, ServiceInputTypes, ServiceOutputTypes } from "../BankSearchClient";
import { commonParams } from "../endpoint/EndpointParameters";
import type { HealthzOutput } from "../models/models_0";
import { Healthz$ } from "../schemas/schemas_0";

/**
 * @public
 */
export type { __MetadataBearer };
export { $Command };
/**
 * @public
 *
 * The input for {@link HealthzCommand}.
 */
export interface HealthzCommandInput {}
/**
 * @public
 *
 * The output of {@link HealthzCommand}.
 */
export interface HealthzCommandOutput extends HealthzOutput, __MetadataBearer {}

/**
 * Lightweight liveness probe for load balancers.
 * @example
 * Use a bare-bones client and the command you need to make an API call.
 * ```javascript
 * import { BankSearchClient, HealthzCommand } from "@knutties/bank-search-client"; // ES Modules import
 * // const { BankSearchClient, HealthzCommand } = require("@knutties/bank-search-client"); // CommonJS import
 * // import type { BankSearchClientConfig } from "@knutties/bank-search-client";
 * const config = {}; // type is BankSearchClientConfig
 * const client = new BankSearchClient(config);
 * const input = {};
 * const command = new HealthzCommand(input);
 * const response = await client.send(command);
 * // { // HealthzOutput
 * //   status: "STRING_VALUE", // required
 * // };
 *
 * ```
 *
 * @param HealthzCommandInput - {@link HealthzCommandInput}
 * @returns {@link HealthzCommandOutput}
 * @see {@link HealthzCommandInput} for command's `input` shape.
 * @see {@link HealthzCommandOutput} for command's `response` shape.
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
export class HealthzCommand extends $Command
  .classBuilder<
    HealthzCommandInput,
    HealthzCommandOutput,
    BankSearchClientResolvedConfig,
    ServiceInputTypes,
    ServiceOutputTypes
  >()
  .ep(commonParams)
  .m(function (this: any, Command: any, cs: any, config: BankSearchClientResolvedConfig, o: any) {
    return [getEndpointPlugin(config, Command.getEndpointParameterInstructions())];
  })
  .s("BankSearch", "Healthz", {})
  .n("BankSearchClient", "HealthzCommand")
  .sc(Healthz$)
  .build() {
  /** @internal type navigation helper, not in runtime. */
  protected declare static __types: {
    api: {
      input: {};
      output: HealthzOutput;
    };
    sdk: {
      input: HealthzCommandInput;
      output: HealthzCommandOutput;
    };
  };
}
