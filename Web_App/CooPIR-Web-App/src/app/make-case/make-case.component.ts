import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { CookieService } from 'ngx-cookie-service';

@Component({
  selector: 'app-make-case',
  templateUrl: './make-case.component.html',
  styleUrls: ['./make-case.component.css']
})
export class MakeCaseComponent implements OnInit {
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

  constructor(private http: HttpClient, private cookieService:CookieService) { }

  submitCase(): void 
  {
    var caseName = (<HTMLInputElement>document.getElementById("caseName")).value;
    var viewAccess = (<HTMLInputElement>document.getElementById("viewAccess")).value;
    var editAccess = (<HTMLInputElement>document.getElementById("editAccess")).value;
    var description = (<HTMLInputElement>document.getElementById("description")).value;
    console.log("Description: ", description);

    var newCase = {
      uuid: '',
      name: caseName,
      description: description,
      dateCreated: '',
      viewAccess: viewAccess,
      editAccess: editAccess,
      collabs: []
    }

    console.log("Our brand new case: ", newCase);

    this.http.post("http://localhost:8080/api/v1/case/new", newCase, {observe: 'response'})
    .subscribe(response => {
      console.log("New Case Response: ", response);
    })
  }
  ngOnInit(): void {
  }

}
