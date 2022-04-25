from pymongo import MongoClient
from random import randint

Users = [

    {
        "name":"Alex Petty",
        "email":"alexpetty@yahoo.org",
        "role":"responder",
        "cases":[],
        "saltedHash":"binka"
    },
    {
        "name":"Andrew Merrow",
        "email":"andrewmerrow@yahoo.org",
        "role":"responder",
        "cases":[],
        "saltedHash":"binka"
    },
    {
        "name":"Sarah Siharath",
        "email":"sarahsiharath@yahoo.org",
        "role":"responder",
        "cases":[],
        "saltedHash":"binka"
    },
    {
        "name":"Heath Tims",
        "email":"timbo@yahoo.org",
        "role":"supervisor",
        "cases":[],
        "saltedHash":"binka"
    }

]

################################################

conn = MongoClient("mongodb://api:bingusmalingus@localhost:27017")


coll = conn["Users"]["User"]

coll.insert_many(Users)


coll = conn["Cases"]["Case"]

for i in range(0,10000):
    coll.insert_one(
        {
            "name": "case"+str(i),

            "dateCreated": "Decembruary " + str(randint(1,31)) + ", 1969",

            "viewaccess": "responder",

            "editAccess": "supervisor",

            "collabs": ["Heath Tims", "Andrew Merrow"]
        }
    )

