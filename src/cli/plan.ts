export const MAINSTREAM_TLDS = ["com", "net", "org", "io", "ai", "co", "app", "dev", "sh"] as const;

export type QueryGroup = {
  input: string;
  root: string;
  domains: string[];
  suggested: boolean;
};

const bareLabelPattern = /^[a-z0-9](?:[a-z0-9-]{0,61}[a-z0-9])?$/;

const hasTld = (value: string): boolean => value.includes(".");

const bareLabel = (value: string): boolean => bareLabelPattern.test(value);

export const createQueryPlan = (
  inputs: string[],
): { groups: QueryGroup[]; lookupDomains: string[] } => {
  const groups = inputs.map((input) => {
    if (!hasTld(input) && bareLabel(input)) {
      return {
        input,
        root: input,
        domains: MAINSTREAM_TLDS.map((tld) => `${input}.${tld}`),
        suggested: true,
      } satisfies QueryGroup;
    }

    return {
      input,
      root: input,
      domains: [input],
      suggested: false,
    } satisfies QueryGroup;
  });

  const lookupDomains = [...new Set(groups.flatMap((group) => group.domains))];

  return { groups, lookupDomains };
};
