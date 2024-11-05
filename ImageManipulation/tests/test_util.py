import sys
import os

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

import util
import unittest
import cv2 as cv
import numpy as np


def imageEqual(image1, image2):
    return image1.shape == image2.shape and not(np.bitwise_xor(image1,image2).any())




TEST_IMG_LINK_JPEG = "https://samplelib.com/lib/preview/jpeg/sample-clouds-400x300.jpg"
TEST_IMG_LINK_PNG = "https://samplelib.com/lib/preview/png/sample-boat-400x300.png"
TEST_GIF_LINK = "https://samplelib.com/lib/preview/gif/sample-animated-400x300.gif"

class TestHTTPToImage(unittest.TestCase):
    def testJpeg(self): 
        image_from_http = util.read_HTTP_into_mat(TEST_IMG_LINK_JPEG)
        actual_image = cv.imread("./test.jpg")
        expected = True
        actual = imageEqual(image_from_http, actual_image)
        self.assertEqual(expected,actual)
     
    def testPng(self):                
        image_from_http = util.read_HTTP_into_mat(TEST_IMG_LINK_PNG)
        actual_image = cv.imread("./test.png")
        expected = True
        actual = imageEqual(image_from_http, actual_image)
        self.assertEqual(expected,actual)

    def testGif(self):
        self.assertRaises(ValueError, lambda: util.read_HTTP_into_mat(TEST_GIF_LINK))
    def test_invalid_url(self):
        self.assertRaises(ValueError, lambda: util.read_HTTP_into_mat("invalid link"))
if __name__ == '__main__':
    unittest.main()
