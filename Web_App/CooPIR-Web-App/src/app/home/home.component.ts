import { Component, Inject, OnInit } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import Swal from 'sweetalert2';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css']
})
export class HomeComponent implements OnInit {
  fileName = '';
  path: string = "/assets/images/CooPIR_Pic.jpg";
  ImageAlt: string;
  constructor(private http: HttpClient) {
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
            Swal.fire(error.error.toString());
          }
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
            Swal.fire(error.error.toString());
          }
        })
    }
    
  }
  
  pingButton() : void {
    //var response;
    this.http.get<any>("http://localhost:8080/api/v1/ping", {observe: 'response'})
    .subscribe(response => {
      if(response.status === 200)
      {
        console.log("Status code 200 received");
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
    const file:File = event.target.files[0];

        if (file) 
        {
            this.fileName = file.name;
            console.log(file);
            const params = new URLSearchParams();
            params.set("caseName", caseName);
            params.set("fileName", this.fileName);

            const formData = new FormData();

            formData.append("testPic", file);

            this.http.post("http://localhost:8080/api/v1/file", formData, {observe: 'response'})
            .subscribe(response => {
              console.log("logging respose");
              console.log(response);
            }, error => {
              console.log("logging error");
              if(error.status === 0)
              {
                this.alertPopup(error);
              }
            });
           
        }
  } 
}
