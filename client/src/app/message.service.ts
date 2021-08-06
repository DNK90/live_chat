import {ComponentFactory, ComponentFactoryResolver, ComponentRef, Injectable, ViewContainerRef} from '@angular/core';
import {ChatMessageComponent} from "./chat-message/chat-message.component";
import {ChatRoomComponent} from "./chat-room/chat-room.component";
import {HttpServiceService} from "./http-service.service";
import {MessagesResponse} from "./types";

@Injectable({
  providedIn: 'root'
})
export class MessageService {
  // @ts-ignore
  rootViewContainer: ViewContainerRef
  constructor(private factoryResolver: ComponentFactoryResolver, private httpService: HttpServiceService) { }

  public setRootViewContainerRef(viewContainerRef: ViewContainerRef): void {
    this.rootViewContainer = viewContainerRef;
  }

  public addComponent(owner: string, content: string, createdTime: number) {
    const factory = this.factoryResolver.resolveComponentFactory(ChatMessageComponent);
    const component: ComponentRef<ChatMessageComponent> = factory.create(this.rootViewContainer.injector);
    component.instance.load({
      content: content,
      owner: owner,
      createdTime: createdTime
    })
    this.rootViewContainer.insert(component.hostView)
  }

  public addDynamicComponent(room: string): void {
    this.httpService.loadMessages(room).subscribe(
      (response: MessagesResponse) => {
        if (response.data.length > 0) {
          for (let i=response.data.length-1; i<=0; i--) {
            this.addComponent(
              response.data[i].username,
              response.data[i].content,
              response.data[i].created_time
            )
          }
        }
      }
    )
  }
}
