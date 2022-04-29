import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { CookieService } from 'ngx-cookie-service';
import Swal from 'sweetalert2';

@Component({
  selector: 'app-make-case',
  templateUrl: './make-case.component.html',
  styleUrls: ['./make-case.component.css']
})
export class MakeCaseComponent implements OnInit {
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

  constructor(private http: HttpClient, private cookieService:CookieService) { }

  submitCase(): void 
  {
    var caseName = (<HTMLInputElement>document.getElementById("caseName")).value;
    var viewaccess = (<HTMLInputElement>document.getElementById("viewaccess")).value;
    var editAccess = (<HTMLInputElement>document.getElementById("editAccess")).value;
    var description = (<HTMLInputElement>document.getElementById("description")).value;
    //console.log("Description: ", description);
  
    const params = new HttpParams()
    .append('uuid', '')
    .append('name', caseName)
    .append('description', description)
    .append('dateCreated', '')
    .append('viewaccess', viewaccess)
    .append('editAccess', editAccess)
    .append('collabs', "Joe")
    .append('collabs', "Momma");

    var newCase = {
      uuid: '',
      name: caseName,
      description: description,
      dateCreated: '',
      viewaccess: viewaccess,
      editAccess: editAccess,
      collabs: []
    }

    //console.log("Our brand new case: ", newCase);

    this.http.post("http://localhost:8080/api/v1/case/new", '', {observe: 'response', params: params})
    .subscribe(response => {
      console.log("New Case Response: ", response);
    }, error => {
      if(error.status === 200)
      {
        Swal.fire({
          icon: 'success',
          title: 'Success',
          text: 'Your case has been created',
          timer: 1500,
          timerProgressBar: true
        })
      }
    })
  }
  ngOnInit(): void {
  }

}
