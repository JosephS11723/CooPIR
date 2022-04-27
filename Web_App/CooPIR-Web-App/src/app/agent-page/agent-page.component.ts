import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-agent-page',
  templateUrl: './agent-page.component.html',
  styleUrls: ['./agent-page.component.css']
})
export class AgentPageComponent implements OnInit {

  menuItems = [
    {
      label: 'Logout',
      icon: 'exit_to_app',
      route: '/login'
    },    
    {
      label: 'Dashboard',
      icon: 'assignment',
      route: '/dashboard'
    }
  ];


  constructor() { }

  launchAgent(): void 
  {

  }
  
  ngOnInit(): void 
  {
  
  }

}
