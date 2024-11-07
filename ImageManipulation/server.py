from manip import Invert, Saturate, edgeDetect, Dilate, Erode, reduce, add_text
from util import read_HTTP_into_mat, image_to_bytes
from fastapi import FastAPI, HTTPException

from fastapi.responses import StreamingResponse
import shutil


from fastapi import FastAPI

app = FastAPI()




@app.get("/api/invertedImage/{image_link:path}")
async def invert_image(image_link: str):
    image = read_HTTP_into_mat(image_link)
    
    inverted_image = Invert(image)
    image_bytes = image_to_bytes(inverted_image)

    return StreamingResponse(image_bytes, media_type="image/png")


@app.get("/api/saturatedImage/{image_link:path}/{saturation}")
async def saturate_image(image_link: str, saturation: float):
    image = read_HTTP_into_mat(image_link)
    
    saturated_image = Saturate(image,saturation)
    image_bytes = image_to_bytes(saturated_image)

    return StreamingResponse(image_bytes, media_type="image/png")

@app.get("/api/edgeImage/{image_link:path}/{lower}/{higher}")
async def edge_detect_image(image_link: str, lower: int, higher: int):
    image = read_HTTP_into_mat(image_link)
    
    if lower <= 0 or higher <= 0:
        raise HTTPException(status_code=400, detail=f'lower and higher bounds and iterations must be greater than 0, got {lower} and {higher}') 
    
    edges = edgeDetect(image,lower,higher)
    image_bytes = image_to_bytes(edges)

    return StreamingResponse(image_bytes, media_type="image/png")



@app.get("/api/dilatedImage/{image_link:path}/{box_size}/{iterations}")
async def dilate_image(image_link:str, box_size: int, iterations: int):
    
    image = read_HTTP_into_mat(image_link)
    
    if box_size <= 0 or iterations <= 0:
        raise HTTPException(status_code=400, detail=f'box size and iterations must be greater than 0, got {box_size} and {iterations}') 
    
    dilated_image = Dilate(image, kernel_size=box_size, iterations=iterations)
    image_bytes = image_to_bytes(dilated_image)

    return StreamingResponse(image_bytes, media_type="image/png")

@app.get("/api/erodedImage/{image_link:path}/{box_size}/{iterations}")
async def erode_image(image_link:str, box_size: int, iterations: int):

    image = read_HTTP_into_mat(image_link)
    
    if box_size <= 0 or iterations <= 0:
        raise HTTPException(status_code=400, detail=f'box size and iterations must be greater than 0, got {box_size} and {iterations}') 
    
    eroded_image = Erode(image, kernel_size=box_size, iterations=iterations)
    image_bytes = image_to_bytes(eroded_image)

    return StreamingResponse(image_bytes, media_type="image/png")



@app.get("/api/textImage/{image_link:path}/{text}/{location}")
async def add_text_to_image(image_link:str, text:str, location:str):
    pass

@app.get("/api/reducedImage/{image_link:path}/{quality}")
async def reduce_image(image_link: str, quality: float):
    
    image = read_HTTP_into_mat(image_link)
    if quality <= 0.0 or quality > 1.0:
        raise HTTPException(status_code=400, detail=f'quality level must be between 0 and 1, got {quality}')

    reduced_image = reduce(image, quality)
    image_bytes = image_to_bytes(reduced_image)

    return StreamingResponse(image_bytes, media_type="image/png")





@app.get("/")
def read_root():
    return {"Hello": "World"}
