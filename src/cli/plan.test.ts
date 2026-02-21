import { describe, expect, test } from "bun:test";

import { MAINSTREAM_TLDS, createQueryPlan } from "@/cli/plan";

describe("query plan", () => {
  test("keeps explicit domain exact", () => {
    const plan = createQueryPlan(["murk.ink"]);

    expect(plan.groups[0]?.suggested).toBeFalse();
    expect(plan.groups[0]?.domains).toEqual(["murk.ink"]);
    expect(plan.lookupDomains).toEqual(["murk.ink"]);
  });

  test("expands bare domain to mainstream tld list", () => {
    const plan = createQueryPlan(["murk"]);
    const expected = MAINSTREAM_TLDS.map((tld) => `murk.${tld}`);

    expect(plan.groups[0]?.suggested).toBeTrue();
    expect(plan.groups[0]?.domains).toEqual(expected);
    expect(plan.lookupDomains).toEqual(expected);
  });
});
