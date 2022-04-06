import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {
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
  ];

  caseList = [
    {
      name: 'case1',
      supervisor: 'Joseph',
      last_modified: 'sometime',
      route: '/case'
    },
    {
      name: 'case2',
      supervisor: 'Joseph',
      last_modified: 'sometime',
      route: '/case'
    },
    {
      name: 'case3',
      supervisor: 'Joseph',
      last_modified: 'sometime',
      route: '/case'
    }
  ];

  
  constructor(private http: HttpClient) { }

  getCases(): void 
  {
    console.log("Running dashboard call")
    //this.http.get("http://localhost:8080/api/v1/gimmecases", { observe: 'response'})
    //  .subscribe(response => {
    //    console.log("Logging response");
    //    console.log(response.body);        
    //  });
  }

  emptyClick(): void
  {

  }

  ngOnInit(): void {
    this.getCases();
  }

}
