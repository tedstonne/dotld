#!/usr/bin/env bun

import { spawnSync } from "node:child_process";
import { rmSync } from "node:fs";
import { homedir } from "node:os";
import { join } from "node:path";

const run = (command: string, args: string[]): void => {
  const result = spawnSync(command, args, {
    stdio: "inherit",
    encoding: "utf-8",
  });

  if (result.status !== 0) {
    process.exit(result.status ?? 1);
  }
};

run("bun", ["remove", "-g", "domain-cli-project"]);

const target = join(homedir(), ".local", "bin", "dotld");
rmSync(target, { force: true });

process.stdout.write(`Removed ${target}\n`);
