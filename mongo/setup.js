conn = new Mongo();


//---------------------------------------------

db = conn.getDB("admin");

db.createUser(
    {
      user: "myUserAdmin",
      pwd: "securepassword",
      roles: [
        { role: "userAdminAnyDatabase", db: "admin" },
        { role: "readWriteAnyDatabase", db: "admin" }
      ]
    }
  )

//---------------------------------------------
db = conn.getDB("Users");

db.createCollection("User");
db.createCollection("Authentication");
db.createCollection("Aggregations");

db.createRole(
    {
      role: "apiUsersDatabase", 
      privileges: [
        { resource: {db: "Users", collection: "User"}, actions: ["find", "insert", "update"] },
        { resource: {db: "Users", collection: "Authentication"}, actions: ["find", "insert", "update"] },
        { resource: {db: "Users", collection: "Aggregation"}, actions: ["find"]}
      ],
      roles: []
    }
 )

//-----------------------------------------------------

db = conn.getDB("Cases");

db.createCollection("Case");
db.createCollection("File");
db.createCollection("Log");
db.createCollection("Aggregations");

db.createRole(
    {
      role: "apiCasesDatabase", 
      privileges: [
        { resource: {db: "Cases", collection: "Case"}, actions: ["find", "insert", "update"] },
        { resource: {db: "Cases", collection: "File"}, actions: ["find", "insert", "update"] },
        { resource: {db: "Cases", collection: "Log"}, actions: ["find", "insert", "update"] },
        { resource: {db: "Cases", collection: "Aggregation"}, actions: ["find"]}
      ],
      roles: []
    }
 )

coll = db.getCollection("Aggregations");

coll.insert(
    {
      "name":"GetFiles",

      "aggregation":[
        {
          "$sort":{
            "$size":"$collaborators"
          }
        },
        {
          "$project":{
            "name":true,
            "number_of_collaborators":{"$size": "$collaborators"},
            "collaborators":true
          }
        }
      ]
    }
  )



//-----------------------------------------------

db = conn.getDB("admin");

db.createUser(
    {
        user: "api",
        pwd: "bingusmalingus",
        roles: [
          { role: "apiCasesDatabase", db: "Cases" },
          { role: "apiUsersDatabase", db: "Users" }
        ]
      }
  )

