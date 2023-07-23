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
  Exact = "exact",
  Regex = "regex",
}
