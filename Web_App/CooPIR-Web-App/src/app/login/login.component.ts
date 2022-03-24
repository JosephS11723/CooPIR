import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

  constructor(private router: Router) { }

  login(){
    var username = (<HTMLInputElement>document.getElementById("username")).value;
    var password = (<HTMLInputElement>document.getElementById("password")).value;
    console.log("Going to home");
    console.log("Username: ", username);
    console.log("Password: ", password);
    this.router.navigateByUrl('/home', { replaceUrl: true});  
}

  ngOnInit(): void {}

}
