// Post-process the generated TypeScript client's runtimeConfig.
//
// The AWS-flavored Smithy runtime resolves `region` eagerly (used to
// build `service.<region>.amazonaws.com` endpoints) and throws "Region
// is missing" if neither the constructor config nor env / shared
// config supply one — even when an explicit `endpoint` is set.
// bank-search has no use for region, so we replace the resolver
// defaults with a static fallback string in both runtime variants.
//
// The README is not generated — clients/typescript/README.md is
// hand-maintained and is excluded from the codegen rsync in the
// Makefile.
//
// Run by `make smithy-publish` after the rsync step.

import { readFileSync, writeFileSync } from "node:fs";

const FALLBACK_REGION = "auto";

function applyPatch(content, pattern, replacement, label) {
    if (pattern instanceof RegExp) {
        if (!pattern.test(content)) {
            console.error(`post-process: no match for ${label} — generator output may have changed`);
            process.exit(1);
        }
        return content.replace(pattern, replacement);
    }
    const parts = content.split(pattern);
    if (parts.length !== 2) {
        console.error(`post-process: expected exactly 1 occurrence of ${label}, found ${parts.length - 1}`);
        process.exit(1);
    }
    return parts.join(replacement);
}

const files = [
    {
        path: "clients/typescript/src/runtimeConfig.ts",
        patches: [
            {
                label: "node region resolver",
                // 4-line `loadNodeConfig(...)` block — ends at the first `),`
                // we encounter after the opening call.
                pattern: /region: config\?\.region \?\? loadNodeConfig\([\s\S]*?\),/,
                replacement: `region: config?.region ?? "${FALLBACK_REGION}",`,
            },
        ],
    },
    {
        path: "clients/typescript/src/runtimeConfig.browser.ts",
        patches: [
            {
                label: "browser region resolver",
                pattern: /region: config\?\.region \?\? invalidProvider\("Region is missing"\),/,
                replacement: `region: config?.region ?? "${FALLBACK_REGION}",`,
            },
        ],
    },
];

for (const { path, patches } of files) {
    let content = readFileSync(path, "utf8");
    for (const { label, pattern, replacement } of patches) {
        content = applyPatch(content, pattern, replacement, `${path}: ${label}`);
    }
    writeFileSync(path, content);
    console.log(`patched ${path}`);
}
