import requests
import inspect
import time
import threading
import testDataGenerator
import time
import random

# uuid for single test.txt file upload, download, and delete
fileUUID : str

# uuid for case test
caseuuid : str = ""

# list for thread return values
parallelThreadReturnValues = []
parallelThreadStartFlag = False

# base path for api version
apiBasePath = "http://localhost:8080/api/v1"

# test credentials
email = "default@coopir.edu"
password = "password"

# case information
casename = "test case" + str(random.randint(0, 1000))
caseDescription = '''This is a test case for the CoopIR API made by the\npython test script'''

# create requests session object for cookies
s = requests.Session() 

def error(reason : str):
    print("[ERROR]: {}".format(reason))

def success():
    print("[Success]")

def loginTest():
    '''Checks the login path and gets a cookie'''
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        # request login page
        r = s.post(
            url = apiBasePath + "/auth/login",
            data = {
                "email":email,
                "password":password
            },
            timeout=20
            )

        # check if good request
        if r.status_code != 200:
            error(r.status_code)
            print(r.content)
        else:
            success()
    except Exception as e:
        error(e)

def pingTest():
    '''Checks for the ping response against the api'''
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        # request ping page
        r = s.get(url=apiBasePath + "/ping", timeout=20)

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

def uploadTest(fileData = None):
    '''Attempts to upload a file to the server'''
    global caseuuid
    try:
        if fileData == None:
            # normal file upload test with sample file
            # print function name
            print(inspect.getframeinfo(inspect.currentframe()).function, end=" ", flush=True)

            # contents of test file
            file = {"file":open("test.txt",'rb')}

            # add params
            params = {
                "fileuuid" : "/home/test/test.txt",
                "caseuuid" : caseuuid,
            }

            # upload file
            r = s.post(url = apiBasePath + "/file", files=file, timeout=20, params=params)

            # check if good request
            if r.status_code != 200:
                error(str(r.status_code) + " " + r.content.decode())
            else:
                success()
                global fileUUID
                fileUUID = r.content.decode()
        else:
            # parallel upload test with a particular file
            # contents of test file
            file = {"file":fileData}

            # upload file
            r = s.post(url = apiBasePath + "/file", files=file, timeout=10)

            # check if good request
            if r.status_code != 200:
                return None, r.status_code
            else:
                return r.content, r.status_code
    except Exception as e:
        error(e)
        return None, r.status_code
        
def downloadTest(filename : str = None):
    '''Attempts to download the file we just uploaded'''
    try:
        global caseuuid
        if filename == None:
            # normal download test with sample file
            # print function name
            print(inspect.getframeinfo(inspect.currentframe()).function, end=" ", flush=True)
            
            url = apiBasePath + '/file'

            params = {
                "fileuuid" : fileUUID,
                "caseuuid" : caseuuid,
            }

            # download file
            r = s.get(url, timeout=20, params=params)

            # check if good request
            if r.status_code != 200:
                error(r.status_code)
            else:
                success()
        else:
            # test if a particular file can be downloaded
            url = apiBasePath + '/file'

            params = {"fileuuid" : filename}

            # download file
            r = s.get(url, timeout=10, params=params)

            # check if good request
            if r.status_code != 200:
                error(r.status_code)
                return False
            else:
                return True
        
    except Exception as e:
        error(e)

def downloadTestWithParameters(filename : str = None):
    '''Attempts to download the file we just uploaded'''
    try:
        global caseuuid
        global fileUUID
        if filename == None:
            # normal download test with sample file
            # print function name
            print(inspect.getframeinfo(inspect.currentframe()).function, end=" ", flush=True)
            
            url = "{}/file/{}/{}".format(apiBasePath, fileUUID, caseuuid)

            # download file
            r = s.get(url, timeout=20)#, params=params)

            # check if good request
            if r.status_code != 200:
                error(str(r.status_code) + " " + r.content.decode())
            else:
                success()
        else:
            # test if a particular file can be downloaded
            url = apiBasePath + '/file/' + fileUUID + '/' + caseuuid

            #params = {"filename" : filename}

            # download file
            r = s.get(url, timeout=10)#, params=params)

            # check if good request
            if r.status_code != 200:
                error(r.status_code)
                return False
            else:
                return True
        
    except Exception as e:
        error(e)

def deleteTest(filename : str = None):
    '''Attempts to delete a file from the server'''
    try:
        global caseuuid
        if filename == None:
            # normal delete test with sample file
            # print function name
            print(inspect.getframeinfo(inspect.currentframe()).function, end=" ", flush=True)

            url = apiBasePath + '/file'

            params = {"fileuuid" : fileUUID}
            
            # request to delete file
            r = s.delete(url, timeout=20, params=params)

            # check if good request
            if r.status_code != 200:
                error(r.status_code)
            else:
                success()
        else:
            # test if a particular file can be deleted
            url = apiBasePath + '/file'

            params = {"filename" : filename}

            # request to delete file
            r = s.delete(url, timeout=20, params=params)

            # check if good request
            if r.status_code != 200:
                return False
            else:
                return True
    except Exception as e:
        error(e)

def parallelUploadTest():
    '''Attemps to upload files at the same time using threads and then download all the data back'''
    def printStatus(status):
        print("\r{} [{}]             ".format("parallelUploadTest", status), end=" ", flush=True)

    def parallelUploadAndGet(fileData, index : int):
        '''Uploads a file and adds the return value to the list'''
        global parallelThreadReturnValues
        global parallelThreadStartFlag

        success, errors = uploadTest(fileData)
        parallelThreadReturnValues[index] = [success, errors]

    # print function information
    print(inspect.getframeinfo(inspect.currentframe()).function, end=" ", flush=True)
    global parallelThreadReturnValues

    # the amount of words to be uploaded in parallel
    wordCount = 100

    # setup thread list
    parallelThreadReturnValues = [None]*wordCount
    
    # get test data
    printStatus("Aquiring test data")
    testData = testDataGenerator.getWords(count=wordCount)

    # create a list of threads
    threads = []

    printStatus("Creating threads")

    # create a thread for each file
    for i in range(len(testData)):
        # create a thread
        t = threading.Thread(target=parallelUploadAndGet(testData[i], i))
        # add the thread to the list
        threads.append(t)
    
    printStatus("Starting threads")

    # start all threads
    for i in range(len(threads)):
        threads[i].start()
    
    printStatus("Waiting for threads to finish")

    # join all the threads
    for t in threads:
        t.join()

    printStatus("Checking for upload errors")

    # check if any uploads failed
    for i in range(len(parallelThreadReturnValues)):
        if parallelThreadReturnValues[i][0] == None:
            print(parallelThreadReturnValues[i][1])
            error("Parallel upload failed")
            return

    printStatus("Downloading files")

    # download all the files
    for i in range(len(testData)):
        if not downloadTest(parallelThreadReturnValues[i][0]):
            # download failed
            error("Parallel download failed")
            return

    printStatus("Deleting files")

    # delete all the files
    for i in range(len(testData)):
        if not deleteTest(parallelThreadReturnValues[i][0]):
            error("Parallel deletion failed")
            return
    else:
        # everything successful
        printStatus("")
        print("\r{}".format(inspect.getframeinfo(inspect.currentframe()).function), end=" ", flush=True)
        success()

def dbPingTest():
    '''Attempts to add a User document into the database'''
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        # request ping page
        r = s.get(url=apiBasePath + "/test")

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
    '''Attempts to add a User document into the database'''
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        # request ping page
        r = s.post(url=apiBasePath + "/test")

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
        r = s.get(url=apiBasePath + "/test/find")

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

def dbNewCaseTest():
    '''Attempts to add a case to the database'''
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        global casename

        # request ping page
        r = s.post(
            url=apiBasePath + "/case/new", json={
                    "uuid":None,
	                "name":casename,
                    "dateCreated":"today :D",
                    "viewAccess":"supervisor",
                    "editAccess":"supervisor",
                    "collaborators":["Brandon Ship", "Me lol"]
                }
            )

        global caseuuid
        caseuuid = r.content.decode()

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

def dbUpdateCaseTest():
    """Attempts to update a case to the database
    """
    try:
        global caseuuid
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        # request ping page
        r = s.post(
            url=apiBasePath + "/case/update", json={
                "filter":{"uuid":caseuuid},
                "update":{
                    "uuid":caseuuid,
	                "name":"testcase",
                    "dateCreated":"June 4th, 1776",
                    "viewAccess":"responder",
                    "editAccess":"responder",
                    "collaborators":["Brandon Ship", "Me lol", "Alex Johnson Petty"]
                    }
                }
            )

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


def dbFindCaseTest():
    """Attempts to find a case in the database
    """
    global casename
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        # request ping page
        r = s.get(
            url=apiBasePath + "/case", json={
                "name":casename
                }
            )

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

def dbNewUserTest():
    """Attempts to add a user to the database
    """
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        # request ping page
        r = s.post(
            url=apiBasePath + "/user/new", json={
	                "name":"testuser",
                    "email":"testemail@emailservice.com",
                    "role":"responder",
                    "cases":["The Case"],
                    "password":"football"
                }
            )

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

def dbUpdateUserTest():
    """Attempts to update a user in the database
    """
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        # request ping page
        r = s.post(
            url=apiBasePath + "/user/update", json={
                "filter":{"name":"testuser"},
                "update":{
                    "name":"testcase",
                    "email":"thenewemail@emailservice.com",
                    "role":"responder",
                    "cases":["The Case", "The OTHER Case ;)"],
                    "password":"football123"
                    }
                }
            )

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


def dbFindUserTest():
    """Attempts to find a case in the database
    """
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        # request ping page
        r = s.get(
            url=apiBasePath + "/user", json={
                "name":"testuser"
                }
            )

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

def dbGetUserCasesTest():
    """Attempts to find all user cases in database
    """
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        # request ping page
        r = s.get(
            url=apiBasePath + "/cases"
            )

        # check if good request
        if r.status_code != 200:
            error(r.status_code)
        else:
            success()
            print(r.content, end="")
            
    except Exception as e:
        error(e)

def createJobTest():
    """Attempts to create a new job
    """
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        # request ping page
        r = s.post(
            url=apiBasePath + "/jobs/new",
            json = {
                "arguments":[],
                "name":"test_job_1",
                "jobtype":"log_analysis",
            }
        )

        # check if good request
        if r.status_code != 200:
            error(r.status_code)
        else:
            success()
            print(r.content, end="")
            
    except Exception as e:
        error(e)


def createSearchJobAndFindByUUIDTest():
    """Attempts to create a new job
    """
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        '''
        # request ping page
        r = s.post(
            url=apiBasePath + "/jobs/new",
            json = {
                "arguments":[],
                "name":"test_job_1",
                "jobtype":"log_analysis",
            }
        )
        '''

        # check if good request
        if r.status_code != 200:
            error(r.status_code)
        else:
            success()
            print(r.content, end="")
            
    except Exception as e:
        error(e)

tests = [loginTest, pingTest, dbNewCaseTest, uploadTest, uploadTest, uploadTest, downloadTest, downloadTestWithParameters, dbUpdateCaseTest, dbFindCaseTest, dbNewUserTest, dbUpdateUserTest, dbFindUserTest, dbGetUserCasesTest]
def runAllTests():
    for test in tests:
        test()
        time.sleep(0.5)