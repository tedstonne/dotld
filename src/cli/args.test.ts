import { describe, expect, test } from "bun:test";

import { parseArgs } from "@/cli/args";

describe("cli argument parser", () => {
  test("parses search domains and flags", () => {
    const parsed = parseArgs([
      "search",
      "example.com",
      "murk.ink",
      "--json",
      "--dynadot-key",
      "key-123",
      "--timeout",
      "5s",
    ]);

    expect(parsed.json).toBeTrue();
    expect(parsed.dynadotKey).toBe("key-123");
    expect(parsed.timeoutMs).toBe(5000);
    expect(parsed.domains).toEqual(["example.com", "murk.ink"]);
  });
});
