import cv2 as cv
import numpy as np
import requests



def read_HTTP_into_mat(url : str):
    """
    read_HTTP_into_mat
    streams a URL to an image into a numpy array.
    will raise a value error if the file type is unsupported or if 
    the url doesn't point to an image.

    :param url: string, url to the image file
    :return: np.ndarray, numpy aray image
    """    
    if ".gif" in url:
        raise ValueError("Gifs not supported")

    try: 
        response = requests.get(url, stream=True).raw
        image = np.asarray(bytearray(response.read()), dtype=np.uint8)
        image = cv.imdecode(image, cv.IMREAD_COLOR)
        if image is None:
            raise ValueError("Could not decode the image from the URL, likely doesn't point to an image file")
        return image
    except requests.exceptions.RequestException as e:
        raise ValueError(f"Error fetching image: {e}")
