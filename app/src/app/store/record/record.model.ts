export interface Record {
  id: string;
  date: Date;
  receivedAt: Date;
  comment: string;
  sender: string;
  pages: Page[];
  category: string;
  tags: string[];
  patientId: string;
  escalated: boolean;
  processed: boolean;
  requiredAction: RequiredAction;
}

export class Page {
  url: string;
  content: string;
}

export enum RequiredAction {
  REVIEW = "review",
  ESCALATED = "escalated",
  OTHER = "other",
}