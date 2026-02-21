export const SUPPORTED_CURRENCIES = ["USD"] as const;

export type Currency = (typeof SUPPORTED_CURRENCIES)[number];

export type SearchResult = {
  domain: string;
  available: boolean;
  price: string | null;
  currency: Currency;
  buyUrl: string | null;
  source: "dynadot";
  cached: boolean;
  quotedAt: string;
  error?: string;
};
