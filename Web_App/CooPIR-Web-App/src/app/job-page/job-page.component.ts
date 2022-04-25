import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { analyzeAndValidateNgModules } from '@angular/compiler';
import { CookieService } from 'ngx-cookie-service';

@Component({
  selector: 'app-job-page',
  templateUrl: './job-page.component.html',
  styleUrls: ['./job-page.component.css']
})
export class JobPageComponent implements OnInit {
  menuItems = [
    {
      label: 'Logout',
      icon: 'exit_to_app',
      route: '/login'
    },
    {
      label: 'Case',
      icon: 'assignment',
      route: '/case'
    }
  ];
  joblist = new Array<any>();
  caseFiles = new Array<any>();
  constructor(private cookieService:CookieService, private http: HttpClient) { }


  getJobs(): void 
  {
    var jobs:any;
    
    this.http.get("http://localhost:8080/api/v1/jobs/types")
    .subscribe(response => {
      jobs = response
      console.log("Response: ", response);
      console.log("Job response", jobs.jobtypes);
      this.joblist = jobs.jobtypes
      console.log("Here is the joblist: ", this.joblist);
      console.log("Here is the job type: ", this.joblist[0].jobtype);
    });

    //get the files the job could be associated with
    const params = new HttpParams()
    .append('uuid', this.cookieService.get("currentUUID"));
    var retrievedFiles:any;
    this.http.get("http://localhost:8080/api/v1/case/files", {params: params, observe: 'response'})
    .subscribe(response => {
      console.log("Here are the files: ", response);
      if(response.body != null)
      {
        retrievedFiles = response.body;
      }
      console.log("Case files: ", retrievedFiles.files);
      this.caseFiles.push(retrievedFiles.files[0]);

    });
  }

  launchJob(): void
  {
    var jobtype = (<HTMLInputElement>document.getElementById("jobtype")).value;
    var fileforjob = (<HTMLInputElement>document.getElementById("fileforjob")).value;
    var jobName = (<HTMLInputElement>document.getElementById("jobName")).value;

    console.log("Job type: ", jobtype);
    console.log("file for job", fileforjob);
    console.log("Job name: ", jobName);
  }
  ngOnInit(): void
  {
    this.getJobs();
  }

}
