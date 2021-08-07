import { Component } from '@angular/core';
import * as word from "random-words";

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  title = 'client';
  constructor() {
    if (!window.localStorage.getItem('username')) {
      window.localStorage.setItem('username', `${word(1)}_${Date.now()}`)
    }
  }
}
