import { describe, expect, test } from "bun:test";

import { affiliateUrl, parseDynadotPrice } from "@/domain/dynadot";

describe("dynadot parser", () => {
  test("extracts registration price", () => {
    const price = parseDynadotPrice(
      "Registration Price: 2.00 in USD and Renewal price: 21.62 in USD and Domain is not a Premium Domain",
    );

    expect(price).toBe("2.00");
  });

  test("extracts generic price format", () => {
    const price = parseDynadotPrice("77.00 in USD");

    expect(price).toBe("77.00");
  });

  test("generates url from placeholder", () => {
    const url = affiliateUrl("murk.ink", "https://example.com/buy?d={domain}");

    expect(url).toBe("https://example.com/buy?d=murk.ink");
  });
});
