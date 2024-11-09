import cv2 as cv
import numpy as np

def Reduce(image : np.ndarray, quality: float) -> np.ndarray:

    original_height, original_width = image.shape[:2]

    new_width = int(original_width * quality)
    new_height = int(original_height * quality)

    scaled_down_image = cv.resize(image, (new_width, new_height), interpolation=cv.INTER_LINEAR)

    scaled_up_image = cv.resize(scaled_down_image, (original_width, original_height), interpolation=cv.INTER_LINEAR) 
    
    return scaled_up_image






