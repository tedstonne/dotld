#!/usr/bin/env bun

import { spawnSync } from "node:child_process";
import { readFileSync, writeFileSync } from "node:fs";

type Bump = "patch" | "minor" | "major";

type Config = {
  version: string | null;
  bump: Bump;
  notes: string | null;
  dryRun: boolean;
};

const run = (command: string, args: string[], message: string): void => {
  process.stdout.write(`${message}\n`);
  const result = spawnSync(command, args, {
    stdio: "inherit",
    encoding: "utf-8",
  });

  if (result.status !== 0) {
    process.exit(result.status ?? 1);
  }
};

const output = (command: string, args: string[]): string => {
  const result = spawnSync(command, args, {
    stdio: ["ignore", "pipe", "pipe"],
    encoding: "utf-8",
  });

  if (result.status !== 0) {
    process.stderr.write(result.stderr ?? "\n");
    process.exit(result.status ?? 1);
  }

  return (result.stdout ?? "").trim();
};

const parseArgs = (): Config => {
  const args = process.argv.slice(2);
  let version: string | null = null;
  let bump: Bump = "patch";
  let notes: string | null = null;
  let dryRun = false;

  for (let index = 0; index < args.length; index += 1) {
    const token = args[index];
    if (!token) continue;

    if (token === "--bump") {
      const value = args[index + 1];
      if (value !== "patch" && value !== "minor" && value !== "major") {
        process.stderr.write("--bump must be patch, minor, or major\n");
        process.exit(1);
      }
      bump = value;
      index += 1;
      continue;
    }

    if (token === "--version") {
      const value = args[index + 1];
      if (!value) {
        process.stderr.write("--version requires a value\n");
        process.exit(1);
      }
      version = value;
      index += 1;
      continue;
    }

    if (token === "--notes") {
      const value = args[index + 1];
      if (!value) {
        process.stderr.write("--notes requires a value\n");
        process.exit(1);
      }
      notes = value;
      index += 1;
      continue;
    }

    if (token === "--dry-run") {
      dryRun = true;
      continue;
    }

    if (token.startsWith("-")) {
      process.stderr.write(`Unknown flag: ${token}\n`);
      process.exit(1);
    }

    version = token;
  }

  return { version, bump, notes, dryRun };
};

const parseSemver = (value: string): [number, number, number] => {
  const match = value.match(/^v?(\d+)\.(\d+)\.(\d+)$/);
  if (!match?.[1] || !match[2] || !match[3]) {
    process.stderr.write(`Invalid semantic version: ${value}\n`);
    process.exit(1);
  }

  return [Number(match[1]), Number(match[2]), Number(match[3])];
};

const latestTag = (): string | null => {
  const tags = output("git", ["tag", "--list", "v*", "--sort=-v:refname"]);
  if (tags.length === 0) return null;

  return tags.split("\n")[0] ?? null;
};

const packageVersion = (): string => {
  const content = readFileSync("package.json", "utf-8");
  const parsed = JSON.parse(content) as { version?: string };
  if (!parsed.version) {
    process.stderr.write("package.json is missing version\n");
    process.exit(1);
  }

  return parsed.version;
};

const nextVersion = (base: string, bump: Bump): string => {
  const [major, minor, patch] = parseSemver(base);
  if (bump === "major") return `${major + 1}.0.0`;
  if (bump === "minor") return `${major}.${minor + 1}.0`;

  return `${major}.${minor}.${patch + 1}`;
};

const ensureCleanGit = (): void => {
  const status = output("git", ["status", "--porcelain"]);
  if (status.length > 0) {
    process.stderr.write(
      "Git working tree is not clean. Commit or stash changes before release.\n",
    );
    process.exit(1);
  }
};

const ensureTagMissing = (tag: string): void => {
  const existing = spawnSync("git", ["rev-parse", "-q", "--verify", `refs/tags/${tag}`], {
    stdio: "ignore",
  });
  if (existing.status === 0) {
    process.stderr.write(`Tag ${tag} already exists. Choose a new version.\n`);
    process.exit(1);
  }
};

const ensureGhAvailable = (): void => {
  const ghCheck = spawnSync("gh", ["--version"], { stdio: "ignore" });
  if (ghCheck.status !== 0) {
    process.stderr.write("GitHub CLI (gh) is required for release automation.\n");
    process.exit(1);
  }
};

const commitSubjectsSince = (tag: string | null): string[] => {
  const range = tag ? `${tag}..HEAD` : "HEAD";
  const raw = output("git", ["log", "--pretty=%s", range]);
  if (raw.length === 0) return [];

  return raw
    .split("\n")
    .map((line) => line.trim())
    .filter((line) => line.length > 0);
};

const updatePackageVersion = (version: string): void => {
  const content = readFileSync("package.json", "utf-8");
  const parsed = JSON.parse(content) as Record<string, unknown>;
  parsed.version = version;
  writeFileSync("package.json", `${JSON.stringify(parsed, null, 2)}\n`, "utf-8");
};

const changelogSection = (tag: string, subjects: string[]): string => {
  const today = new Date().toISOString().slice(0, 10);
  const bullets = subjects.map((subject) => `- ${subject}`).join("\n");

  return `## ${tag} - ${today}\n${bullets}`;
};

const updateChangelog = (section: string): void => {
  let current = "";
  try {
    current = readFileSync("CHANGELOG.md", "utf-8").trim();
  } catch {
    current = "# Changelog";
  }

  const header = "# Changelog";
  const normalized = current.startsWith(header) ? current : `${header}\n\n${current}`;
  const rest = normalized.replace(/^# Changelog\s*/, "").trim();
  const next = `${header}\n\n${section}${rest.length > 0 ? `\n\n${rest}` : ""}\n`;

  writeFileSync("CHANGELOG.md", next, "utf-8");
};

const createRelease = (tag: string, notes: string): void => {
  const assets = [
    "dist/dotld-linux-x64",
    "dist/dotld-linux-arm64",
    "dist/dotld-darwin-x64",
    "dist/dotld-darwin-arm64",
    "dist/checksums.txt",
  ];

  run(
    "gh",
    ["release", "create", tag, ...assets, "--title", tag, "--notes", notes],
    `Creating GitHub release ${tag}...`,
  );
};

const main = (): void => {
  const config = parseArgs();
  ensureGhAvailable();
  ensureCleanGit();

  const previousTag = latestTag();
  const baseVersion = config.version ?? nextVersion(previousTag ?? packageVersion(), config.bump);
  const [major, minor, patch] = parseSemver(baseVersion);
  const version = `${major}.${minor}.${patch}`;
  const tag = `v${version}`;

  ensureTagMissing(tag);

  const subjects = commitSubjectsSince(previousTag);
  if (subjects.length === 0) {
    process.stderr.write("No commits found since the previous release.\n");
    process.exit(1);
  }

  const generatedNotes = changelogSection(tag, subjects);

  if (config.dryRun) {
    process.stdout.write("Dry run\n");
    process.stdout.write(`- previous tag: ${previousTag ?? "(none)"}\n`);
    process.stdout.write(`- next version: ${version}\n`);
    process.stdout.write(`- next tag: ${tag}\n`);
    process.stdout.write(`- commits in release: ${subjects.length}\n`);
    process.stdout.write(`- first note line: ${generatedNotes.split("\n")[0] ?? ""}\n`);
    return;
  }

  run("bun", ["run", "lint"], "Running lint...");
  run("bun", ["run", "check"], "Running type-check...");
  run("bun", ["run", "test"], "Running tests...");
  run("bun", ["run", "build:release"], "Building release binaries...");

  updatePackageVersion(version);
  updateChangelog(generatedNotes);

  run("git", ["add", "package.json", "CHANGELOG.md"], "Staging release metadata...");
  run("git", ["commit", "-m", `release ${tag}`], `Committing release ${tag}...`);
  run("git", ["tag", tag], `Creating git tag ${tag}...`);
  run("git", ["push", "origin", "main"], "Pushing main...");
  run("git", ["push", "origin", tag], `Pushing tag ${tag}...`);
  createRelease(tag, config.notes ?? generatedNotes);

  process.stdout.write(`Release complete: ${tag}\n`);
};

main();
