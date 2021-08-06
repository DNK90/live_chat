import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HomeComponent } from './home/home.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { ChatRoomComponent } from './chat-room/chat-room.component';
import { ChatMessageComponent } from './chat-message/chat-message.component';
import { HttpClientModule } from '@angular/common/http';
import { MaterialModule } from "./material/material.module";
import { FormsModule } from '@angular/forms';
import { HttpErrorHandler } from "./http-error-handler.service";
import { HttpServiceService } from "./http-service.service";
import { WsService } from "./ws.service";
import { MessageService } from "./message.service";

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    ChatRoomComponent,
    ChatMessageComponent
  ],
  imports: [
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    HttpClientModule,
    MaterialModule,
    FormsModule,
  ],
  providers: [
    HttpErrorHandler,
    HttpServiceService,
    WsService,
    MessageService
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
