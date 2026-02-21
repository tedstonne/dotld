#!/usr/bin/env bun

import { parseArgs } from "@/cli/args";
import { createQueryPlan } from "@/cli/plan";
import { renderTable } from "@/cli/table";
import { searchDynadot } from "@/domain/dynadot";
import ora from "ora";

const ANSI_RESET = "\x1b[0m";
const DYNADOT_KEY_URL = "https://www.dynadot.com/account/domain/setting/api.html";
const useAnsiColor = process.env.NO_COLOR === undefined;

const redText = (value: string): string =>
  useAnsiColor ? `\x1b[1;38;2;255;95;95m${value}\x1b[0m` : value;

const resolveDynadotKey = (
  fromFlag: string | undefined,
): { key: string; source: "flag" | "env" } => {
  if (fromFlag && fromFlag.length > 0) return { key: fromFlag, source: "flag" };

  const fromEnv = process.env.DYNADOT_API_PRODUCTION_KEY;
  if (fromEnv && fromEnv.length > 0) return { key: fromEnv, source: "env" };

  throw new Error("Missing Dynadot key.");
};

const keyWarnings = (rawKey: string, source: "flag" | "env"): string[] => {
  const warnings: string[] = [];
  const trimmed = rawKey.trim();

  if (trimmed !== rawKey) {
    warnings.push("Key has leading/trailing whitespace; trim the value before using it.");
  }

  if (/^[a-f0-9]{64}$/i.test(trimmed)) {
    warnings.push(
      "Key looks like a secret/signing token, not the Dynadot production API key from Tools -> API.",
    );
  }

  if (!/^[a-z0-9]+$/i.test(trimmed)) {
    warnings.push("Key contains unusual characters; Dynadot API keys are typically alphanumeric.");
  }

  if (trimmed.length < 16) {
    warnings.push("Key looks too short; confirm you pasted the full production API key.");
  }

  if (warnings.length > 0) {
    warnings.push(`Source: ${source}. Export DYNADOT_API_PRODUCTION_KEY or pass --dynadot-key.`);
    warnings.push("Fix key:");
    warnings.push(DYNADOT_KEY_URL);
  }

  return warnings;
};

export const runCli = async (argv: string[]): Promise<number> => {
  let spinner: ReturnType<typeof ora> | null = null;
  let spinnerTimer: ReturnType<typeof setTimeout> | null = null;
  let parsed: ReturnType<typeof parseArgs> | null = null;
  let plan: ReturnType<typeof createQueryPlan> | null = null;

  try {
    parsed = parseArgs(argv);
    plan = createQueryPlan(parsed.domains);
    const resolved = resolveDynadotKey(parsed.dynadotKey);
    const dynadotKey = resolved.key.trim();
    const warnings = keyWarnings(resolved.key, resolved.source);
    if (warnings.length > 0) {
      process.stderr.write("Warning: possible key format issue\n");
      for (const warning of warnings) {
        process.stderr.write(`- ${warning}\n`);
      }
    }
    const spinnerEnabled = process.stderr.isTTY && !parsed.json;
    if (spinnerEnabled) {
      spinner = ora({
        text: "",
        isEnabled: spinnerEnabled,
      });

      spinnerTimer = setTimeout(() => {
        spinner?.start();
      }, 120);
    }

    const results = await searchDynadot({
      apiKey: dynadotKey,
      domains: plan.lookupDomains,
      currency: parsed.currency,
      timeoutMs: parsed.timeoutMs,
      affiliateTemplate: process.env.AFFILIATE_URL_TEMPLATE,
    });

    if (spinnerTimer) clearTimeout(spinnerTimer);
    spinner?.stop();
    if (spinner) {
      process.stderr.write(ANSI_RESET);
      process.stdout.write(ANSI_RESET);
    }

    if (parsed.json) {
      process.stdout.write(`${JSON.stringify({ results }, null, 2)}\n`);

      return 0;
    }

    process.stdout.write(`${renderTable(results, plan.groups)}\n`);
    return 0;
  } catch (error) {
    if (spinnerTimer) clearTimeout(spinnerTimer);
    spinner?.stop();
    if (spinner) {
      process.stderr.write(ANSI_RESET);
      process.stdout.write(ANSI_RESET);
    }

    const message = error instanceof Error ? error.message : "Unknown error";
    const hasSuggested = plan?.groups.some((group) => group.suggested) ?? false;
    if (hasSuggested && parsed) {
      const fallback = parsed.domains.join("\n");
      process.stdout.write(`${fallback}\n`);
    }
    if (message === "Missing Dynadot key." || message === "Invalid Dynadot key.") {
      process.stderr.write(
        `${redText(message === "Missing Dynadot key." ? "Missing Dynadot key. Export DYNADOT_API_PRODUCTION_KEY or pass --dynadot-key." : "Invalid Dynadot key. Get your production key here:")}\n`,
      );
      process.stderr.write(`${redText(DYNADOT_KEY_URL)}\n`);
    } else {
      process.stderr.write(`${message}\n`);
    }

    return 1;
  }
};

const run = async (): Promise<void> => {
  const code = await runCli(process.argv.slice(2));
  process.exit(code);
};

if (import.meta.main) {
  void run();
}
