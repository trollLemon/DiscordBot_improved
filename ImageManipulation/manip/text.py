import cv2 as cv
import numpy as np

def Add_text(image : np.ndarray, text, font_scale:float, x_perc:float ,y_perc:float) -> np.ndarray:
    
    font = cv.FONT_HERSHEY_COMPLEX
    foreground = (255,255,255)
    outline    = (0,0,0)
    rows, cols,_ = image.shape
    x, y = rows * x_perc, cols * y_perc
    print(x,y)
    foreground_thickness = 2 
    outline_thickness = 4

    image_with_text = cv.putText(image,text,(int(x),int(y)),font,font_scale, outline, outline_thickness, cv.LINE_AA, False)
    image_with_text = cv.putText(image_with_text,text,(int(x),int(y)),font,font_scale, foreground, foreground_thickness, cv.LINE_AA, False)

    return image_with_text






