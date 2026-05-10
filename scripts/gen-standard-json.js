#!/usr/bin/env node
"use strict";

const fs   = require("fs");
const path = require("path");

const contractPath = process.argv[2];
if (!contractPath) {
  console.error("usage: gen-standard-json.js <contract-path>");
  process.exit(1);
}

const root    = process.cwd();
const visited = new Set();
const sources = {};

function resolveImport(importPath, fromDir) {
  if (importPath.startsWith(".")) {
    return path.resolve(fromDir, importPath);
  }
  return path.resolve(root, importPath);
}

function collect(filePath) {
  const rel = path.relative(root, path.resolve(filePath));
  if (visited.has(rel)) return;
  visited.add(rel);

  const content = fs.readFileSync(filePath, "utf8");
  sources[rel]  = { content };

  const dir         = path.dirname(filePath);
  const importRegex = /import\s+(?:\{[^}]*\}\s+from\s+)?["']([^"']+)["']/g;
  let match;
  while ((match = importRegex.exec(content)) !== null) {
    collect(resolveImport(match[1], dir));
  }
}

collect(contractPath);

const contractName = path.basename(contractPath, ".sol");
const outDir       = path.join(root, "build", contractName);
fs.mkdirSync(outDir, { recursive: true });

const input = {
  language: "Solidity",
  sources,
  settings: {
    optimizer:       { enabled: false, runs: 200 },
    outputSelection: { "*": { "*": ["abi", "evm.bytecode", "evm.deployedBytecode", "metadata"] } },
  },
};

const outPath = path.join(outDir, `${contractName}.standard-input.json`);
fs.writeFileSync(outPath, JSON.stringify(input, null, 2));
console.log(`wrote ${outPath}`);