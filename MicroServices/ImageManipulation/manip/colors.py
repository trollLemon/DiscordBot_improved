import cv2 as cv
import numpy as np




def Invert(image) -> np.ndarray:
    return cv.bitwise_not(image)



def Saturate(image, value: float) -> np.ndarray:
   
    sat_channel = 1
    hsv_image = cv.cvtColor(image, cv.COLOR_BGR2HSV)

    hsv_image[:,:,sat_channel ] = hsv_image[:,:,sat_channel] * value 
    hsv_image[:, :, sat_channel] = np.clip(hsv_image[:, :, sat_channel], 0, 255)
    
    saturated_image = cv.cvtColor(hsv_image, cv.COLOR_HSV2BGR)
    
    return saturated_image


