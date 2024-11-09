"""
Image Manipulation Server
This file contains implementations of a REST API for image manipulation.
The server expects encoded URLs to image files, rather than accepting an image file from HTTP,
though the server does return bytes of a png image after processing.

All images are stored im memory without any disk IO operations.


The server provides a RESTful api with the following endpoints:
    - /api/randomFilteredImage/{image_link:path}/{kernel_size}/{low}/{high}/{kernel_type}
    - /api/invertedImage/{image_link:path}
    - /api/saturatedImage/{image_link:path}/{saturation}
    - /api/edgeImage/{image_link:path}/{lower}/{higher}
    - /api/dilatedImage/{image_link:path}/{box_size}/{iterations}
    - /api/erodedImage/{image_link:path}/{box_size}/{iterations}
    - /api/textImage/{image_link:path}/{text}/{font_scale}/{x}/{y}
    - /api/reducedImage/{image_link:path}/{quality}

The API supports image manipulation operations such as:
    - morphology
    - quality 
    - adding text
    - saturation
    - edge detection
    - random filtering
    - color inverting
"""


from manip import Invert, Saturate, EdgeDetect, Dilate, Erode, Reduce, Add_text, RandomFilter
from util import read_HTTP_into_mat, image_to_bytes
from fastapi import FastAPI, HTTPException

from fastapi.responses import StreamingResponse


from fastapi import FastAPI

app = FastAPI()

@app.get("/api/randomFilteredImage/{image_link:path}/{kernel_size}/{low}/{high}/{kernel_type}")
async def filter_image_random(image_link: str, kernel_size: int, low: int, high: int, kernel_type: str):
    """
    filter_image_random
    returns bytes in png format of an image with a random filter applied to each color channel
    :param image_link:string  : encoded url to an image file
    :param kernel_size:int    : size for the kernel ( also called a filter)
    :param low:int            : lower bound for RNG
    :param high:int           : upper bound for RNG
    :param kernel_type:string : type of kernel, normalized(norm) or not normalized (raw)
    """
    image = read_HTTP_into_mat(image_link)
    
    if kernel_size <= 0:
        raise HTTPException(status_code=400, detail=f'Kernel size must be greater than 0, got {kernel_size}')
    
    should_norm = kernel_type == "norm"
    random_filtered_image = RandomFilter(image,kernel_size,low,high, normalize=should_norm)
    image_bytes = image_to_bytes(random_filtered_image)

    return StreamingResponse(image_bytes, media_type="image/png")

@app.get("/api/invertedImage/{image_link:path}")
async def invert_image(image_link: str):
    """
    invert_image
    returns bytes in png format of the inverted version of an image
    :param image_link:string : encoded url to an image file
    """
    image = read_HTTP_into_mat(image_link)
    inverted_image = Invert(image)
    image_bytes = image_to_bytes(inverted_image)

    return StreamingResponse(image_bytes, media_type="image/png")


@app.get("/api/saturatedImage/{image_link:path}/{saturation}")
async def saturate_image(image_link: str, saturation: float):
    """
    saturate_image
    returns bytes in png format of the saturated version of an image
    :param image_link:string : encoded url to an image file
    :param saturation:int    : magnitude of saturation
    """
    image = read_HTTP_into_mat(image_link)
    saturated_image = Saturate(image,saturation)
    image_bytes = image_to_bytes(saturated_image)

    return StreamingResponse(image_bytes, media_type="image/png")

@app.get("/api/edgeImage/{image_link:path}/{lower}/{higher}")
async def edge_detect_image(image_link: str, lower: int, higher: int):
    """
    edge_detect_image
    returns bytes in png format of the edge detected version of an image
    :param image_link:string : encoded url to an image file
    :param lower:int         : lower bound for edge values
    :param higher:int        : upper bound for edge values
    """
    image = read_HTTP_into_mat(image_link)

    if lower <= 0 or higher <= 0:
        raise HTTPException(status_code=400, detail=f'lower and higher bounds and iterations must be greater than 0, got {lower} and {higher}')

    edges = EdgeDetect(image,lower,higher)
    image_bytes = image_to_bytes(edges)

    return StreamingResponse(image_bytes, media_type="image/png")




@app.get("/api/dilatedImage/{image_link:path}/{box_size}/{iterations}")
async def dilate_image(image_link:str, box_size: int, iterations: int):
    """
    dilate_image
    returns bytes in png format of the  the dilated version of an image, A⊕ B. For more information on what dilation is, see https://en.wikipedia.org/wiki/Dilation_(morphology)
    :param image_link:string : encoded url to an image file
    :param box_size:int      : width and height of the structuring element (see above link)
    :param iterations:int   : number of dilation operations to perform
    """
    image = read_HTTP_into_mat(image_link)

    if box_size <= 0 or iterations <= 0:
        raise HTTPException(status_code=400, detail=f'box size and iterations must be greater than 0, got {box_size} and {iterations}')

    dilated_image = Dilate(image, kernel_size=box_size, iterations=iterations)
    image_bytes = image_to_bytes(dilated_image)

    return StreamingResponse(image_bytes, media_type="image/png")



@app.get("/api/erodedImage/{image_link:path}/{box_size}/{iterations}")
async def erode_image(image_link:str, box_size: int, iterations: int):
    """
    erode_image
    returns bytes in png format of the eroded version of an image, A⊕ B. For more information on what erosion is, see https://en.wikipedia.org/wiki/Erosion_(morphology)
    :param image_link:string : encoded url to an image file
    :param box_size:int      : width and height of the structuring element (see above link)
    :param iterations:int   : number of erosion operations to perform
    """
    image = read_HTTP_into_mat(image_link)

    if box_size <= 0 or iterations <= 0:
        raise HTTPException(status_code=400, detail=f'box size and iterations must be greater than 0, got {box_size} and {iterations}')

    eroded_image = Erode(image, kernel_size=box_size, iterations=iterations)
    image_bytes = image_to_bytes(eroded_image)

    return StreamingResponse(image_bytes, media_type="image/png")



@app.get("/api/textImage/{image_link:path}/{text}/{font_scale}/{x}/{y}")
async def add_text_to_image(image_link:str, text:str, font_scale:float, x:float, y:float):
    """
    add_text_to_image
    returns bytes in png format of an image with text drawn on
    :param image_link:string : encoded url to an image file
    :param text:string       : text to add
    :font_scale:float        : scale factor for font size
    :x:float                 : percentage of the width, where the x coordinate for the text shall be
    :y:float                 : percentage of the height, where the y coordinate for the text shall be
    """
    image = read_HTTP_into_mat(image_link)
    if x < 0.0 or y < 0.0 or x>1.0 or y >1.0:
        raise HTTPException(status_code=400, detail=f'x and y percentages must be between 0 and 1, got {x} and {y}')
    
    if font_scale <= 0.0:
        raise HTTPException(status_code=400, detail=f'font scale must be greater than 0.0, got {font_scale}')


    reduced_image = Add_text(image,text,font_scale,x,y)
    image_bytes = image_to_bytes(reduced_image)

    return StreamingResponse(image_bytes, media_type="image/png")


 

@app.get("/api/reducedImage/{image_link:path}/{quality}")
async def reduce_image(image_link: str, quality: float):
    """
    reduce_image
    returns bytes in png format of an image with reduced quality
    :param image_link:string : encoded url to an image file
    :param quality:float     : percentage of image quality relative to the original image
    """
    image = read_HTTP_into_mat(image_link)
    if quality <= 0.0 or quality > 1.0:
        raise HTTPException(status_code=400, detail=f'quality level must be between 0 and 1, got {quality}')

    reduced_image = Reduce(image, quality)
    image_bytes = image_to_bytes(reduced_image)

    return StreamingResponse(image_bytes, media_type="image/png")

