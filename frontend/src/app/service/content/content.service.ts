import { Injectable } from '@angular/core';

export class Nugget {
  public title: string
  public link: string
  public description: string
  public timestamp: Date
}

export class Content {
  public nuggets: Nugget[]
  public cached:boolean
}

@Injectable({
  providedIn: 'root'
})
export class ContentService {

  constructor() { }
}

