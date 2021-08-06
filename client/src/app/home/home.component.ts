import { Component, OnInit } from '@angular/core';
import * as word from 'random-words'
import {HttpServiceService} from "../http-service.service";

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css'],
  providers: [HttpServiceService]
})
export class HomeComponent implements OnInit {

  defaultValue: string
  username: string

  constructor(private HttpService:HttpServiceService) {
    this.defaultValue = ''
    this.username = window.localStorage.getItem('username') || `${word(1)}_${Date.now()}`
  }

  ngOnInit(): void {
  }

  createRoom() {
    this.HttpService.createRoom(this.username).subscribe(
      ({ id }) => {
        this.defaultValue = `${window.location.protocol}//${window.location.host}/room/${id}`
      }
    )
  }

  joinRoom(url:string) {
    window.open(url, '_blank')
  }

}
