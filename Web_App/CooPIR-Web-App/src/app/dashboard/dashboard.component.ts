import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { GlobalConstants } from '../common/global-constraints';
import { CookieService } from 'ngx-cookie-service';
import { analyzeAndValidateNgModules } from '@angular/compiler';
import { Router } from '@angular/router';
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
    
  ];

  //cases the user has access to
  caseList = [
    {
      name: '',
      uuid: '',
      supervisor: '',
      date_created: ''
    }
  ];
  
  constructor(private http: HttpClient, private cookieService:CookieService, private router: Router) { }




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
            var fileParams = new HttpParams()
            .append('uuid', retrievedCases.cases[index]);
            this.http.get("http://localhost:8080/api/v1/case", {params: fileParams, observe: 'response'})
            .subscribe(response => {
              var caseData: any;
              caseData = response.body;
              //console.log("Case metadata: ", response);
              //console.log("Case name: ", caseData.case.name);
              //console.log("Case uuid: ", caseData.case.uuid);

              this.caseList.push({
                name: caseData.case.name,
                uuid: caseData.case.uuid,
                supervisor: 'Joseph',
                date_created: caseData.case.dateCreated
              });
            })
            //console.log(retrievedCases.cases[index]);
            
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
    this.router.navigateByUrl('/case', { replaceUrl: true});
    GlobalConstants.currentCase = uuid;

  }

  //cases are retrieved as soon as someone logs in
  ngOnInit(): void 
  {
    //console.log("Running oninit")
    //console.log("Can the user make a case (before)", GlobalConstants.canUserMakeCase)
    let caseMaker: any;
    this.http.get("http://localhost:8080/api/v1/case/make", {observe: 'response'})
    .subscribe( response =>
      { 
        caseMaker = response.body
        GlobalConstants.canUserMakeCase = caseMaker.allow;
        //console.log("Can the user make a case (after):", GlobalConstants.canUserMakeCase);

      });
    if(GlobalConstants.canUserMakeCase === true)
    {
      this.menuItems.push(
        {
          label: "Make Case",
          icon: 'library_add',
          route: '/makeCase'
        }
      )
    }
    this.getCases();
  }

}
