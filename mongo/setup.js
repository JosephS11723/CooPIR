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

db.createCollection("UserMetadata");

db.createRole(
    {
      role: "apiUsersDatabase", 
      privileges: [
        { resource: {db: "Users", collection: "UserMetadata"}, actions: ["find", "insert", "update", "remove", "listCollections"] },
      ],
      roles: []
    }
 )

coll = db.getCollection("UserMetadata");

coll.insert(
  { 
    "uuid":"00000000-0000-0000-0000-000000000001",
    "name":"default",
    "email":"default@coopir.edu",
    "role":"admin",
    "cases":[],
    "saltedhash":"$2a$10$uogPgLI8U5T1I5Ud7lvbMeb6i5nDApNYU4YVhv/.GP16gUqTTfVsC"
  }
)

//-----------------------------------------------------

db = conn.getDB("Cases");

db.createCollection("CaseMetadata");


db.createRole(
    {
      role: "apiCasesDatabase", 
      privileges: [
        { resource: {db: "Cases", collection: ""}, actions: ["find", "insert", "update", "remove", "listCollections"] },
      ],
      roles: []
    }
 )

//-----------------------------------------------
db = conn.getDB("Jobs");

db.createCollection("JobQueue");
db.createCollection("WorkerInfo");
db.createCollection("JobResults");

coll = db.getCollection("JobQueue")

coll.createIndex(
  {"jobtype":1},
  {"starttime":-1}
);


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


// create logs database and collection
db = conn.getDB("Logs")
db.createCollection("Logs")
