import { Injectable } from '@angular/core';
import {environment} from "../environments/environment";
// @ts-ignore
import {payload} from "./types";

@Injectable({
  providedIn: 'root'
})
export class WsService {

  // @ts-ignore
  conn: WebSocket
  constructor() { }

  public connect(room: string) {
    this.conn = new WebSocket(`ws://${environment.apiUrl}/ws/${room}`)
    this.conn.onclose = (evt:CloseEvent) => {
      console.log('connection closed')
    }
  }

  public sendMessage(msg: payload) {
    this.conn.send(JSON.stringify(msg))
  }
}
