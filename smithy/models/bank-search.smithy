$version: "2.0"

namespace io.knutties.banksearch

use aws.api#service
use aws.protocols#restJson1
use smithy.rules#endpointRuleSet

// ---------------------------------------------------------------------------
// Service
// ---------------------------------------------------------------------------

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
    resources: [Branch, Bank]
    operations: [
        Healthz
        Status
    ]
    errors: [BadRequest]
}

// ---------------------------------------------------------------------------
// Branch resource — single-record lookup + search
// ---------------------------------------------------------------------------

/// A bank branch keyed by its IFSC code. The properties are the wire
/// fields exposed by every operation that returns a branch.
resource Branch {
    identifiers: {
        ifsc: String
    }

    properties: {
        bank_code: String
        bank_name: String
        branch: String
        centre: String
        district: String
        state: String
        address: String
        city: String
        contact: String
        micr: String
        swift: String
        upi: Boolean
        neft: Boolean
        rtgs: Boolean
        imps: Boolean
    }

    read: GetBranch
    collectionOperations: [Search]
}

/// Look up a single branch by IFSC code.
@readonly
@http(method: "GET", uri: "/ifsc/{ifsc}", code: 200)
operation GetBranch {
    input := for Branch {
        @httpLabel
        @required
        $ifsc
    }

    output := @references([{ resource: Bank }]) for Branch {
        @required $ifsc
        @required $bank_code
        @required $bank_name
        @required $branch
        $centre
        $district
        $state
        $address
        $city
        $contact
        $micr
        $swift
        $upi
        $neft
        $rtgs
        $imps
    }

    errors: [BranchNotFound]
}

@error("client")
@httpError(404)
structure BranchNotFound {
    @required
    error: String
}

/// Free-text plus structured search across the indexed branches.
/// At least one of {bank, q, ifsc, state, district, city} is required.
@readonly
@http(method: "GET", uri: "/search", code: 200)
operation Search {
    input := {
        /// 4-char IFSC bank code or fuzzy bank name.
        @httpQuery("bank")
        bank: String

        /// Free-text query over branch, city, address and IFSC code prefix.
        @httpQuery("q")
        q: String

        /// Case-insensitive IFSC prefix, e.g. "HDFC0CAG".
        @httpQuery("ifsc")
        ifsc: String

        /// Strict, case-insensitive state filter.
        @httpQuery("state")
        state: String

        /// Strict, case-insensitive district filter.
        @httpQuery("district")
        district: String

        /// Strict, case-insensitive city filter.
        @httpQuery("city")
        city: String

        /// Page size; defaults to 20 and is clamped to 100.
        @httpQuery("limit")
        limit: Integer

        /// Result offset for pagination; defaults to 0.
        @httpQuery("offset")
        offset: Integer
    }

    output := {
        @required total: Integer
        @required limit: Integer
        @required offset: Integer
        @required results: ResultItemList
    }

    errors: [BadRequest]
}

list ResultItemList {
    member: ResultItem
}

/// A search hit — a Branch's properties plus its relevance score.
@references([
    { resource: Bank }
])
structure ResultItem for Branch {
    @required $ifsc
    @required $bank_code
    @required $bank_name
    @required $branch
    $centre
    $district
    $state
    $address
    $city
    $contact
    $micr
    $swift
    $upi
    $neft
    $rtgs
    $imps

    @required
    score: Double
}

// ---------------------------------------------------------------------------
// Bank resource — currently list-only, ready for Get/Put/Delete
// ---------------------------------------------------------------------------

/// A bank registered with the index. Identified by its 4-char IFSC bank
/// code. Today the only bound operation is List; future Get/Put/Delete
/// will hang off this resource without changing existing URLs.
//
// `properties` is intentionally omitted: Smithy requires every declared
// property to be referenced by a create or instance (read/update/delete)
// operation, and we have none yet. When GetBank lands, `bank_name` should
// move from BankSummary into a `properties` block here, and BankSummary
// can be re-projected from the resource via `for Bank` + `$bank_name`.
resource Bank {
    identifiers: {
        bank_code: String
    }

    list: ListBanks
}

/// List the distinct banks present in the index, sorted by bank_code.
@readonly
@http(method: "GET", uri: "/list", code: 200)
operation ListBanks {
    output := {
        @required
        total: Integer

        @required
        banks: BankList
    }
}

list BankList {
    member: BankSummary
}

/// Wire shape for a bank in the listing.
structure BankSummary {
    @required
    bank_code: String

    @required
    bank_name: String
}

// ---------------------------------------------------------------------------
// Operational probes
// ---------------------------------------------------------------------------

/// Lightweight liveness probe for load balancers.
@readonly
@http(method: "GET", uri: "/healthz", code: 200)
operation Healthz {
    output := {
        @required
        status: String
    }
}

/// Index version metadata and document count.
@readonly
@http(method: "GET", uri: "/status", code: 200)
operation Status {
    output := {
        @required
        status: String

        indexed_docs: Long
        release_tag: String
        rbi_update_date: String
        built_at: String
    }
}

// ---------------------------------------------------------------------------
// Shared errors
// ---------------------------------------------------------------------------

@error("client")
@httpError(400)
structure BadRequest {
    @required
    error: String
}
