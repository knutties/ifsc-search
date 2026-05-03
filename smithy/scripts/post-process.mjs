// Post-process the generated TypeScript client.
//
// 1. runtimeConfig: the AWS-flavored Smithy runtime resolves `region`
//    eagerly (used to build `service.<region>.amazonaws.com` endpoints)
//    and throws "Region is missing" if neither the constructor config
//    nor env / shared config supply one — even when an explicit
//    `endpoint` is set. bank-search has no use for region, so we
//    replace the resolver defaults with a static fallback string in
//    both runtime variants.
// 2. README: the generator emits an AWS-SDK-flavored README with a
//    `region: "REGION"` quickstart and many links into the
//    aws-sdk-js-v3 repo / docs.aws.amazon.com that don't apply to this
//    service. Patch those out.
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
    {
        path: "clients/typescript/README.md",
        patches: [
            {
                label: "AWS SDK framing in description",
                pattern: "AWS SDK for JavaScript BankSearch Client for Node.js, Browser and React Native.\n\n",
                replacement: "",
            },
            {
                label: "install instructions (point at GitHub Releases tarball)",
                pattern:
                    "## Installing\n\n" +
                    "To install this package, use the CLI of your favorite package manager:\n\n" +
                    "- `npm install @knutties/bank-search-client`\n" +
                    "- `yarn add @knutties/bank-search-client`\n" +
                    "- `pnpm add @knutties/bank-search-client`\n",
                replacement:
                    "## Installing\n\n" +
                    "This package isn't published to npm. Install the latest tarball from " +
                    "[GitHub Releases](https://github.com/knutties/bank-search/releases):\n\n" +
                    "```sh\n" +
                    "# latest release:\n" +
                    "npm install https://github.com/knutties/bank-search/releases/latest/download/bank-search-client.tgz\n\n" +
                    "# pinned to a specific release (replace vYYYYMMDDHHMM with a tag from the releases page):\n" +
                    "npm install https://github.com/knutties/bank-search/releases/download/vYYYYMMDDHHMM/bank-search-client.tgz\n" +
                    "```\n\n" +
                    "`yarn add` / `pnpm add` accept the same tarball URL.\n\n" +
                    "Each release ships with the in-tarball version suffixed by the calver " +
                    "(e.g. `0.1.0-202605031200`) so lockfiles can pin a specific build.\n",
            },
            {
                label: "\"AWS SDK is modulized\" intro",
                pattern: "The AWS SDK is modulized by clients and commands.\n",
                replacement: "The client is modular by clients and commands.\n",
            },
            {
                label: "credentials/region usage bullet + AWS doc links",
                pattern:
                    "- Instantiate a client with configuration (e.g. credentials, region).\n" +
                    "  - See [docs/CLIENTS](https://github.com/aws/aws-sdk-js-v3/blob/main/supplemental-docs/CLIENTS.md) for configuration details.\n" +
                    "  - See [@aws-sdk/config](https://github.com/aws/aws-sdk-js-v3/blob/main/packages/config/README.md) for additional options.\n",
                replacement: "- Instantiate a client with configuration (e.g. endpoint).\n",
            },
            {
                label: "BankSearchClient region example",
                pattern: `const client = new BankSearchClient({ region: "REGION" });`,
                replacement: `const client = new BankSearchClient({ endpoint: "https://your-bank-search-host" });`,
            },
            {
                label: "AWS SDK v2 / bundling blog references",
                pattern:
                    "This style may be familiar to you from the AWS SDK for JavaScript v2.\n\n" +
                    "If you are bundling the AWS SDK, we recommend using only the bare-bones client (`BankSearchClient`).\n" +
                    "More details are in the blog post on\n" +
                    "[modular packages in AWS SDK for JavaScript](https://aws.amazon.com/blogs/developer/modular-packages-in-aws-sdk-for-javascript/).\n",
                replacement: "If you are bundling, prefer the bare-bones client (`BankSearchClient`).\n",
            },
            {
                label: "BankSearch (aggregated) region example",
                pattern: `const client = new BankSearch({ region: "REGION" });`,
                replacement: `const client = new BankSearch({ endpoint: "https://your-bank-search-host" });`,
            },
            {
                label: "AWS error-handling doc link",
                pattern:
                    "\nSee also [docs/ERROR_HANDLING](https://github.com/aws/aws-sdk-js-v3/blob/main/supplemental-docs/ERROR_HANDLING.md).\n",
                replacement: "",
            },
            {
                label: "Getting Help section (AWS resources)",
                pattern:
                    "## Getting Help\n\n" +
                    "Please use these community resources for getting help.\n" +
                    "We use GitHub issues for tracking bugs and feature requests, but have limited bandwidth to address them.\n\n" +
                    "- Visit the [Developer Guide](https://docs.aws.amazon.com/sdk-for-javascript/v3/developer-guide/welcome.html)\n" +
                    "  or [API Reference](https://docs.aws.amazon.com/AWSJavaScriptSDK/v3/latest/index.html).\n" +
                    "- Check out the blog posts tagged with [`aws-sdk-js`](https://aws.amazon.com/blogs/developer/tag/aws-sdk-js/)\n" +
                    "  on AWS Developer Blog.\n" +
                    "- Ask a question on [StackOverflow](https://stackoverflow.com/questions/tagged/aws-sdk-js) and tag it with `aws-sdk-js`.\n" +
                    "- Join the AWS JavaScript community on [gitter](https://gitter.im/aws/aws-sdk-js-v3).\n" +
                    "- If it turns out that you may have found a bug, please [open an issue](https://github.com/aws/aws-sdk-js-v3/issues/new/choose).\n\n" +
                    "To test your universal JavaScript code in Node.js, browser and react-native environments,\n" +
                    "visit our [code samples repo](https://github.com/aws-samples/aws-sdk-js-tests).\n",
                replacement:
                    "## Getting Help\n\n" +
                    "Please [open an issue](https://github.com/knutties/bank-search/issues) on the bank-search repo for bugs or feature requests.\n",
            },
            {
                label: "Contributing section AWS scripts link",
                pattern:
                    "To contribute to client you can check our [generate clients scripts](https://github.com/aws/aws-sdk-js-v3/tree/main/scripts/generate-clients).\n",
                replacement:
                    "The Smithy IDL and codegen config live in the `smithy/` directory of the [bank-search repo](https://github.com/knutties/bank-search).\n",
            },
            {
                label: "Operations List section (broken AWS doc links)",
                // Drop everything from this heading to end of file.
                pattern: /\n## Client Commands \(Operations List\)\n[\s\S]*$/,
                replacement: "\n",
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
