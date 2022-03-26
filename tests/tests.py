import requests
import inspect
import time
import threading
import testDataGenerator
import time

# uuid for single test.txt file upload, download, and delete
fileUUID : str

# list for thread return values
parallelThreadReturnValues = []
parallelThreadStartFlag = False

# base path for api version
apiBasePath = "http://localhost:8080/api/v1"

def error(reason : str):
    print("[ERROR]: {}".format(reason))

def success():
    print("[Success]")

def pingTest():
    '''Checks for the ping response against the api'''
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        # request ping page
        r = requests.get(url=apiBasePath + "/ping", timeout=20)

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
    try:
        if fileData == None:
            # normal file upload test with sample file
            # print function name
            print(inspect.getframeinfo(inspect.currentframe()).function, end=" ", flush=True)

            # contents of test file
            file = {"file":open("test.txt",'rb')}

            # add params
            params = {
                "filename" : "/home/test/test.txt",
                "casename" : "testcase",
            }

            # upload file
            r = requests.post(url = apiBasePath + "/file", files=file, timeout=20, params=params)

            # check if good request
            if r.status_code != 200:
                error(r.status_code)
            else:
                success()
                global fileUUID
                fileUUID = r.content
        else:
            # parallel upload test with a particular file
            # contents of test file
            file = {"file":fileData}

            # upload file
            r = requests.post(url = apiBasePath + "/file", files=file, timeout=10)

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
        if filename == None:
            # normal download test with sample file
            # print function name
            print(inspect.getframeinfo(inspect.currentframe()).function, end=" ", flush=True)
            
            url = apiBasePath + '/file'

            params = {
                "filename" : fileUUID,
                "casename" : "testcase",
            }

            # download file
            r = requests.get(url, timeout=20, params=params)

            # check if good request
            if r.status_code != 200:
                error(r.status_code)
            else:
                success()
        else:
            # test if a particular file can be downloaded
            url = apiBasePath + '/file'

            params = {"filename" : filename}

            # download file
            r = requests.get(url, timeout=10, params=params)

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
        if filename == None:
            # normal delete test with sample file
            # print function name
            print(inspect.getframeinfo(inspect.currentframe()).function, end=" ", flush=True)

            url = apiBasePath + '/file'

            params = {"filename" : fileUUID}
            
            # request to delete file
            r = requests.delete(url, timeout=20, params=params)

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
            r = requests.delete(url, timeout=20, params=params)

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
        r = requests.get(url=apiBasePath + "/db/test")

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
        r = requests.post(url=apiBasePath + "/db/test")

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
        r = requests.get(url=apiBasePath + "/db/test/find")

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
    """Attempts to add a case to the database
    """
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        # request ping page
        r = requests.get(url=apiBasePath + "/db/test/find")

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
        time.sleep(0.5)