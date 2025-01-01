import cv2 as cv
import numpy as np
import requests
import io
from fastapi import HTTPException

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
        raise HTTPException(status_code=415, detail=f".gif files are not supported")
    try: 
        response = requests.get(url, stream=True).raw
    except requests.exceptions.RequestException as e:
        raise HTTPException(status_code=422, detail=f"given url doesn't point to an image: {e}")
    image = np.asarray(bytearray(response.read()), dtype=np.uint8)
    image = cv.imdecode(image, cv.IMREAD_COLOR)
    if image is None:
         raise HTTPException(status_code=500, detail=f"failed to decode image")
    return image


def image_to_bytes(image : np.ndarray) -> io.BytesIO:
     """
     image_to_bytes
     encodes a numpy image into bytes in the png format
     :param image: numpy array, our image, in any format
     """
     _, buffer = cv.imencode('.png', image)
     image_bytes = io.BytesIO(buffer)
     image_bytes.seek(0)
     return image_bytes


