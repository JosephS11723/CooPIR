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
    # print function name
    print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

    # contents of test file
    file = {"file":open("test.txt",'rb')}

    # upload file
    r = requests.post(url = "http://localhost:8080/file", files=file)

    # check if good request
    if r.status_code != 200:
        error(r.status_code)
        pprint(requests.post(url = "http://localhost:8080/file", files=file).request.body)
    else:
        success()

def downloadTest():
    """Attemps to download the file we just uploaded
    """
    # print function name
    print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")
    
    url = 'http://localhost:8080/file/test.txt'

    myfile = requests.get(url)
    print(myfile.content.decode())
    print(myfile.status_code)


tests = [pingTest, uploadTest, downloadTest]

def runAllTests():
    for test in tests:
        test()