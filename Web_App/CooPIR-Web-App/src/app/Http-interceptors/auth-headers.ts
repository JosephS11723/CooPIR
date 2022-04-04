import { Injectable } from '@angular/core';
import { HttpInterceptor, HttpRequest, HttpHandler } from '@angular/common/http';
import { CookieService } from 'ngx-cookie-service';

@Injectable()
export class AuthHeadersInterceptor implements HttpInterceptor {

  constructor(private cookieService:CookieService) { }

  intercept(request: HttpRequest<any>, next: HttpHandler) {
    console.log('Auth interceptor dude');
    console.log(request.url);
    const authToken = {token: this.cookieService.get('token')};
    console.log("Auth token: ", authToken);
    var authReq = request.clone({ setHeaders: { Authorization: authToken.toString() }});
    authReq = authReq.clone({withCredentials: true})
    return next.handle(authReq);
  }
}