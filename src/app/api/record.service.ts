import { ManyQueryResult, NgrxJsonApiService, StoreResource } from "ngrx-json-api";
import { Injectable } from "@angular/core";

@Injectable()
export class RecordService {
  constructor(private ngrxService: NgrxJsonApiService) {}

  public find(id: string) {
    return this.ngrxService.findOne({query: {
        type: 'Record',
        id: id,
      }});
  }

  public all() {
    return this.ngrxService.findMany({query: {type: 'Record'}})
      .map<ManyQueryResult,Record[]>(r => r.data && r.data.map(data => this.toRecord(data)) || undefined)
  }


  private toRecord(data: StoreResource) {
    return new Record(data.id, new Date(data.attributes.date), data.attributes.comment, data.attributes.sender);
  }
}

export class Record {
  constructor(public id: string,
              public date: Date,
              public comment: string,
              public sender: string) {}
}
