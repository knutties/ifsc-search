$version: "2.0"

namespace io.knutties.banksearch

use aws.api#service
use aws.protocols#restJson1

/// HTTP search service for Indian bank branches.
@service(sdkId: "BankSearch")
@restJson1
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
