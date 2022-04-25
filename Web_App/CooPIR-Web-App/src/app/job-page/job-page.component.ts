import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';

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
      label: 'Dashboard',
      icon: 'assignment',
      route: '/dashboard'
    }
  ];
  constructor(private http: HttpClient) { }


  getJobs(): void 
  {
    this.http.get("http://localhost:8080/api/v1/jobs/types")
    .subscribe(response => {
      console.log("Job response", response);
    });
  }


  ngOnInit(): void
  {
    this.getJobs();
  }

}
