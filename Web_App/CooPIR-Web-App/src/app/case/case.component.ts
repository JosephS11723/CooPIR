import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-case',
  templateUrl: './case.component.html',
  styleUrls: ['./case.component.css']
})
export class CaseComponent implements OnInit {
  file = '';
  menuItems = [
    {
      label: 'Logout',
      icon: 'exit_to_app',
      route: '/login'
    },    
    {
      label: 'Home',
      icon: 'home',
      route: '/home'
    },
    {
      label: 'Upload',
      icon: 'note_add',
      route: '/home'
    },
  ];
  constructor() { }

  ngOnInit(): void {
  }

  onFileSelected(event:any): void 
  {
    var submit_button = (<HTMLInputElement>document.getElementById("submit_file"));
    this.file = event.target.files[0];
    if (this.file)
    {
      console.log("Received file")
      if (submit_button.disabled === true)
      {
        submit_button.disabled = false;
      }
    }
  }
  
  submitFile(event:any)
  {
    if (this.file)
    {
      console.log("File ready to send");
    }
  }

}
