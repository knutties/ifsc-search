// Post-process the generated TypeScript client so consumers don't have
// to pass a `region` to the constructor. The AWS-flavored Smithy
// runtime resolves region eagerly (used to build
// `service.<region>.amazonaws.com` endpoints) and throws "Region is
// missing" if neither the constructor config nor env / shared config
// supply one — even when an explicit `endpoint` is set.
//
// bank-search has no use for region, so we replace the resolver
// defaults with a static fallback string in both runtime variants.
//
// Run by `make smithy-publish` after the rsync step.

import { readFileSync, writeFileSync } from "node:fs";

const FALLBACK = "auto";

const replacements = [
    {
        file: "clients/typescript/src/runtimeConfig.ts",
        // 4-line `loadNodeConfig(...)` block — ends at the first `),`
        // we encounter after the opening call.
        pattern: /region: config\?\.region \?\? loadNodeConfig\([\s\S]*?\),/,
        replacement: `region: config?.region ?? "${FALLBACK}",`,
    },
    {
        file: "clients/typescript/src/runtimeConfig.browser.ts",
        pattern: /region: config\?\.region \?\? invalidProvider\("Region is missing"\),/,
        replacement: `region: config?.region ?? "${FALLBACK}",`,
    },
];

for (const { file, pattern, replacement } of replacements) {
    const before = readFileSync(file, "utf8");
    const after = before.replace(pattern, replacement);
    if (after === before) {
        console.error(`patch-runtime-config: no match in ${file} — generator output may have changed`);
        process.exit(1);
    }
    writeFileSync(file, after);
    console.log(`patched ${file}`);
}
