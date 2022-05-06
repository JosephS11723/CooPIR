# CooPIR

## Table of Contents

- [What Is CooPIR?](##what-is-coopir)
- [Quick Start](##quick-start)
- [Components](##components)
  - [Docker](###docker)
  - [API](###api)
  - [FileStorage](###file-storage)
  - [MongoDB](###mongodb)
  - [Web Interface](###web-interface)

## What is CooPIR?

  The Coordination Platform for Incident Response (CooPIR) is a software that combines many existing
  workflows into a single platform and creates a new way to collaboratively link evidence together.
  The platform will allow for quick response and coordination when investigating and reacting to
  an incident.

## Quick Start
  **MUST HAVE DOCKER INSTALLED**
  Run runCompose.sh or runCompose.bat depending on the OS to build the docker containers.
  To run the dev web application, run the webAppServeWithReload.bat which allows for automatic reload when changes are made to the web application files.

## Components

### Docker

  This application utilizes docker containers that contain the services needed for the operation of the application.

### API

  The API is the intermediary interface that connects all the services together and handles all traffic to and from the server and the services. It has a job server and job worker that can process and filter information that could help users find specific information about a file and its contents.

### File Storage

  SeaweedFS is used as the file storage which allows for practically unlimited distributed storage and encrypted storing of files and data.

### MongoDB

  MongoDB is the database used by CooPIR to store a variety of data, including logs, case metadata, and user login information.

### Web Interface

  Using Angular, CooPIR can be used via a web application which allows users to create, modify, and view cases within a web browser. It also contains a relations map which can show relations between files and allows for downloading of files
