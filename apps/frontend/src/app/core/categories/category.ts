export interface Category {
  id: string;
  name: string;
  match: MatchConfig;
}

export interface MatchConfig {
  type: MatchType;
  expression: string;
}

export enum MatchType {
  None= "",
  All = "all",
  Any = "any",
  Exact = "exact",
  Regex = "regex",
}
