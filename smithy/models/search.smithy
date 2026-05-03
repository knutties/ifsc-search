$version: "2.0"

namespace io.knutties.banksearch

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
