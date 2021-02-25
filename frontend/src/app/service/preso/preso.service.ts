import { Injectable } from '@angular/core';
import { Observable, of } from 'rxjs';
import { MessageService } from '../message/message.service';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { catchError, map, tap } from 'rxjs/operators';
import { environment } from '../../../environments/environment';
import { Content, Nugget} from '../content/content.service'

export class Preso {
  public title: string
  public link: string
  public description: string
  public published:Date
}

@Injectable({
  providedIn: 'root'
})
export class PresoService {

  constructor(
    private http: HttpClient,
    private messageService: MessageService
  ) { }

  private presoUrl: string = environment.preso_url;
  httpOptions = {
    headers: new HttpHeaders({ 'Content-Type': 'application/atom+xml' })
  };

  getPresos (): Observable<Content> {
    return this.http.get<Content>(this.presoUrl)
      .pipe(
        tap(_ => this.log('fetched repos'))
        // catchError(this.handleError<Repo[]>('getRepos', []))
      );
  }


  private handleError (error: any) {
    // In a real world app, we might use a remote logging infrastructure
    // We'd also dig deeper into the error to get a better message
    let errMsg = (error.message) ? error.message :
      error.status ? `${error.status} - ${error.statusText}` : 'Server error';
    console.error(errMsg); // log to console instead
    return Observable.throw(errMsg);
  }

  private log(message: string) {
    this.messageService.add(`PresoService: ${message}`);
  }


}

