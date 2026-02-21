import { describe, expect, test } from "bun:test";

import { createQueryPlan } from "@/cli/plan";
import { renderTable } from "@/cli/table";
import type { SearchResult } from "@/shared/types";

describe("table rendering", () => {
  test("renders compact availability output", () => {
    const results: SearchResult[] = [
      {
        domain: "murk.ink",
        available: true,
        price: "2.00",
        currency: "USD",
        buyUrl: "https://example.com/?d=murk.ink",
        source: "dynadot",
        cached: false,
        quotedAt: "2026-02-20T00:00:00.000Z",
      },
    ];

    const plan = createQueryPlan(["murk.ink"]);

    const table = renderTable(results, plan.groups);

    expect(table).toContain(" · ");
    expect(table).toContain("murk.ink");
    expect(table).toContain("2.00");
  });

  test("renders suggestion tree with connectors", () => {
    const results: SearchResult[] = [
      {
        domain: "murk.com",
        available: false,
        price: null,
        currency: "USD",
        buyUrl: null,
        source: "dynadot",
        cached: false,
        quotedAt: "2026-02-20T00:00:00.000Z",
      },
      {
        domain: "murk.sh",
        available: true,
        price: "39.99",
        currency: "USD",
        buyUrl: "https://example.com/?d=murk.sh",
        source: "dynadot",
        cached: false,
        quotedAt: "2026-02-20T00:00:00.000Z",
      },
    ];

    const plan = createQueryPlan(["murk"]);
    const output = renderTable(results, plan.groups);

    expect(output).toContain("murk");
    expect(output).toContain("├─");
    expect(output).toContain("└─");
    expect(output).toContain("Taken");
  });
});
