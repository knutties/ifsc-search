import { suite, test } from "node:test";
import { strict as assert } from "node:assert";

import {
    BankSearchClient,
    BadRequest,
    BranchNotFound,
    GetBranchCommand,
    HealthzCommand,
    ListBanksCommand,
    SearchCommand,
    StatusCommand,
} from "@knutties/bank-search-client";

const baseUrl = process.env.BANK_SEARCH_E2E_BASE_URL ?? "http://localhost:8080";

const client = new BankSearchClient({ endpoint: baseUrl });

suite("bank-search e2e via SDK", () => {
    test("healthz returns ok", async () => {
        const res = await client.send(new HealthzCommand({}));
        assert.equal(res.status, "ok");
    });

    test("status reports indexed_docs", async () => {
        const res = await client.send(new StatusCommand({}));
        assert.ok(typeof res.indexed_docs === "number" && res.indexed_docs > 0);
    });

    test("getBranch returns the requested IFSC", async () => {
        const res = await client.send(new GetBranchCommand({ ifsc: "HDFC0000001" }));
        assert.equal(res.ifsc, "HDFC0000001");
        assert.equal(res.bank_code, "HDFC");
        assert.equal(res.branch, "ANDHERI WEST");
        assert.equal(res.state, "MAHARASHTRA");
    });

    test("getBranch on unknown IFSC throws BranchNotFound (404)", async () => {
        await assert.rejects(
            () => client.send(new GetBranchCommand({ ifsc: "ZZZZ0000000" })),
            (err: unknown) => err instanceof BranchNotFound,
        );
    });

    test("search by bank+q matches HDFC Andheri", async () => {
        const res = await client.send(new SearchCommand({ bank: "HDFC", q: "andheri" }));
        assert.ok(res.total && res.total >= 1);
        assert.ok(res.results?.some((r) => r.ifsc === "HDFC0000001"));
    });

    test("search by IFSC prefix is case-insensitive", async () => {
        const res = await client.send(new SearchCommand({ ifsc: "hdfc0" }));
        assert.ok(res.total && res.total >= 1);
        assert.ok(res.results?.every((r) => r.ifsc?.startsWith("HDFC0")));
    });

    test("search by state is strict and case-insensitive", async () => {
        const res = await client.send(new SearchCommand({ state: "Karnataka" }));
        assert.ok(res.total && res.total >= 1);
        assert.ok(res.results?.every((r) => r.state === "KARNATAKA"));
        assert.ok(res.results?.some((r) => r.ifsc === "HDFC0000003"));
    });

    test("search by city does not bleed across substrings", async () => {
        const res = await client.send(new SearchCommand({ city: "Bangalore" }));
        assert.ok(res.total && res.total >= 1);
        assert.ok(res.results?.every((r) => r.city === "BANGALORE"));
    });

    test("search with no signals throws BadRequest (400)", async () => {
        await assert.rejects(
            () => client.send(new SearchCommand({})),
            (err: unknown) => err instanceof BadRequest,
        );
    });

    test("listBanks contains seeded banks", async () => {
        const res = await client.send(new ListBanksCommand({}));
        const codes = (res.banks ?? []).map((b) => b.bank_code);
        assert.ok(codes.includes("HDFC"));
        assert.ok(codes.includes("SBIN"));
    });
});
