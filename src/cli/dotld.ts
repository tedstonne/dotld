#!/usr/bin/env bun

import { runCli } from "@/cli/index";

const withImplicitSearch = (argv: string[]): string[] => {
  if (argv.length === 0) return argv;

  const first = argv[0];
  if (!first) return argv;
  if (first === "search" || first === "--help" || first === "-h") return argv;

  return ["search", ...argv];
};

const run = async (): Promise<void> => {
  const code = await runCli(withImplicitSearch(process.argv.slice(2)));
  process.exit(code);
};

void run();
