import { Component, OnInit } from '@angular/core';
import { GlobalConstants } from '../common/global-constraints';
import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { CookieService } from 'ngx-cookie-service';
import { saveAs } from 'file-saver';
import * as FileSaver from 'file-saver';
import { Router } from '@angular/router';
import * as Vis from 'vis';
//import { writeFile } from 'fs';

@Component({
  selector: 'app-case',
  templateUrl: './case.component.html',
  styleUrls: ['./case.component.css']
})
export class CaseComponent implements OnInit {
  file:any;
  fileName = '';
  doc = ""

  //nodes:any;
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

  fileList = [
    {
      name: '',
      uuid: '',
      created: '',
      md5: '',
      relations: '',
      route: ''
    }
  ];
  constructor(private http: HttpClient, private cookieService:CookieService,  private router: Router) { }


  getFiles(): void
  {
    //var nodes = [{id: '', value: 0, label: 'fake'}];
    var nodes = new Array<any>();
    console.log("Empty node list: ", nodes);
    //var edges = [{from: 'df', to: 'df', value: 0}];
    var edges = [{}];
    
    const params = new HttpParams()
    .append('uuid', this.cookieService.get("currentUUID"));

    console.log("Getting files for: ", GlobalConstants.currentCase);
    //get all the files in the selected case
    this.http.get("http://localhost:8080/api/v1/case/files", {params: params, observe: 'response'})
    .subscribe( response => {
      //console.log("Logging response");
      //console.log(response.body);
      let retrievedFiles: any
      //add each file to the list that is displayed in a table
        if(response.body != null)
        {
          retrievedFiles = response.body;
          //for each file, get the file info and push it into the list to be displayed
          for(var index = 0; index < retrievedFiles.files.length; index++)
          {
            //retrieve info for each file
            var fileParams = new HttpParams()
            .append('caseUUID', this.cookieService.get("currentUUID"))
            .append('fileUUID', retrievedFiles.files[index]);
            let fileInfo: any;
            //get the info for the file
            this.http.get("http://localhost:8080/api/v1/file/info", {params: fileParams, observe: 'response'})
            .subscribe( response => {
                console.log("Here is the file info: ", response.body);
                fileInfo = response.body;

                //push file into nodes to be displayed by the map
                //nodes.push({ id: fileInfo.file.filename.split("/").pop(), value: 1, label: fileInfo.file.filename.split("/").pop()});
                nodes.push({ id: fileInfo.file.uuid, value: 1, label: fileInfo.file.filename.split("/").pop()});               
                var relations = fileInfo.file.relations;
                
                if(relations != '')
                {
                  console.log("Pushing edge");
                  console.log("File relation: ", relations[0]);
                  edges.push({from: fileInfo.file.uuid, to: relations[0].split(":")[0], value: 1});
                }
                console.log("Nodes: ", nodes);
                console.log("Edges: ", edges);

                //display the map
                var container = document.getElementById("mynetwork");
                var data = {
                  nodes: nodes,
                  edges: edges,
                };
                var options = {
                  nodes: {
                    shape: "dot",
                    scaling: {
                      customScalingFunction: function(min:any, max:any, total:any, value:any) {
                        return value / total;
                      },
                      min: 0,
                      max: 150,
                    },
                  },
                };
                if(container != null)
                {
                  var network = new Vis.Network(container, data, options);
                }

                //push file and its info to be displayed by the table
                //console.log("This is the selected file's name: ", fileInfo.file.filename.split("/").pop());
                this.fileList.push({
                  name: fileInfo.file.filename.split("/").pop(),
                  uuid: fileInfo.file.uuid,
                  created: fileInfo.file.uploadDate,
                  md5: fileInfo.file.md5,
                  relations: '',
                  route: '/case'
                });
            });
          
          }
         
        }
    });

  }

  getFileInfo(uuid: any, name: any): void
  {
    this.http.get("http://localhost:8080/api/v1/file/" + uuid  + "/" + this.cookieService.get("currentUUID"), {observe: 'response'})
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
        //console.log("Blob body2: ", blob);
        //FileSaver.saveAs(blob, 'downloadtest.txt');
        //console.log("After download test");
        //let link = document.createElement('a');
        //console.log("Link: ", testData.url);
        //link.download = testData.url;
        //let blob = new Blob([testData], {type: 'text/plain'});
        //link.href = URL.createObjectURL(blob);
        //link.click();

        //URL.revokeObjectURL(link.href);
      });

      var downloadFile = this.http.get("http://localhost:8080/api/v1/file/"   + this.cookieService.get("currentUUID") + "/"+ uuid, {observe: 'response', responseType: 'blob'});
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
            FileSaver.saveAs(blob, name);
          }
          else
          {
            console.log("Subscriber body is null");
          }     
        });
      //var blobTest = new Blob(blurg, {type: 'text/plain'});
      //console.log("Blob test: ", blobTest);
      //FileSaver.saveAs(blobTest, 'newDownloadTest.txt');



    var testFile2: any;
    testFile2 = this.http.get("http://localhost:8080/api/v1/file/" + uuid  + "/" + this.cookieService.get("currentUUID"), {responseType: 'blob'});
    //console.log("Here is the testFile2: ", testFile2);

    var testFile: any;
    testFile = this.http.get("http://localhost:8080/api/v1/file/" + uuid  + "/" + this.cookieService.get("currentUUID"), {observe: 'response'})
    .subscribe( response => {
      //console.log("Response to file download: ", response);
    });

   // console.log("Here is the testFile: ", testFile);

    //this.doc = "https://www.google.com/"
    
      
  }

  onFileSelected(event:any): void 
  {
    var submit_button = (<HTMLInputElement>document.getElementById("submit_file"));
    this.file = event.target.files[0];
    if (this.file)
    {
      //console.log("Received file")
      if (submit_button.disabled === true)
      {
        submit_button.disabled = false;
      }
    }
  }

  emptyFunc(): void
  {

  }

  ngOnInit(): void 
  {
    console.log("Cookie test: ", this.cookieService.get("currentUUID"));
    console.log("Inside case: ", GlobalConstants.currentCase);
    setTimeout(this.emptyFunc, 1000);
    this.getFiles();
    
  }
  
  submitFile(event:any)
  {
    if (this.file)
    {
      //console.log("File ready to send");
      
      this.fileName = this.file.name;
      var caseuuid = this.cookieService.get("currentUUID");
      var relations = (<HTMLInputElement>document.getElementById("relations")).value;
      console.log("Relation choice: ", relations);


      const params = new HttpParams()
      .append('caseuuid', caseuuid)
      .append('fileuuid', '/'+this.fileName)
      .append('relations', relations + ":contains")
      
      .append('tags', "testtag");

      if(relations != '')
      {
        console.log("Adding relation", relations);
        params.append('relations', relations + ":contains")
        .append('relations', relations + ":used");

      }
      else
      {
        console.log("Relations is null");
      }
      
      
      const formData = new FormData();
      formData.append("file", this.file);
      //post the new file
      this.http.post("http://localhost:8080/api/v1/file", formData, 
      {
        params: params,
        observe: 'response'})
        .subscribe(response => {
            console.log("logging file post respose");
            console.log(response);
          }, error => {
            console.log("logging file post error");
            console.log(error);
          });

      
      this.fileList.push(
        {
        name: this.fileName,
        uuid: '',
        created: '',
        md5: '',
        relations: relations,
        route: ''
        });
        //refresh the page so the new file shows up
        //window.location.reload();
    }
  }

  emptyClick(): void
  {
    console.log("Button function");
  }
}
