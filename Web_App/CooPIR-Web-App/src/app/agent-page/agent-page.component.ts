import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';

@Component({
  selector: 'app-agent-page',
  templateUrl: './agent-page.component.html',
  styleUrls: ['./agent-page.component.css']
})
export class AgentPageComponent implements OnInit {

  agentList = new Array<any>();
  menuItems = [
    {
      label: 'Logout',
      icon: 'exit_to_app',
      route: '/login'
    },    
    {
      label: 'Dashboard',
      icon: 'list',
      route: '/dashboard'
    }
  ];


  constructor(private http: HttpClient) { }

  getAgents(): void
  {;
    //var retrievedAgents: any;
    var retrievedAgents: any;
    var jsonTest = new Map<string, JSON>();
    
    //jsonTest.set("test", bluh);
    var agentKeys: any;
    this.http.get("http://localhost:8080/api/v1/agent/all", {observe: 'response'})
    .subscribe(response => {
      console.log("Agent get response: ", response);
      if(response.body != null)
      {
        console.log("Intial response: ", response.body);
        retrievedAgents = response.body;
      }
      //store the keys for each agent
      agentKeys = Object.keys(retrievedAgents);
      console.log("Retrived keys: ", agentKeys);
      jsonTest = retrievedAgents;

      console.log("Retrieved agents: ", Object.values(retrievedAgents))
      var agentInfo = Object.values(retrievedAgents);
      for(var index = 0; index < agentInfo.length; index++)
      {
        var testAgent:any;
        testAgent = agentInfo[index];
        this.agentList.push(
          {
            name: testAgent.name,
            os: testAgent.os,
            arch: testAgent.arch
          }
        )
        console.log("Agent name: ", testAgent.name);
        //console.log("Key retrieved: ", "'" + agentKeys[index] + "'");
        //console.log("Retrieved agents: ", retrievedAgents);
        //console.log("Values in retrieved agents: ", Object.values(retrievedAgents));
        //console.log("Loop: ", retrievedAgents.get("'" + agentKeys[index] + "'"));
      }
      
    })
  }


  launchAgent(): void 
  {

  }
  
  ngOnInit(): void 
  {
    this.getAgents();
  }

}
