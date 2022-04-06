import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { CookieService } from 'ngx-cookie-service';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

  constructor(private cookieService:CookieService, private http: HttpClient, private router: Router) { }

  login()
  {
    const headers = new HttpHeaders();
    const formData = new FormData();
    
    var email = (<HTMLInputElement>document.getElementById("email")).value;
    var password = (<HTMLInputElement>document.getElementById("password")).value;
    
    //const params = new HttpParams()
    //.append('email', username)
    //.append('password', password);

    formData.append("email", email);
    formData.append("password", password);

    //console.log(formData.getAll("email"));
    
    this.http.post("http://localhost:8080/api/v1/auth/login", formData, { 
      observe: 'response', responseType: 'text'})
      .subscribe(response => {
        //console.log("Logging response");
        //console.log(response.body);
        if(response.body != null)
        {
          var cookieText = response.body.slice(10, response.body.length-2);
          //console.log("Sliced: ", cookieText);
          this.cookieService.set('k', cookieText);
        }
        if(response.status === 200)
        {
          this.router.navigateByUrl('/dashboard', { replaceUrl: true});
          //console.log(this.cookieService.get('test'));
        }
      });
    
    //console.log("Email: ", email);
    //console.log("Password: ", password);
    //this.router.navigateByUrl('/dashboard', { replaceUrl: true});  
}

toggleVisibility()
{
  var x = (<HTMLInputElement>document.getElementById("password"));
  if (x.type === "password")
  {
    x.type = "text";
  }
  else
  {
    x.type = "password";
  }
}

  ngOnInit(): void {}

}
