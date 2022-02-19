from PyInquirer import prompt
from examples import custom_style_2
from prompt_toolkit.validation import Validator, ValidationError
import argparse
import requests
import inspect
import json
from pprint import pprint
import config
import pyfiglet

#not exactly sure how the validator works yet
class NumberValidator(Validator):

    def validate(self, document):
        try:
            int(document.text)
        except ValueError:
            raise ValidationError(message="Please enter a number",
                                  cursor_position=len(document.text))


#options to be given to the user
questions = [
    {
        'type': 'list',
        'name': 'user_option',
        'message': 'Welcome to CooPIR',
        'choices': ["ping","upload","delete", "exit"]
    },

]

uploadQuestions = [
    {
        'type': 'input',
        'name': 'file_to_upload',
        'message': 'Enter the file name to upload:'
    },
    {
        'type': 'input',
        'name': 'case_name',
        'message': 'Enter the case name to upload to:'
    }
]

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
        r = requests.get(url=config.apiBasePath + "/ping")

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

def uploadTest(fileName, caseName):
    """Attempts to upload a file to the server
    """
    try:
        # print function name
        print(inspect.getframeinfo(inspect.currentframe()).function, end=" ")

        # contents of test file
        file = {"file":open(fileName,'rb')}

        # add params
        params = {
            "casename" : caseName,
            "filename" : fileName,
        }

        # upload file
        r = requests.post(url = config.apiBasePath + "/file", files=file, params=params)

        # check if good request
        if r.status_code != 200:
            error(r.status_code)
        else:
            success()
    except Exception as e:
        error(e)

def main():
    # print the ascii art
    print(pyfiglet.figlet_format("CooPIR") + '''                     /^._
       ,___,--~~~~--' /'
       `~--~\ )___,)/'
           (/\\_  (/\\_''')

    while(True):
        #get function choice
        answers = prompt(questions)

        if answers.get("user_option") == "ping":
            pingTest()
        
        elif answers.get("user_option") == "upload":
            #ask user for the file 
            uploadAnswers = prompt(uploadQuestions, style=custom_style_2)
            #retrieve the specified file name
            uploadFile = uploadAnswers.get("file_to_upload")
            caseName = uploadAnswers.get("case_name")
            uploadTest(uploadFile, caseName)
        
        #to be implemented if approved
        elif answers.get("user_option") == "delete":
            pass
        
        elif answers.get("user_option") == "exit":
            print("Thank you for using CooPIR!")
            break
        


if __name__ == "__main__":
    main()