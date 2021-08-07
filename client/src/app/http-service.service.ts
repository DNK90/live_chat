import { Injectable } from '@angular/core';
import {HttpClient, HttpErrorResponse, HttpHeaders} from "@angular/common/http";
import { Observable, throwError } from 'rxjs';
import { catchError, retry } from 'rxjs/operators';
// @ts-ignore
import { environment } from "../environments/environment";
import { HttpErrorHandler } from "./http-error-handler.service";
import {MessagesResponse, newRoomResponse} from "./types";

const httpOptions = {
  headers: new HttpHeaders({
    'Content-Type':  'application/json',
  })
};

@Injectable({
  providedIn: 'root'
})
export class HttpServiceService {

  constructor(private httpClient: HttpClient, private httpErrorHandler: HttpErrorHandler) { }

  public createRoom(username: string): Observable<newRoomResponse> {
    return this.httpClient.post<newRoomResponse>(`http://${environment.apiUrl}/v1/room`, {username}, httpOptions).pipe(
      catchError(
        this.httpErrorHandler.handleError<newRoomResponse>('createRoom', 'id')
      )
    )
  }

  public loadMessages(room: string):Observable<MessagesResponse> {
    console.log(`[httpService][loadMessages] roomId=${room}`)
    return this.httpClient.get<MessagesResponse>(`http://${environment.apiUrl}/v1/room/${room}/messages`, httpOptions).pipe(
      catchError(
        this.httpErrorHandler.handleError<MessagesResponse>('loadMessages', '')
      )
    )
  }
}
