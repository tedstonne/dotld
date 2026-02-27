#!/usr/bin/env bun

import { $ } from "bun";

const IMAGE = "dotld-smoke";

console.log("Building test image...\n");
await $`docker build -f Dockerfile.test -t ${IMAGE} .`;

console.log("\nRunning smoke tests...\n");
await $`docker run --rm -e DYNADOT_API_PRODUCTION_KEY ${IMAGE}`;
