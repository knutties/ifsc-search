<!-- generated file, do not edit directly -->

# @knutties/bank-search-client

## Description

HTTP search service for Indian bank branches.

## Installing

To install this package, use the CLI of your favorite package manager:

- `npm install @knutties/bank-search-client`
- `yarn add @knutties/bank-search-client`
- `pnpm add @knutties/bank-search-client`

## Getting Started

### Import

The client is modular by clients and commands.
To send a request, you only need to import the `BankSearchClient` and
the commands you need, for example `ListBanksCommand`:

```js
// ES5 example
const { BankSearchClient, ListBanksCommand } = require("@knutties/bank-search-client");
```

```ts
// ES6+ example
import { BankSearchClient, ListBanksCommand } from "@knutties/bank-search-client";
```

### Usage

To send a request:

- Instantiate a client with configuration (e.g. endpoint).
- Instantiate a command with input parameters.
- Call the `send` operation on the client, providing the command object as input.

```js
const client = new BankSearchClient({ endpoint: "https://your-bank-search-host" });

const params = { /** input parameters */ };
const command = new ListBanksCommand(params);
```

#### Async/await

We recommend using the [await](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Operators/await)
operator to wait for the promise returned by send operation as follows:

```js
// async/await.
try {
  const data = await client.send(command);
  // process data.
} catch (error) {
  // error handling.
} finally {
  // finally.
}
```

#### Promises

You can also use [Promise chaining](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Guide/Using_promises#chaining).

```js
client
  .send(command)
  .then((data) => {
    // process data.
  })
  .catch((error) => {
    // error handling.
  })
  .finally(() => {
    // finally.
  });
```

#### Aggregated client

The aggregated client class is exported from the same package, but without the "Client" suffix.

`BankSearch` extends `BankSearchClient` and additionally supports all operations, waiters, and paginators as methods.
If you are bundling, prefer the bare-bones client (`BankSearchClient`).

```ts
import { BankSearch } from "@knutties/bank-search-client";

const client = new BankSearch({ endpoint: "https://your-bank-search-host" });

// async/await.
try {
  const data = await client.listBanks(params);
  // process data.
} catch (error) {
  // error handling.
}

// Promises.
client
  .listBanks(params)
  .then((data) => {
    // process data.
  })
  .catch((error) => {
    // error handling.
  });

// callbacks (not recommended).
client.listBanks(params, (err, data) => {
  // process err and data.
});
```

### Troubleshooting

When the service returns an exception, the error will include the exception information,
as well as response metadata (e.g. request id).

```js
try {
  const data = await client.send(command);
  // process data.
} catch (error) {
  const { requestId, cfId, extendedRequestId } = error.$metadata;
  console.log({ requestId, cfId, extendedRequestId });
  /**
   * The keys within exceptions are also parsed.
   * You can access them by specifying exception names:
   * if (error.name === 'SomeServiceException') {
   *     const value = error.specialKeyInException;
   * }
   */
}
```

## Getting Help

Please [open an issue](https://github.com/knutties/bank-search/issues) on the bank-search repo for bugs or feature requests.

## Contributing

This client code is generated automatically. Any modifications will be overwritten the next time the `@knutties/bank-search-client` package is updated.
The Smithy IDL and codegen config live in the `smithy/` directory of the [bank-search repo](https://github.com/knutties/bank-search).

## License

This SDK is distributed under the
[Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0),
see LICENSE for more information.

