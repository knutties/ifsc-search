$version: "2.0"

namespace io.knutties.banksearch

use aws.api#service
use aws.protocols#restJson1
use smithy.rules#endpointRuleSet

/// HTTP search service for Indian bank branches.
//
// The endpointRuleSet trait below replaces the AWS SDK's default
// region/FIPS/DualStack endpoint params with a single `endpoint` builtIn,
// which is all this non-AWS service needs. The generated client therefore
// drops region/useFipsEndpoint/useDualstackEndpoint from its public input
// types. Note: this only cleans up the *type-level* surface — the runtime
// region resolver in `runtimeConfig.ts` is still emitted by the AWS-flavored
// codegen, and is patched out post-generation by
// `smithy/scripts/patch-runtime-config.mjs`.
@service(sdkId: "BankSearch")
@restJson1
@endpointRuleSet(
    version: "1.0"
    parameters: {
        Endpoint: {
            type: "string"
            builtIn: "SDK::Endpoint"
            required: true
            documentation: "The HTTP endpoint to send requests to."
        }
    }
    rules: [
        {
            conditions: []
            endpoint: { url: "{Endpoint}" }
            type: "endpoint"
        }
    ]
)
service BankSearch {
    version: "2026-05-02"
    resources: [Branch]
    operations: [
        ListBanks
        Healthz
        Status
    ]
    errors: [BadRequest]
}

@error("client")
@httpError(400)
structure BadRequest {
    @required
    error: String
}
