#!/usr/bin/env bun

import { spawnSync } from "node:child_process";
import { mkdirSync, rmSync, symlinkSync } from "node:fs";
import { homedir } from "node:os";
import { join } from "node:path";

const run = (command: string, args: string[]): string => {
  const result = spawnSync(command, args, {
    stdio: ["ignore", "pipe", "inherit"],
    encoding: "utf-8",
  });

  if (result.status !== 0) {
    process.exit(result.status ?? 1);
  }

  return (result.stdout ?? "").trim();
};

const runWithOutput = (command: string, args: string[]): void => {
  const result = spawnSync(command, args, {
    stdio: "inherit",
    encoding: "utf-8",
  });

  if (result.status !== 0) {
    process.exit(result.status ?? 1);
  }
};

runWithOutput("bun", ["add", "-g", process.cwd()]);

const bunGlobalBin = run("bun", ["pm", "bin", "-g"]);
if (!bunGlobalBin) {
  process.stderr.write("Could not determine Bun global bin path.\n");
  process.exit(1);
}

const source = join(bunGlobalBin, "dotld");
const localBinDir = join(homedir(), ".local", "bin");
const target = join(localBinDir, "dotld");

mkdirSync(localBinDir, { recursive: true });
rmSync(target, { force: true });
symlinkSync(source, target);

process.stdout.write(`Linked dotld to ${target}\n`);

const pathEntries = (process.env.PATH ?? "").split(":");
if (!pathEntries.includes(localBinDir)) {
  process.stdout.write(`Add ${localBinDir} to PATH if dotld is not found.\n`);
  process.stdout.write(`zsh/bash: export PATH=\"${localBinDir}:$PATH\"\n`);
  process.stdout.write(`fish: fish_add_path \"${localBinDir}\"\n`);
}
