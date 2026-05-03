// smithy-typescript generated code
import { BinaryDecisionDiagram } from "@smithy/util-endpoints";

const _data={
  conditions: [
  ],
  results: [
    [-1],
    ["{Endpoint}",{}]
  ]
};

const root = 100000001;
const r = 100_000_000;
const nodes = new Int32Array([
  -1, 1, -1,
]);
export const bdd = BinaryDecisionDiagram.from(
  nodes, root, _data.conditions, _data.results
);
