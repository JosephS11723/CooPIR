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

  db.createUser(
    {
        user: "api",
        pwd: "bingusmalingus",
        roles: [
          { role: "readWrite", db: "Cases" },
          { role: "readWrite", db: "Users" }
        ]
      }
  )

//---------------------------------------------
db = conn.getDB("Users");

db.createCollection("User");
db.createCollection("Authentication");
db.createCollection("Aggregations");

db = conn.getDB("Cases");

db.createCollection("Case");
db.createCollection("Aggregations");

//-----------------------------------------------