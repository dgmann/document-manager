import { Action } from "@ngrx/store";

export interface PayloadAction extends  Action{
  payload?: any;
}
