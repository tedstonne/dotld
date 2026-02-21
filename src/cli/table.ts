import { stripVTControlCharacters } from "node:util";

import pc from "picocolors";

import type { QueryGroup } from "@/cli/plan";
import type { SearchResult } from "@/shared/types";

const divider = pc.dim(" · ");
const useAnsiColor = process.env.NO_COLOR === undefined;

const visibleLength = (value: string): number => stripVTControlCharacters(value).length;

const padVisible = (value: string, width: number): string => {
  const padding = Math.max(0, width - visibleLength(value));

  return `${value}${" ".repeat(padding)}`;
};

const noValue = pc.dim("N/A");

const redAlert = (value: string): string =>
  useAnsiColor ? `\x1b[1;38;2;255;95;95m${value}\x1b[0m` : value;

const availableLine = (result: SearchResult): string => {
  const price = result.price ? pc.green(pc.bold(`$${result.price}`)) : noValue;
  const buy = result.buyUrl ? pc.cyan(result.buyUrl) : noValue;

  return `${pc.bold(result.domain)}${divider}${price}${divider}${buy}`;
};

const unavailableLine = (domain: string): string =>
  `${pc.bold(domain)}${divider}${redAlert("Taken")}`;

const renderSuggestedGroup = (group: QueryGroup, byDomain: Map<string, SearchResult>): string => {
  const lines = [pc.bold(group.root)];
  const domainWidth = Math.max(...group.domains.map((domain) => visibleLength(domain)));

  const detailValues = group.domains.map((domain) => {
    const result = byDomain.get(domain);
    if (!result) return redAlert("Lookup failed");
    if (!result.available) return redAlert("Taken");

    return result.price ? pc.green(pc.bold(`$${result.price}`)) : noValue;
  });
  const detailWidth = Math.max(...detailValues.map((value) => visibleLength(value)));

  for (let index = 0; index < group.domains.length; index += 1) {
    const domain = group.domains[index];
    if (!domain) continue;

    const connector = index === group.domains.length - 1 ? pc.dim("└─ ") : pc.dim("├─ ");
    const paddedDomain = padVisible(pc.bold(domain), domainWidth);
    const result = byDomain.get(domain);
    if (!result) {
      lines.push(
        `${connector}${paddedDomain}${divider}${padVisible(redAlert("Lookup failed"), detailWidth)}`,
      );
      continue;
    }

    if (!result.available) {
      lines.push(
        `${connector}${paddedDomain}${divider}${padVisible(redAlert("Taken"), detailWidth)}`,
      );
      continue;
    }

    const price = result.price ? pc.green(pc.bold(`$${result.price}`)) : noValue;
    const buy = result.buyUrl ? pc.cyan(result.buyUrl) : noValue;

    lines.push(
      `${connector}${paddedDomain}${divider}${padVisible(price, detailWidth)}${divider}${buy}`,
    );
  }

  return lines.join("\n");
};

const renderExactGroup = (group: QueryGroup, byDomain: Map<string, SearchResult>): string => {
  const domain = group.domains[0];
  if (!domain) return unavailableLine(group.input);

  const result = byDomain.get(domain);
  if (!result) return `${pc.bold(domain)}${divider}${redAlert("Lookup failed")}`;

  return result.available ? availableLine(result) : unavailableLine(domain);
};

export const renderTable = (results: SearchResult[], groups: QueryGroup[]): string => {
  const byDomain = new Map(results.map((result) => [result.domain, result]));

  return groups
    .map((group) =>
      group.suggested ? renderSuggestedGroup(group, byDomain) : renderExactGroup(group, byDomain),
    )
    .join("\n\n");
};
