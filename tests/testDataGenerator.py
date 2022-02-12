import os
from random import shuffle

def getWords(count : int = 100):
    '''Returns a certain number of words from lorem ipsum'''
    # open the file
    with open("LoremIpsum.txt", "r") as f:
        # read the first count words
        words = f.read().replace("\n","").replace(",","").replace(".","").split(" ")[:count]
    
    # return the words
    return words


if __name__ == "__main__":
    k : list
    with open("LoremIpsum.txt", 'r') as f:
        k = f.read().replace("\n", "").split(" ")

    print("Loaded {} words".format(len(k)))
    print("Writing to {}".format(os.path.join("Data", "parallelUploadTest")))

    # put the words into Data/parallelUploadTest in different files
    for i in range(len(k)):
        with open(os.path.join("Data", os.path.join("parallelUploadTest","{}.txt".format(i))), 'w') as f:
            f.write(k[i])