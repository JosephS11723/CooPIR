import { Injectable } from '@angular/core';
import { HttpInterceptor, HttpRequest, HttpHandler } from '@angular/common/http';
import { CookieService } from 'ngx-cookie-service';

@Injectable()
export class AuthHeadersInterceptor implements HttpInterceptor {

  constructor(private cookieService:CookieService) { }

  intercept(request: HttpRequest<any>, next: HttpHandler) {
    console.log('Auth interceptor dude');
    console.log(request.url);
    const authToken = this.cookieService.get('test');
    console.log("Auth token: ", authToken);
    const authReq = request.clone({ setHeaders: { Authorization: authToken } });
    return next.handle(authReq);
  }
}