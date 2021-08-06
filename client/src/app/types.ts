export interface newRoomResponse {
  id?:string
}

export interface messageData {
  messageId?: number
  owner: string
  connectionId ?: number
  content: string
  status?: number
  createdTime?: number
}

export interface payload {
  username?: string
  connectionId?: number
  messageType: number
  data: messageData
}

export const MESSAGE_TYPE = {
  SendMessage: 0,
  ReceiveMessage:1,
  DeleteMessage:2,
  NewConnection:3,
  ErrorMessage:4,
}

export interface MessageModel {
  id: number,
  content: string,
  username: string,
  room: string,
  created_time: number,
  status: number,
}

export interface MessagesResponse {
  data: MessageModel[]
}

