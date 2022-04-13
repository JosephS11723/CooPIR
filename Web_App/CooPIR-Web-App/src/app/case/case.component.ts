import { Component, OnInit } from '@angular/core';
import { GlobalConstants } from '../common/global-constraints';
import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { CookieService } from 'ngx-cookie-service';

@Component({
  selector: 'app-case',
  templateUrl: './case.component.html',
  styleUrls: ['./case.component.css']
})
export class CaseComponent implements OnInit {
  file:any;
  fileName = '';
  doc = ""
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
    {
      label: 'Dashboard',
      icon: 'assignment',
      route: '/dashboard'
    }
  ];

  fileList = [
    {
      name: '',
      created: '',
      last_modified: '',
      route: ''
    }
  ];
  constructor(private http: HttpClient, private cookieService:CookieService) { }


  getFiles(): void
  {
    //console.log("Getting files");
    const params = new HttpParams()
    .append('uuid', this.cookieService.get("currentUUID"));
    //.append('uuid', GlobalConstants.currentCase);

    console.log("Getting files for: ", GlobalConstants.currentCase);
    //get all the files in the selected case
    this.http.get("http://localhost:8080/api/v1/case/files", {params: params, observe: 'response'})
    .subscribe( response => {
      console.log("Logging response");
      console.log(response.body);
      let retrievedFiles: any
        if(response.body != null)
        {
          retrievedFiles = response.body;
          //push each file in the case into the files to be displayed
          for(var index = 0; index < retrievedFiles.files.length; index++)
          {
           // this.http.get("http://localhost:8080/api/v1/file/info", {params: fileParams, observe: 'response'})
           // .subscribe( response => {
           //   console.log("Here is the file info: ", response);
           // });
            //console.log(retrievedCases.cases[index]);
            this.fileList.push({
              name: retrievedFiles.files[index],
              created: 'at some point',
              last_modified: 'sometime',
              route: '/case'
            });
          }
        }
    });

  }

  getFileInfo(uuid: any): void
  {
    var testFile: any;
    this.http.get("http://localhost:8080/api/v1/file/" + GlobalConstants.currentCase + "/" + uuid, {observe: 'response'})
    .subscribe( response => {
      console.log("Response to file down load: ", response);
    });

    //this.doc = "https://www.google.com/"
    var fileParams = new HttpParams()
    .append('caseUUID', GlobalConstants.currentCase)
    .append('fileUUID', uuid);
    this.http.get("http://localhost:8080/api/v1/file/info", {params: fileParams, observe: 'response'})
    .subscribe( response => {
        console.log("Here is the file info: ", response);
      });
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
  ngOnInit(): void 
  {
    console.log("Cookie test: ", this.cookieService.get("currentUUID"));
    console.log("Inside case: ", GlobalConstants.currentCase);
    this.getFiles();
  }
  
  submitFile(event:any)
  {
    if (this.file)
    {
      console.log("File ready to send");
      
      this.fileName = this.file.name;
      var caseuuid = GlobalConstants.currentCase;

      const params = new HttpParams()
      .append('caseuuid', caseuuid)
      .append('filename', this.fileName);
      
      const formData = new FormData();
      formData.append("file", this.file);

      this.http.post("http://localhost:8080/api/v1/file", formData, 
      {
        params: params,
        observe: 'response'})
        .subscribe(response => {
            console.log("logging respose");
            console.log(response);
          }, error => {
            console.log("logging error");
            console.log(error);
          });
      this.fileList.push(
        {
        name: this.fileName,
        created: '',
        last_modified: '',
        route: ''
        });
    }
  }

  emptyClick(): void
  {
    console.log("Button function");
  }
}
