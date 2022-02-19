FROM mongo:latest

EXPOSE 27017

COPY setup.js /docker-entrypoint-initdb.d/

COPY mongod.conf /etc/mongod.conf