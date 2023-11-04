from PyDictionary import PyDictionary

dictionary = PyDictionary()

def handle(word):
    """handle a request to the function
    Args:
        req (str): request body
    """

    return dictionary.meaning(word)
