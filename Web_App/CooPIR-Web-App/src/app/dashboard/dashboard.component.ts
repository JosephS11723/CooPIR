import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {
  menuItems = [
    {
      label: 'Logout',
      icon: 'login',
      route: '/login'
    },    
    {
      label: 'Home',
      icon: 'notes',
      route: '/home'
    },
  ];

  constructor() { }

  ngOnInit(): void {
  }

}
