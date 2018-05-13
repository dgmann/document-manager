import {Moment} from "moment";

export interface Record {
  id: string;
  date: Moment;
  receivedAt: Moment;
  comment: string;
  sender: string;
  pages: Page[];
  categoryId: string;
  tags: string[];
  patientId: string;
  escalated: boolean;
  processed: boolean;
  status: Status;
}

export class Page {
  id: string;
  url: string;
  content: string;
}

export class PageUpdate {
  public rotate: number = 0;

  constructor(public id: string, public url: string) {
  }

  public static FromPage(page: Page) {
    return new PageUpdate(page.id, page.url)
  }
}


export enum Status {
  INBOX = "inbox",
  REVIEW = "review",
  ESCALATED = "escalated",
  OTHER = "other",
  DONE = "done",
  NONE = ""
}
