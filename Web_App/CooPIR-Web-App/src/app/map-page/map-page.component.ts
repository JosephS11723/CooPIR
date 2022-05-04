import { Component, OnInit } from '@angular/core';
import { HttpClient, HttpParams, HttpHeaders } from '@angular/common/http';
import { CookieService } from 'ngx-cookie-service';
import * as Vis from 'vis';
import Swal from 'sweetalert2';

@Component({
  selector: 'app-map-page',
  templateUrl: './map-page.component.html',
  styleUrls: ['./map-page.component.css']
})
export class MapPageComponent implements OnInit {
  menuItems = [
    {
      label: 'Logout',
      icon: 'exit_to_app',
      route: '/login'
    },
    {
      label: 'Back to Case',
      icon: 'keyboard_tab',
      route: '/case'
    }
  ];
  nodes = new Array<any>();
  //nodes = new Vis.DataSet();
  edges = new Array<any>();
  //edges = new Vis.DataSet();
  constructor(private http: HttpClient, private cookieService:CookieService) { }


  getMapInfo(): void 
  {
    //Swal.fire({
    //  icon: 'info',
    //  title: 'Building Map',
    //  didOpen: () => {
    //    Swal.showLoading()
    //  }
    //});
    //var nodes = new Array<any>();
    //var edges = [{}];

    const params = new HttpParams()
    .append('uuid', this.cookieService.get("currentUUID"));

    this.http.get("http://localhost:8080/api/v1/case/files", {params: params, observe: 'response'})
    .subscribe(response => {
      let retrievedFiles: any
      if(response.body != null)
      {
        retrievedFiles = response.body;
        for(var index = 0; index < retrievedFiles.files.length; index++)
        {
          var fileParams = new HttpParams()
          .append('caseUUID', this.cookieService.get("currentUUID"))
          .append('fileUUID', retrievedFiles.files[index]);
          let fileInfo: any;

          this.http.get("http://localhost:8080/api/v1/file/info", {params: fileParams, observe: 'response'})
          .subscribe(response => {
            fileInfo = response.body;

            //this.nodes.add({ id: fileInfo.file.uuid, value: 10, label: fileInfo.file.filename.split("/").pop()}); 
            this.nodes.push({ id: fileInfo.file.uuid, value: 10, label: fileInfo.file.filename.split("/").pop()});               
            var relations = fileInfo.file.relations;

            if(relations != '')
            {
              for(var index = 1; index < relations.length; index++)
              {
                //this.edges.add({from: fileInfo.file.uuid, to: relations[index].split(":")[0], value: 1})
                this.edges.push({from: fileInfo.file.uuid, to: relations[index].split(":")[0], value: 1})
              }
              //this.edges.add({from: fileInfo.file.uuid, to: relations[0].split(":")[0], value: 1});
              this.edges.push({from: fileInfo.file.uuid, to: relations[0].split(":")[0], value: 1});
            }
            //console.log("redrawing network")
    
          
            //Swal.close();
            //display the map
            var container = document.getElementById("mynetwork");
            var data = {
              nodes: this.nodes,
              edges: this.edges,
            };
            var options = {
              nodes: {
                shape: "dot",
                scaling: {
                  customScalingFunction: function(min:any, max:any, total:any, value:any) {
                    return value / total;
                  },
                  min: 10,
                  max: 150,
                },
              },
              edges: {
                scaling: {
                  customScalingFunction: function(min:any, max:any, total:any, value:any) {
                    return value / total;
                  },
                  min: 0,
                  max: 150,
                },
              },                
              physics: {
                enabled: true,
                barnesHut: {
                  gravitationalConstant: -100000
                }
              },
              layout:
              {
                improvedLayout: false
              }
            };
            if(container != null)
            {
              var network = new Vis.Network(container, data, options);
              //console.log("Map is ready")
            }

          })
        }
      }
    })
    
  }

  buildMap(): void 
  {
    console.log("Nodes: ", this.nodes);
    console.log("Eges: ", this.edges);
    //display the map
    var container = document.getElementById("mynetwork");
    var data = {
       nodes: this.nodes,
       edges: this.edges,
     };
     var options = {
       nodes: {
         shape: "dot",
         scaling: {
           customScalingFunction: function(min:any, max:any, total:any, value:any) {
             return value / total;
           },
           min: 0,
           max: 150,
         },
       },
       physics: {
         enabled: true,
         barnesHut: {
           gravitationalConstant: -100000
         }
       }
     };
     if(container != null)
     {
       var network = new Vis.Network(container, data, options);
       console.log("Map is ready")
     }
     
  }

  ngOnInit(): void 
  {
    //console.log("Gathering info");
    this.getMapInfo();
    //console.log("Map info gathered");
    //console.log("Displaying map");
    //this.buildMap();
  }

}
