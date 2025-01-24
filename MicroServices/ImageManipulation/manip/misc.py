import numpy as np
from math import floor, sqrt

def Shuffle(image : np.ndarray, partitions: int = 10) -> np.ndarray:
    
    if partitions <= 1:
        raise ValueError(f'Expected partitions to be greater than 1, got {partitions}')

    height, width, chan = image.shape

    rows = floor(sqrt(partitions))  
    cols = floor(partitions / rows)  

    slice_width  = width // cols
    slice_height = height // rows

    slices = []
    for i in range(rows):
        for j in range(cols):
            row_range = i * slice_height
            col_range = j * slice_width

            img_slice = image[row_range:min(row_range + slice_height, height),
                             col_range:min(col_range + slice_width, width)]

            slices.append(img_slice)

    np.random.shuffle(slices)

    new_height = min(rows * slice_height, height)
    new_width = min(cols * slice_width, width)

    shuffled_image = np.zeros((new_height, new_width, chan))

    for idx, img_slice in enumerate(slices):
        row_idx = idx // cols
        col_idx = idx % cols

        row_start = row_idx * slice_height
        col_start = col_idx * slice_width

        shuffled_image[row_start:row_start + img_slice.shape[0],
                       col_start:col_start + img_slice.shape[1]] = img_slice
    
    return shuffled_image
