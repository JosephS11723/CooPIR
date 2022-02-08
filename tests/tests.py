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

def dbFindTest():
    """Attempts to add a User document into the database
    """
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        # request ping page
        r = requests.get(url="http://localhost:8080/db/test/find")

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


tests = [pingTest, uploadTest, downloadTest, deleteTest, dbPingTest, dbInsertTest, dbFindTest]

def runAllTests():
    for test in tests:
        test()