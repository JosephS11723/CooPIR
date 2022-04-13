import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { GlobalConstants } from '../common/global-constraints';
import { CookieService } from 'ngx-cookie-service';
//import { getPackedSettings } from 'http2';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})

export class DashboardComponent implements OnInit {
  //items placed in the toolbar
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

  //cases the user has access to
  caseList = [
    {
      name: '',
      supervisor: '',
      last_modified: '',
      route: '',
      function: this.emptyClick()
    }
  ];
  
  constructor(private http: HttpClient, private cookieService:CookieService) { }

  //gets the cases that the user has access to
  getCases(): void
  {
    //get the case list from the db
    this.http.get("http://localhost:8080/api/v1/cases", { observe: 'response'})
      .subscribe(response => {
        console.log("Logging response");
        console.log(response.body);
        //store the response as an any type so we can access the data inside
        let retrievedCases: any
        if(response.body != null)
        {
          retrievedCases = response.body;
          //push each case the user has access to into the cases to be displayed
          for(var index = 0; index < retrievedCases.cases.length; index++)
          {
            //console.log(retrievedCases.cases[index]);
            this.caseList.push({
              name: retrievedCases.cases[index],
              supervisor: 'Joseph',
              last_modified: 'sometime',
              route: '/case',
              function: this.goToCase(retrievedCases.cases[index])
            });
          }
        }
      }, error => {
        console.log("logging error");
        console.log(error);
      });
      
  }

  emptyClick(): void
  {

  }

  goToCase(uuid: string): void 
  {
    console.log("Going to: ", uuid);
    this.cookieService.set("currentUUID", GlobalConstants.currentCase);
    GlobalConstants.currentCase = uuid;

  }

  //cases are retrieved as soon as someone logs in
  ngOnInit(): void 
  {
    //console.log("Running oninit")
    this.getCases();
  }

}
