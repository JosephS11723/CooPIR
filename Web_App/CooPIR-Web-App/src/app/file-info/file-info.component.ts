import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { CookieService } from 'ngx-cookie-service';
import * as FileSaver from 'file-saver';
//import { info } from 'console';

@Component({
  selector: 'app-file-info',
  templateUrl: './file-info.component.html',
  styleUrls: ['./file-info.component.css']
})
export class FileInfoComponent implements OnInit {
  menuItems = [
    {
      label: 'Logout',
      icon: 'exit_to_app',
      route: '/login'
    },
    {
      label: 'Back to Case',
      icon: 'keyboard_tab',
      route: '/case'
    }

  ];
  infoList = new Array<any>();
  constructor(private http: HttpClient, private cookieService:CookieService,  private router: Router) { }

  ngOnInit(): void 
  {
    var fileParams = new HttpParams()
    .append('caseUUID', this.cookieService.get("currentUUID"))
    .append('fileUUID', this.cookieService.get("fileUUID"));

    var fileInfo:any;
    var infoKeys:any;
    this.http.get("http://localhost:8080/api/v1/file/info", {params: fileParams, observe: 'response'})
    .subscribe( response => {
      console.log("File info request: ", response.body)
      fileInfo = response.body
      console.log("fileInfo: ", fileInfo.file.uuid)
      infoKeys = Object.keys(fileInfo.file).forEach(e => 
      { 
        if(e == 'filename')
        {
          this.infoList.push(
            {
              key: e,
              info: fileInfo.file[e].split("/").pop()
            }
          )
        }
        else if(e == 'relations')
        {
          if(fileInfo.file[e] == ":contains")
          {
            this.infoList.push(
              {
                key: e,
                info: ''
              }
            )
          }
          else
          {
            this.infoList.push(
              {
                key: e,
                info: fileInfo.file[e]}
              )
          }
        }
        else
        {
          this.infoList.push(
            {
              key: e,
              info: fileInfo.file[e]}
            )
        }
      })
      //console.log("Keys: ", infoKeys)
      console.log(this.infoList)
    })
  }


  passTheSalt(fileuuid:any)
  {
    
    //this.cookieService.set('fileUUID', fileuuid)
    console.log("link clicked")
    console.log("Passed uuid: ", fileuuid)

    var finaluuid = fileuuid[0].split(":")
    console.log("UUID parse: ", finaluuid[0])
    this.cookieService.set('fileUUID', finaluuid[0])
    window.location.reload()
  }

  downloadFile(name: any): void
  {
    //this.cookieService.set('fileUUID', uuid);
    this.http.get("http://localhost:8080/api/v1/file/" + this.cookieService.get('fileUUID')  + "/" + this.cookieService.get("currentUUID"), {observe: 'response'})
    .subscribe(response =>
      {
        if(response.body === null)
        {
          console.log("Response body is null");
        }
        console.log("response from download call: ", response);
        //console.log("New download test");
        let testData:any;
        testData = response.body;
        //console.log("Blob body: ", testData);
        let blob = new Blob([testData], {type: 'text/plain;charset=utf-8'});           
      });

      var downloadFile = this.http.get("http://localhost:8080/api/v1/file/"   + this.cookieService.get("currentUUID") + "/"+ this.cookieService.get("fileUUID"), {observe: 'response', responseType: 'blob'});
      console.log("DownloadFile: ", downloadFile);
      //var blurg = new BlobPart;
      downloadFile.subscribe(subscriber =>
        {
          console.log("Subscriber: ", subscriber);
          if(subscriber.body != null)
          {
            //console.log("Subscriber body: ", subscriber.body);
            //console.log("Subscriber header: ", subscriber.headers)
            const blob = new Blob([subscriber.body], {type: 'application/octetstream'});
            console.log("Blob test: ", blob);
            FileSaver.saveAs(blob, this.cookieService.get('fileName'));
          }
          else
          {
            console.log("Subscriber body is null");
          }     
        });
         
  }
}
