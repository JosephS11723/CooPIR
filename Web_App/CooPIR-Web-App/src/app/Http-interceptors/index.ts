import { HTTP_INTERCEPTORS } from "@angular/common/http";
import { AuthHeadersInterceptor } from './auth-headers';

export const httpInterceptProviders = [
  { provide: HTTP_INTERCEPTORS, useClass: AuthHeadersInterceptor, multi: true }
];