import cv2 as cv
import numpy as np


def edgeDetect(image : np.ndarray, t_lower: int = 100, t_upper : int = 200) -> np.ndarray: 
    gray_image = cv.cvtColor(image, cv.COLOR_BGR2GRAY)

    edges = cv.Canny(gray_image, threshold1=t_lower, threshold2=t_upper)

    return edges





