import argparse
import requests
import inspect
import json
from pprint import pprint

def error(reason : str):
    print("[ERROR]: {}".format(reason))

def success():
    print("[Success]")

def pingTest():
    """Checks for the ping response against the api
    """
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        # request ping page
        r = requests.get(url="http://localhost:8080/ping")

        # check if good request
        if r.status_code != 200:
            error(r.status_code)
        
        # check if returned value is correct
        if r.json()["data"] == "pong":
            success()
        else:
            error()
    except Exception as e:
        error(e)

def uploadTest():
    """Attempts to upload a file to the server
    """
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        # contents of test file
        file = {"file":open("test.txt",'rb')}

        # upload file
        r = requests.post(url = "http://localhost:8080/file", files=file)

        # check if good request
        if r.status_code != 200:
            error(r.status_code)
        else:
            success()
    except Exception as e:
        error(e)
        
def downloadTest():
    """Attempts to download the file we just uploaded
    """
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")
        
        url = 'http://localhost:8080/file/test.txt'

        # download file
        r = requests.get(url)

        # check if good request
        if r.status_code != 200:
            error(r.status_code)
        else:
            success()
        
    except Exception as e:
        error(e)
    


def deleteTest():
    """Attempts to delete a file from the server and then test to see if it exists
    """
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        url = 'http://localhost:8080/file/test.txt'

        # request to delete file
        r = requests.delete(url)

        # check if good request
        if r.status_code != 200:
            error(r.status_code)
        else:
            success()

    except Exception as e:
        error(e)

def dbPingTest():
    """Attempts to add a User document into the database
    """
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        # request ping page
        r = requests.get(url="http://localhost:8080/db/test")

        # check if good request
        if r.status_code != 200:
            error(r.status_code)

        # check if good request
        if r.status_code != 200:
            error(r.status_code)
        else:
            success()
            
    except Exception as e:
        error(e)

def dbInsertTest():
    """Attempts to add a User document into the database
    """
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        # request ping page
        r = requests.post(url="http://localhost:8080/db/test")

        # check if good request
        if r.status_code != 200:
            error(r.status_code)

        # check if good request
        if r.status_code != 200:
            error(r.status_code)
        else:
            success()
            
    except Exception as e:
        error(e)

DEBUG = False

#create parser
parser = argparse.ArgumentParser()

#specify all possible arguments
parser.add_argument('-ping', action="store_true", default=False, help="Enable Ping Test")
parser.add_argument('-upload', action="store_true", default=False, help="Enable Upload Test")
parser.add_argument('-download', action="store_true", default=False, help="Enable Download Test")
parser.add_argument('-delete', action="store_true", default=False, help="Enable Delete Test")
parser.add_argument('-dbping', action="store_true", default=False, help="Enable DBping Test")
parser.add_argument('-dbinsert', action="store_true", default=False, help="Enable DBinsert Test")

# Parse and print the results
args = parser.parse_args()
#Print the status of each flag
if DEBUG:
    print("Ping: {}".format(args.ping))
    print("Upload: {}".format(args.upload))
    print("Download: {}".format(args.download))
    print("Delete: {}".format(args.delete))
    print("DBping: {}".format(args.dbping))
    print("DBinsert: {}".format(args.dbinsert))

#all true flads add their test to the test list
tests = []
if args.ping:
    tests.append(pingTest)
if args.upload:
    tests.append(uploadTest)
if args.download:
    tests.append(downloadTest)
if args.delete:
    tests.append(deleteTest)
if args.dbping:
    tests.append(dbPingTest)
if args.dbinsert:
    tests.append(dbInsertTest)

#run specified tests
def runAllTests():
    for test in tests:
        test()
