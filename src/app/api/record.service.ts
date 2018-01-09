import { ManyQueryResult, NgrxJsonApiService, StoreResource } from "ngrx-json-api";
import { Injectable } from "@angular/core";

@Injectable()
export class RecordService {
  constructor(private ngrxService: NgrxJsonApiService) {}

  public find(id: string) {
    return this.ngrxService.findOne({query: {
        type: 'Record',
        id: id,
      }}).map(r => this.toRecord(r.data));
  }

  public all() {
    return this.ngrxService.findMany({query: {type: 'Record'}})
      .map<ManyQueryResult,Record[]>(r => r.data && r.data.map(data => this.toRecord(data)) || undefined)
  }

  public update(id: string, attributes: any) {
    let resource = {
      type: 'records',
      id: id,
      attributes: attributes
    };

    this.ngrxService.patchResource({resource: resource,
      toRemote: false})
  }

  private toRecord(data: StoreResource) {
    if (!data) {
      return null
    }
    const pages = data.attributes.pages.map(page => new Page( page.url, page.content));
    return new Record(data.id, new Date(data.attributes.date), data.attributes.comment, data.attributes.sender, pages);
  }
}

export class Record {
  constructor(public id: string,
              public date: Date,
              public comment: string,
              public sender: string,
              public pages: Page[]) {}
}

export class Page {
  constructor(public url: string,
              public content: string) {}
}
