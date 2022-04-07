FROM node:16 as node

WORKDIR /app

COPY . .

RUN npm install

RUN npm audit fix --force

RUN npm run build --prod

# RUN ng serve -o --liveReload=true

#create a container to serve the app
FROM nginx:alpine

COPY --from=node /app/dist/CooPIR-Web-App /usr/share/nginx/html

# COPY --from=node /app/nginx.conf /etc/nginx/nginx.conf

EXPOSE 80

# CMD ["nginx", "-g", "daemon off;"]