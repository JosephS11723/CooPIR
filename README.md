# CooPIR 

## What is CooPIR?

  The Coordination Platform for Incident Response (CooPIR) is a software that combines many existing
  workflows into a single platform and creates a new way to collaboratively link evidence together.
  The platform will allow for quick response and coordination when investigating and reacting to
  an incident.
  
## Components

### Docker:
  
  This application utilizes docker containers that contain the services needed for the operation of the application.
  
### MongoDB:
  
  MongoDB is used for the database which contains all the user data and case data.
  
### File Storage:
  SeaweedFS is sed as the file storage which allows for practically unlimited distributed storage.
  
### API:
  TBD
  
### Web Interface:
  Using Angular, CooPIR can be used via a web application which allows users to create, modify, and view cases within a web browser.
  
### Quick Start:
  
  Run the runCompose.sh or .bat depending on the OS to build the docker containers.
  Run the Angular command (TBD)
