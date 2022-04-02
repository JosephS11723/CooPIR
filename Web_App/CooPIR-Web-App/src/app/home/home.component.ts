import { Component, Inject, OnInit } from '@angular/core';
import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import Swal from 'sweetalert2';
import { Router } from '@angular/router';
import { CookieService } from 'ngx-cookie-service';
//import { Server } from 'http';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
  fileName = '';
  path: string = "/assets/images/CooPIR_Pic.jpg";
  ImageAlt: string;
  constructor(private cookieService:CookieService, private http: HttpClient, private router: Router) {
    //this.ImagePath = 'src/app/img/CooPir_Pic.jpg'
    this.ImageAlt = 'Fox dude'
   }

  ngOnInit(): void {}

  alertPopup(error:any) {
    switch(error.status)
    {
      case 0:
        Swal.fire({
          icon: 'error',
          title: 'Unknown Error Occured',
          showCancelButton: true,
          confirmButtonText: 'More Info'
        }).then((result) => 
        {
          if(result.isConfirmed)
          {
            Swal.fire("Error status code: ", error.status.toString());
          }
        })
        break;
      case 200:
          Swal.fire({
            icon: 'success',
            title: 'Success!'
          })
          break;
      case 400:
        Swal.fire({
          icon: 'error',
          title: 'Bad Request',
          showCancelButton: true,
          confirmButtonText: 'More Info'
        }).then((result) => 
        {
          if(result.isConfirmed)
          {
            Swal.fire("Error status code: ", error.status.toString());
          }
        })
        break;
    }
    
  }
  
  pingButton() : void {
    //var response;

    const header = new HttpHeaders()
    .set('content-type', 'application/json');
    this.http.get<any>("http://localhost:8080/api/v1/ping", {
      observe: 'response',
      'headers': header
    })
    .subscribe(response => {
      if(response.status === 200)
      {
        this.alertPopup(response);
      }
      console.log("Logging status");
      console.log(response.status);
  }, error => {
    console.log("Logging error");
    console.log(error);
    this.alertPopup(error);    
  });

  }

  onFileSelected(event:any) {
    var caseName = (<HTMLInputElement>document.getElementById("caseName")).value;
    console.log(caseName);
    const file = event.target.files[0];

        if (file) 
        {
            this.fileName = file.name;
            console.log(file);

            console.log(this.cookieService.get('test'));

            //const params = new URLSearchParams();
            //params.set("caseName", caseName);
            //params.set("fileName", this.fileName);

            const params = new HttpParams()
            .append('casename', caseName)
            .append('filename', this.fileName);

            const headers = new HttpHeaders()
            .append('Cookie', this.cookieService.get('test'));
            //.set('Access-Control-Allow-Origin', '*');


            const formData = new FormData();

            formData.append("file", file);
            //formData.

            this.http.post("http://localhost:8080/api/v1/file", formData, 
            {
              params: params,
              headers: headers,
              withCredentials: true,
              observe: 'response'})
            .subscribe(response => {
              console.log("logging respose");
              console.log(response);
            }, error => {
              console.log("logging error");
              this.alertPopup(error);
              console.log(error.status);
            });
           
        }
  } 
  goToLogin(){
    console.log("Going to login");
    this.router.navigateByUrl('/login', { replaceUrl: true});  
}
goToDashboard(){
  console.log("Going to Dashboard");
  this.router.navigateByUrl('/dashboard', { replaceUrl: true});  
}
}
