import requests
import inspect
import json

def error(reason : str):
    print("[ERROR]: {}")

def success():
    print("[Success]")

def pingTest():
    """Checks for the ping response against the api
    """
    # print function name
    print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

    # request ping page
    r = requests.get(url="http://localhost:8080/ping")

    # check if returned value
    if r.status_code != 200:
        error(r.status_code)
    
    # check if returned value is correct
    if r.json()["data"] == "pong":
        success()
    else:
        error()




tests = [pingTest]

def runAllTests():
    for test in tests:
        test()