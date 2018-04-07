import {Category} from "../../shared/category-service";

export interface Patient {
  id: string;
  firstName: string;
  lastName: string;
  birthDate: Date;
  tags?: string[];
  categories?: Category[];
}
