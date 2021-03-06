import { Component, OnInit } from '@angular/core';
import * as Vis from 'vis';


@Component({
  selector: 'app-map-test',
  templateUrl: './map-test.component.html',
  styleUrls: ['./map-test.component.css']
})
export class MapTestComponent implements OnInit {

  menuItems = [
    {
      label: 'Logout',
      icon: 'exit_to_app',
      route: '/login'
    },    
    {
      label: 'Home',
      icon: 'home',
      route: '/home'
    },
    {
      label: 'Upload',
      icon: 'note_add',
      route: '/home'
    },
    {
      label: 'Dashboard',
      icon: 'assignment',
      route: '/dashboard'
    }
  ];

  fileList = [
    {
      name: '',
      created: '',
      md5: '',
      route: ''
    }
  ];
//nodes = [];
//edges = null;
//network = null;

  draw(): void
  {
    // create people.
    // value corresponds with the age of the person
    //var nodes = [
    //  { id: 1, value: 2, label: "Algie" },
    //  { id: 2, value: 31, label: "Alston" },
    //  { id: 3, value: 12, label: "Barney" },
    //  { id: 4, value: 16, label: "Coley" },
    //  { id: 5, value: 17, label: "Grant" },
    //  { id: 6, value: 15, label: "Langdon" },
    //  { id: 7, value: 6, label: "Lee" },
    //  { id: 8, value: 5, label: "Merlin" },
    //  { id: 9, value: 30, label: "Mick" },
    //  { id: 10, value: 18, label: "Tod" },
    //];

    var nodes = [
      {id: "test1.txt", value: 1, label: "Test1"},
      {id: "test2.txt", value: 1, label: "Test2"},
      {id: "test3.txt", value: 1, label: "Test3"}
    ]
  
    // create connections between people
    // value corresponds with the amount of contact between two people
    //var edges = [
    //  { from: 2, to: 8, value: 3 },
    //  { from: 2, to: 9, value: 5 },
    //  { from: 2, to: 10, value: 1 },
    //  { from: 4, to: 6, value: 8 },
    //  { from: 5, to: 7, value: 2 },
    //  { from: 4, to: 5, value: 1 },
    //  { from: 9, to: 10, value: 2 },
    //  { from: 2, to: 3, value: 6 },
    //  { from: 3, to: 9, value: 4 },
    //  { from: 5, to: 3, value: 1 },
    //  { from: 2, to: 7, value: 4 },
    //];

    var edges = [
      {from: "test1.txt", to: "test2.txt", value: 5}
    ]
  
    // Instantiate our network object.
    var container = document.getElementById("mynetwork");
    var data = {
      nodes: nodes,
      edges: edges,
    };
    var options = {
      nodes: {
        shape: "dot",
        scaling: {
          customScalingFunction: function(min:any, max:any, total:any, value:any) {
            return value / total;
          },
          min: 5,
          max: 150,
        },
      },
    };
    console.log("Container is: ", container);
    if(container != null)
    {
      var network = new Vis.Network(container, data, options);
    }
  }

  constructor() {}
 
  ngOnInit() 
  {
    this.draw();
  }
 
  
}
