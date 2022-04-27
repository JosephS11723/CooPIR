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
      for(var index = 0; index < retrievedFiles.files.length; index++)
      {
        var fileParams = new HttpParams()
        .append('caseUUID', this.cookieService.get("currentUUID"))
        .append('fileUUID', retrievedFiles.files[index]);
        let fileInfo: any;
        //get the name for the file so that the user doesn't have to look at uuids
        this.http.get("http://localhost:8080/api/v1/file/info", {params: fileParams, observe: 'response'})
        .subscribe(response => {
          //console.log("Here is the job file info: ", response.body);
          fileInfo = response.body;
          this.caseFiles.push({
            filename: fileInfo.file.filename.split("/").pop(),
            fileuuid: fileInfo.file.uuid
          });
        });
        //this.caseFiles.push(retrievedFiles.files[index]);
      }
      //this.caseFiles.push(retrievedFiles.files[0]);

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

    var params = new HttpParams()
    .append('caseuuid', this.cookieService.get("currentUUID"))
    .append('arguments', '')
    .append('files', fileforjob)
    .append('name', jobName)
    .append('jobtype', jobtype);

    this.http.post("http://localhost:8080/api/v1/jobs/new", '', {params: params, observe: 'response'})
    .subscribe(response =>{
      console.log("Response from job start post: ", response);
    })
  }
  ngOnInit(): void
  {
    this.getJobs();
  }

}
