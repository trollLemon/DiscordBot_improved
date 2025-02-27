import cv2 as cv
import numpy as np


def EdgeDetect(image : np.ndarray, t_lower: int = 100, t_upper : int = 200) -> np.ndarray: 
    
    if t_lower < 0 or t_upper < 0:
        raise ValueError(f'Expected lower and upper values to be greater or equal to 0, got {t_lower} and {t_upper}')
    
    color_channel = 2
    rgb_channels  = 3
    
    gray_image = cv.cvtColor(image, cv.COLOR_BGR2GRAY) if image.shape[color_channel] == rgb_channels else image
    
    edges = cv.Canny(gray_image, threshold1=t_lower, threshold2=t_upper)

    return edges





