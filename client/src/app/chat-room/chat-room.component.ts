import {AfterViewInit, Component, OnInit, ViewChild, ViewContainerRef} from '@angular/core';
import {ActivatedRoute} from "@angular/router";
import {WsService} from "../ws.service";
import {HttpServiceService} from "../http-service.service";
import {MESSAGE_TYPE, messageData, MessagesResponse, payload} from "../types";
import {MessageService} from "../message.service";
// @ts-ignore

@Component({
  selector: 'app-chat-room',
  templateUrl: './chat-room.component.html',
  styleUrls: ['./chat-room.component.css'],
  providers: [
    WsService
  ]
})
export class ChatRoomComponent implements OnInit, AfterViewInit {

  // @ts-ignore
  id: string
  // @ts-ignore
  username: string
  // @ts-ignore
  connectionId: number

  // @ts-ignore
  @ViewChild('messageList', {read: ViewContainerRef}) messageList: ViewContainerRef

  constructor(private route: ActivatedRoute, private ws: WsService, private httpServiceService: HttpServiceService, private messageService:MessageService) {
    this.username = window.localStorage.getItem('username') || ''
  }

  ngOnInit(): void {
    this.route.params.subscribe(params => {
      console.log(`query params ${JSON.stringify(params)}`)
      this.id = params.id
    })
  }

  send(message: string) {
    const data: messageData = {
      connectionId: this.connectionId,
      owner: this.username,
      content: message,
      createdTime: Date.now()
    }
    this.ws.sendMessage({
      username: this.username,
      connectionId: this.connectionId || -1,
      messageType: MESSAGE_TYPE.SendMessage,
      data: data
    })
    // add message to view child
    this.addMessage(this.username, message, data.createdTime || 0)
  }

  onKeyDown(event: KeyboardEvent, content: HTMLTextAreaElement) {
    if (event.key === "Enter") {
      if (!event.shiftKey) {
        this.sendMessage(content)
        content.value = ''
      }
    }
  }

  addMessage(username: string, message: string, createdTime: number) {
    this.messageService.addComponent(this.username, message, createdTime || 0)
  }

  loadMessages() {
    this.messageService.addDynamicComponent(this.id)
  }

  sendMessage(content: HTMLTextAreaElement) {
    this.send(content.value)
    content.value = ''
  }

  ngAfterViewInit() {
    this.ws.connect(this.id)
    this.ws.conn.onmessage = (evt:MessageEvent) => {
      const messages = evt.data.split('\n');
      for (const msgStr of messages) {
        const message: payload = JSON.parse(msgStr)
        switch (message.messageType) {
          case MESSAGE_TYPE.NewConnection: {
            this.connectionId = message.connectionId || -1
          } break
          case MESSAGE_TYPE.ErrorMessage: {
            console.error(message.data.content)
          } break
          case MESSAGE_TYPE.ReceiveMessage: {
            console.log(message.data.content)
            this.addMessage(message.data.owner, message.data.content, message.data.createdTime || 0)
          } break
          default: console.error(`invalid message - ${JSON.stringify(message)}`)
        }
      }
    }
    this.messageService.setRootViewContainerRef(this.messageList)
    this.loadMessages()
  }

}
