$version: "2.0"

namespace io.knutties.banksearch

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
