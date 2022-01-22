conn = new Mongo();

db = conn.getDB("Users");

db.createCollection("User");
db.createCollection("Authentication");

db = conn.getDB("Cases");

db.createCollection("Case");