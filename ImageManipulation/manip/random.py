import cv2 as cv
import numpy as np




def RandomFilter(image: np.ndarray, kernel_size = 3, min = 0, max = 1, normalize = False) -> np.ndarray:

    random_kernel_r = np.random.uniform(min,max,(kernel_size,kernel_size))
    random_kernel_g = np.random.uniform(min,max,(kernel_size,kernel_size))
    random_kernel_b = np.random.uniform(min,max,(kernel_size,kernel_size))

    if normalize:
        random_kernel_r /= kernel_size ** 2
        random_kernel_g /= kernel_size ** 2
        random_kernel_b /= kernel_size ** 2

    ddepth = -1
    b_chan,g_chan,r_chan = cv.split(image)
    filtered_image_r = cv.filter2D(r_chan, ddepth,random_kernel_r)
    filtered_image_g = cv.filter2D(g_chan, ddepth,random_kernel_g)
    filtered_image_b = cv.filter2D(b_chan, ddepth,random_kernel_b)

    filtered_image = cv.merge([filtered_image_b, filtered_image_g, filtered_image_r])
    return filtered_image
