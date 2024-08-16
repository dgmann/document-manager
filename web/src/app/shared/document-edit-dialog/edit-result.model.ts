export interface EditResult {
  id: string;
  change: {
    patientId: string,
    date: Date,
    tags: string[],
    category: string
  };
}
