import cv2 as cv
import numpy as np




def Invert(image) -> np.ndarray:
    return cv.bitwise_not(image)



def Saturate(image, value: float) -> np.ndarray:
    
    if value <=0.0:
        raise ValueError(f'Expected saturation value to be greater than 0, got {value}')
    
    color_channel = 2

    if image.shape[color_channel] != 3:
        raise ValueError(f'Expected an RGB image but got {image.shape[color_channel]} channels')

    sat_channel = 1
    hsv_image = cv.cvtColor(image, cv.COLOR_BGR2HSV)

    hsv_image[:,:,sat_channel ] = hsv_image[:,:,sat_channel] * value 
    hsv_image[:, :, sat_channel] = np.clip(hsv_image[:, :, sat_channel], 0, 255)
    
    saturated_image = cv.cvtColor(hsv_image, cv.COLOR_HSV2BGR)
    
    return saturated_image


