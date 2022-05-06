import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { CookieService } from 'ngx-cookie-service';

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
      label: 'Dashboard',
      icon: 'list',
      route: '/dashboard'
    },
    {
      label: 'Jobs',
      icon: 'assignment',
      route: '/jobs'
    },
    {
      label: 'Map',
      icon: 'map',
      route: '/map'
    },
    {
      label: 'File info',
      icon: 'map',
      route: '/fileinfo'
    }
  ];
  constructor(private http: HttpClient, private cookieService:CookieService,  private router: Router) { }

  ngOnInit(): void 
  {
    var fileParams = new HttpParams()
    .append('caseUUID', this.cookieService.get("currentUUID"))
    .append('fileUUID', this.cookieService.get("fileUUID"));

    var fileInfo:any;
    this.http.get("http://localhost:8080/api/v1/file/info", {params: fileParams, observe: 'response'})
    .subscribe( response => {
      console.log("File info request: ", response.body)
      fileInfo = response.body
      console.log("fileInfo: ", fileInfo.file.length)
    })
  }

}
