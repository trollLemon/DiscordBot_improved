import cv2 as cv
import numpy as np

def Reduce(image : np.ndarray, quality: float) -> np.ndarray:
    
    if quality <= 0.0:
        raise ValueError(f'Expected quality percentage to be greater than 0%, got {quality}')

    original_height, original_width = image.shape[:2]

    new_width = int(original_width * quality)
    new_height = int(original_height * quality)
    print(image.shape)
    print(quality)
    if new_height == 0 or new_width == 0:
        raise ValueError("Shrinking provided image resulted in a 0 width or height, cannot continue")

    scaled_down_image = cv.resize(image, (new_width, new_height), interpolation=cv.INTER_LINEAR)

    scaled_up_image = cv.resize(scaled_down_image, (original_width, original_height), interpolation=cv.INTER_LINEAR) 
    
    return scaled_up_image






