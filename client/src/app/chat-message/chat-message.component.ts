import { Component, OnInit } from '@angular/core';
import {messageData} from "../types";

@Component({
  selector: 'app-chat-message',
  templateUrl: './chat-message.component.html',
  styleUrls: ['./chat-message.component.css']
})
export class ChatMessageComponent implements OnInit {
  // @ts-ignore
  owner: string
  username: string
  // @ts-ignore
  createdTime: number
  // @ts-ignore
  isYou: string

  // @ts-ignore
  content: string

  constructor() {
    this.username = window.localStorage.getItem('username') || ''
  }

  ngOnInit(): void {
  }

  public formatDate(): string {
    const date = new Date(`${this.createdTime}`.length === 10 ? this.createdTime*1000 : this.createdTime)
    return `${date.getDay()+1}/${date.getMonth()+1}/${date.getFullYear()} ${date.getHours()}:${date.getMinutes()}:${date.getSeconds()}`
  }

  public load(msg: messageData) {
    console.log(`[ChatMessageComponent][load] msg=${JSON.stringify(msg)}`)
    this.owner = msg.owner
    this.createdTime = msg.createdTime || 0
    this.content = msg.content
    this.isYou = this.owner === this.username ? '(You)' : ''
  }

}
