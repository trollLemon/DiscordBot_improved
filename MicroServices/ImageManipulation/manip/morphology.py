import cv2 as cv
import numpy as np


def Erode(image: np.ndarray, kernel_size = 3, iterations=1) -> np.ndarray: 
    
    if kernel_size <= 0 or iterations <= 0:
        raise ValueError(f'Expected kernel size or interations to be greater than 0, got {kernel_size} and {iterations}')

    kernel = np.ones((kernel_size,kernel_size), np.uint8)
    eroded_image = cv.erode(image, kernel, iterations=iterations)

    return eroded_image



def Dilate(image: np.ndarray, kernel_size = 3, iterations=1) -> np.ndarray: 
    
    if kernel_size <= 0 or iterations <= 0:
        raise ValueError(f'Expected kernel size or interations to be greater than 0, got {kernel_size} and {iterations}')

    kernel = np.ones((kernel_size, kernel_size ), np.uint8) 
    dilated_image = cv.dilate(image, kernel, iterations=iterations)
   
    return dilated_image





