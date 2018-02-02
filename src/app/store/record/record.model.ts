export interface Record {
  id: string;
  date: Date;
  comment: string;
  sender: string;
  pages: Page[];
  tags: string[];
  patientId: string;
}

export class Page {
  url: string;
  content: string;
}
