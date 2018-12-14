import { Moment } from "moment";

export interface EditResult {
  id: string,
  change: {
    patientId: string,
    date: Moment,
    tags: string[],
    category: string
  }
}
