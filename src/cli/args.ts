import { readFileSync } from "node:fs";

export type CliOptions = {
  json: boolean;
  currency: "USD";
  dynadotKey: string | undefined;
  timeoutMs: number;
  domains: string[];
};

const usage = `Usage:
  domain search <domains...>
  domain search --file domains.txt

Flags:
  --json
  --currency USD
  --dynadot-key <key>
  --timeout 10s`;

const parseTimeout = (raw: string): number => {
  const seconds = raw.match(/^([0-9]+)s$/i);
  if (seconds?.[1]) return Number(seconds[1]) * 1000;

  const millis = raw.match(/^([0-9]+)ms$/i);
  if (millis?.[1]) return Number(millis[1]);

  const direct = Number(raw);
  if (Number.isFinite(direct) && direct > 0) return direct;

  throw new Error(`Invalid timeout value: ${raw}`);
};

const fromFile = (filePath: string): string[] =>
  readFileSync(filePath, "utf-8")
    .split(/\r?\n/)
    .map((line) => line.trim())
    .filter((line) => line.length > 0);

export const parseArgs = (argv: string[]): CliOptions => {
  const [command, ...rest] = argv;
  if (command !== "search") {
    throw new Error(usage);
  }

  const options: Omit<CliOptions, "domains"> = {
    json: false,
    currency: "USD",
    dynadotKey: undefined,
    timeoutMs: 10_000,
  };

  const positional: string[] = [];
  let filePath: string | null = null;

  for (let index = 0; index < rest.length; index += 1) {
    const token = rest[index];
    if (!token) continue;

    if (token === "--json") {
      options.json = true;
      continue;
    }

    if (token === "--file") {
      filePath = rest[index + 1] ?? null;
      index += 1;
      continue;
    }

    if (token === "--currency") {
      const value = rest[index + 1];
      if (!value || value !== "USD") throw new Error("Only USD is supported in v1");
      options.currency = value;
      index += 1;
      continue;
    }

    if (token === "--dynadot-key") {
      const value = rest[index + 1];
      if (!value) throw new Error("--dynadot-key requires a value");
      options.dynadotKey = value;
      index += 1;
      continue;
    }

    if (token === "--timeout") {
      const value = rest[index + 1];
      if (!value) throw new Error("--timeout requires a value");
      options.timeoutMs = parseTimeout(value);
      index += 1;
      continue;
    }

    positional.push(token);
  }

  const fromPath = filePath ? fromFile(filePath) : [];
  const domains = [
    ...new Set([...positional, ...fromPath].map((value) => value.trim().toLowerCase())),
  ].filter((value) => value.length > 0);

  if (domains.length === 0) throw new Error("No domains provided");

  return { ...options, domains };
};
