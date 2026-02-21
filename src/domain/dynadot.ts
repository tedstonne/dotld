import type { Currency, SearchResult } from "@/shared/types";

const DYNADOT_API_URL = "https://api.dynadot.com/api3.json";

type DynadotSearchRow = {
  DomainName?: string;
  Available?: string;
  Price?: string;
  Status?: string;
};

type DynadotSearchResponse = {
  SearchResponse?: {
    ResponseCode?: string;
    Error?: string;
    SearchResults?: DynadotSearchRow[];
  };
  Response?: {
    ResponseCode?: string;
    Error?: string;
  };
};

const timeoutSignal = (timeoutMs: number): AbortSignal => {
  const controller = new AbortController();
  setTimeout(() => controller.abort(), timeoutMs);

  return controller.signal;
};

export const parseDynadotPrice = (value: string | undefined): string | null => {
  if (!value) return null;

  const registrationMatch = value.match(/Registration Price:\s*([0-9]+(?:\.[0-9]+)?)/i);
  if (registrationMatch?.[1]) return registrationMatch[1];

  const genericMatch = value.match(/([0-9]+(?:\.[0-9]+)?)\s+in\s+USD/i);

  return genericMatch?.[1] ?? null;
};

export const affiliateUrl = (domain: string, template?: string): string => {
  const fallback = `https://www.dynadot.com/domain/search?domain=${encodeURIComponent(domain)}`;
  if (!template || template.trim().length === 0) return fallback;

  if (template.includes("{domain}")) {
    return template.replaceAll("{domain}", encodeURIComponent(domain));
  }

  try {
    const url = new URL(template);
    url.searchParams.set("domain", domain);

    return url.toString();
  } catch {
    return fallback;
  }
};

const requestOne = async (params: {
  apiKey: string;
  domain: string;
  currency: Currency;
  timeoutMs: number;
}): Promise<DynadotSearchResponse> => {
  const query = new URLSearchParams({
    key: params.apiKey,
    command: "search",
    show_price: "1",
    currency: params.currency,
    domain0: params.domain,
  });

  try {
    const response = await fetch(`${DYNADOT_API_URL}?${query.toString()}`, {
      method: "GET",
      signal: timeoutSignal(params.timeoutMs),
    });

    return (await response.json()) as DynadotSearchResponse;
  } catch {
    throw new Error("Dynadot request timed out");
  }
};

const mapResult = (params: {
  domain: string;
  payload: DynadotSearchResponse;
  currency: Currency;
  affiliateTemplate: string | undefined;
}): SearchResult => {
  const rows = params.payload.SearchResponse?.SearchResults ?? [];
  const row = rows.find((item) => item.DomainName?.toLowerCase() === params.domain) ?? rows[0];
  const available = row?.Available?.toLowerCase() === "yes";

  return {
    domain: params.domain,
    available,
    price: available ? parseDynadotPrice(row?.Price) : null,
    currency: params.currency,
    buyUrl: available ? affiliateUrl(params.domain, params.affiliateTemplate) : null,
    source: "dynadot",
    cached: false,
    quotedAt: new Date().toISOString(),
    ...(row?.Status && row.Status !== "success" ? { error: row.Status } : {}),
  };
};

const ensureSuccess = (payload: DynadotSearchResponse): void => {
  const authCode = payload.Response?.ResponseCode;
  if (authCode === "-1") {
    const message = payload.Response?.Error ?? "Dynadot authentication failed";
    if (message.toLowerCase().includes("invalid key")) {
      throw new Error("Invalid Dynadot key.");
    }

    throw new Error(message);
  }

  const searchCode = payload.SearchResponse?.ResponseCode;
  if (searchCode && searchCode !== "0") {
    const message = payload.SearchResponse?.Error ?? "Dynadot search failed";
    if (message.toLowerCase().includes("invalid key")) {
      throw new Error("Invalid Dynadot key.");
    }

    throw new Error(message);
  }
};

export const searchDynadot = async (params: {
  apiKey: string;
  domains: string[];
  currency: Currency;
  timeoutMs: number;
  affiliateTemplate: string | undefined;
}): Promise<SearchResult[]> => {
  const results: SearchResult[] = [];

  for (const domain of params.domains) {
    const payload = await requestOne({
      apiKey: params.apiKey,
      domain,
      currency: params.currency,
      timeoutMs: params.timeoutMs,
    });
    ensureSuccess(payload);
    results.push(
      mapResult({
        domain,
        payload,
        currency: params.currency,
        affiliateTemplate: params.affiliateTemplate,
      }),
    );
  }

  return results;
};
