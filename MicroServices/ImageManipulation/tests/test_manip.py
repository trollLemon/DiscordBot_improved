"""
test_manip.py

This file tests the image manipulation functions in the Manip module.
Tests mostly cover incorrect input paremeters, which should throw a value error.

We test with both greyscale and rgb images.
"""

import sys
import os

sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from manip import *
import unittest
import numpy as np


class TestColors(unittest.TestCase):
   
    def setUp(self):
        # initialize test images
        self.arr = np.random.randint(low=0, high=255, size=(64, 64, 3), dtype=np.uint8)
        self.arr_bigger = np.random.randint(low=0, high=255, size=(64, 64, 3), dtype=np.uint8)
        self.arr_g = np.random.randint(low=0, high=255, size=(64, 64, 1), dtype=np.uint8)
        self.arr_bigger_g = np.random.randint(low=0, high=255, size=(64, 64, 1), dtype=np.uint8)

    def test_invert(self):
        
        # these should run without error
        Invert(self.arr)
        Invert(self.arr_bigger)

    def test_invert_greyscale(self):
        Invert(self.arr_g)
        Invert(self.arr_bigger_g)
    
    def test_saturate(self):
        bad_value = 0.0
        another_bad_value = -2.0
        good_value = 4.0
        good_value_other = 40.0
        
        Saturate(self.arr,good_value)
        Saturate(self.arr,good_value_other)
        Saturate(self.arr_bigger,good_value)
        Saturate(self.arr_bigger,good_value_other)
        
        with self.assertRaises(ValueError):
            # should throw value errors due to bad parameters
            Saturate(self.arr,bad_value)
            Saturate(self.arr,another_bad_value)
            Saturate(self.arr_bigger,bad_value)
            Saturate(self.arr_bigger,another_bad_value)

    def test_saturate_greyscale(self):
        bad_value = 0.0
        another_bad_value = -2.0
        good_value = 4.0
        good_value_other = 40.0
        
       
        with self.assertRaises(ValueError):
            # should throw value errors due to bad parameters
            Saturate(self.arr_g,good_value)
            Saturate(self.arr_g,good_value_other)
            Saturate(self.arr_bigger_g,good_value)
            Saturate(self.arr_bigger_g,good_value_other)
            Saturate(self.arr_g,bad_value)
            Saturate(self.arr_g,another_bad_value)
            Saturate(self.arr_bigger_g,bad_value)
            Saturate(self.arr_bigger_g,another_bad_value)

class TestEdges(unittest.TestCase):
        
    def setUp(self):
        # initialize test images
        self.arr = np.random.randint(low=0, high=255, size=(64, 64, 3), dtype=np.uint8)
        self.arr_bigger = np.random.randint(low=0, high=255, size=(64, 64, 3), dtype=np.uint8)
        self.arr_g = np.random.randint(low=0, high=255, size=(64, 64, 1), dtype=np.uint8)
        self.arr_bigger_g = np.random.randint(low=0, high=255, size=(64, 64, 1), dtype=np.uint8)


    def test_edge_detection(self):
        bad_value = 0
        another_bad_value = -1
        good_lower = 100
        good_upper = 200
        

        EdgeDetect(self.arr, good_lower, good_upper)
        EdgeDetect(self.arr_bigger, good_lower, good_upper)
        EdgeDetect(self.arr_g, good_lower, good_upper)
        EdgeDetect(self.arr_bigger_g, good_lower, good_upper)
    
        with self.assertRaises(ValueError):
            # if any value passed as the lower or upperbound is <= 0, we should get a value error
            EdgeDetect(self.arr, bad_value, good_upper)
            EdgeDetect(self.arr_bigger, another_bad_value, good_upper)
            EdgeDetect(self.arr_g, bad_value, good_upper)
            EdgeDetect(self.arr_bigger_g, another_bad_value, good_upper)
 
            EdgeDetect(self.arr, good_lower, bad_value)
            EdgeDetect(self.arr_bigger, good_lower, another_bad_value)
            EdgeDetect(self.arr_g, good_lower, bad_value)
            EdgeDetect(self.arr_bigger_g, good_lower, another_bad_value)
   


class TestMorphology(unittest.TestCase):
        
    def setUp(self):
        # initialize test images
        self.arr = np.random.randint(low=0, high=255, size=(64, 64, 3), dtype=np.uint8)
        self.arr_bigger = np.random.randint(low=0, high=255, size=(64, 64, 3), dtype=np.uint8)
        self.arr_g = np.random.randint(low=0, high=255, size=(64, 64, 1), dtype=np.uint8)
        self.arr_bigger_g = np.random.randint(low=0, high=255, size=(64, 64, 1), dtype=np.uint8)


    def test_dilate(self):
        bad_value = 0
        good_itr = 4
        good_size = 3

        Dilate(self.arr, good_itr, good_size)
        Dilate(self.arr_bigger, good_itr, good_size)
        Dilate(self.arr_g, good_itr, good_size)
        Dilate(self.arr_bigger_g, good_itr, good_size)
        
        with self.assertRaises(ValueError):
            Dilate(self.arr, good_itr, bad_value)
            Dilate(self.arr_bigger, bad_value, good_size)
            Dilate(self.arr_g, good_itr, bad_value)
            Dilate(self.arr_bigger_g, bad_value, good_size)

    def test_erode(self):
        bad_value = 0
        good_itr = 4
        good_size = 3

        Erode(self.arr, good_itr, good_size)
        Erode(self.arr_bigger, good_itr, good_size)
        Erode(self.arr_g, good_itr, good_size)
        Erode(self.arr_bigger_g, good_itr, good_size)

        with self.assertRaises(ValueError):
            Erode(self.arr, bad_value, good_size)
            Erode(self.arr_bigger, bad_value, good_size)
            Erode(self.arr_g, bad_value, good_size)
            Erode(self.arr_bigger_g, bad_value, good_size)

class TestQuality(unittest.TestCase):
 
    def setUp(self):
        # initialize test images
        self.arr = np.random.randint(low=0, high=255, size=(64, 64, 3), dtype=np.uint8)
        self.arr_bigger = np.random.randint(low=0, high=255, size=(64, 64, 3), dtype=np.uint8)
        self.arr_g = np.random.randint(low=0, high=255, size=(64, 64, 1), dtype=np.uint8)
        self.arr_bigger_g = np.random.randint(low=0, high=255, size=(64, 64, 1), dtype=np.uint8)


    def test_reduce(self):
        bad_value = 2.0
        good_value = 0.6
        other_good_value = 0.01
        other_bad_value = 0.0
        

        Reduce(self.arr, good_value)
        Reduce(self.arr_bigger, good_value)
        Reduce(self.arr, other_good_value)
        Reduce(self.arr_bigger, other_good_value)

        with self.assertRaises(ValueError): 
            Reduce(self.arr, bad_value)
            Reduce(self.arr_bigger, bad_value)
            Reduce(self.arr, other_bad_value)
            Reduce(self.arr_bigger, other_bad_value)

class TestRandom(unittest.TestCase):
    
    def setUp(self):
        # initialize test images
        self.arr = np.random.randint(low=0, high=255, size=(64, 64, 3), dtype=np.uint8)
        self.arr_bigger = np.random.randint(low=0, high=255, size=(64, 64, 3), dtype=np.uint8)
        self.arr_g = np.random.randint(low=0, high=255, size=(64, 64, 1), dtype=np.uint8)
        self.arr_bigger_g = np.random.randint(low=0, high=255, size=(64, 64, 1), dtype=np.uint8)


    def test_random_filter(self):
        with self.assertRaises(ValueError):
            # we can't have a kernel size less than 1
            RandomFilter(image=self.arr, kernel_size=-1, min=0, max=1, normalize=False)
        
        RandomFilter(image=self.arr, kernel_size=3, min=0, max=1, normalize=False)
        RandomFilter(image=self.arr_bigger, kernel_size=4, min=0, max=1, normalize=False)
        
        with self.assertRaises(ValueError):
            # we cant use this random filter method with greyscale images, as its supposed to work on RGB images
            RandomFilter(image=self.arr_g, kernel_size=3, min=0, max=1, normalize=False)
            RandomFilter(image=self.arr_bigger_g, kernel_size=9, min=0, max=1, normalize=False)        



class TestText(unittest.TestCase):
    def setUp(self):
        # initialize test images
        self.arr = np.random.randint(low=0, high=255, size=(64, 64, 3), dtype=np.uint8)
        self.arr_bigger = np.random.randint(low=0, high=255, size=(64, 64, 3), dtype=np.uint8)
        self.arr_g = np.random.randint(low=0, high=255, size=(64, 64, 1), dtype=np.uint8)
        self.arr_bigger_g = np.random.randint(low=0, high=255, size=(64, 64, 1), dtype=np.uint8)

    def test_add_text(self):
            # these should run without error
            AddText(image=self.arr, text="Test", font_scale=0.5, x_perc=0.2, y_perc=0.3)
            AddText(image=self.arr_bigger, text="Test", font_scale=0.5, x_perc=0.2, y_perc=0.3)
            AddText(image=self.arr_g, text="Test", font_scale=0.5, x_perc=0.2, y_perc=0.3)
            AddText(image=self.arr_bigger_g, text="Test", font_scale=0.5, x_perc=0.2, y_perc=0.3)
            
            with self.assertRaises(ValueError):
                # these should error due to percentages or font scale cannot be negative
                AddText(image=self.arr, text="Test", font_scale=-0.5, x_perc=0.2, y_perc=0.3)
                AddText(image=self.arr_bigger, text="Test", font_scale=0.5, x_perc=-0.2, y_perc=0.3)
                AddText(image=self.arr_g, text="Test", font_scale=0.5, x_perc=0.2, y_perc=-0.3)
                AddText(image=self.arr_bigger_g, text="Test", font_scale=0.5, x_perc=-0.2, y_perc=0.3)
            


class TestMisc(unittest.TestCase):

    def setUp(self):
        # initialize test images
        self.arr = np.random.randint(low=0, high=255, size=(64, 64, 3), dtype=np.uint8)
        self.arr_bigger = np.random.randint(low=0, high=255, size=(64, 64, 3), dtype=np.uint8)
        self.arr_g = np.random.randint(low=0, high=255, size=(64, 64, 1), dtype=np.uint8)
        self.arr_bigger_g = np.random.randint(low=0, high=255, size=(64, 64, 1), dtype=np.uint8)
    
    def test_shuffle(self):
        
        bad_value = 1
        other_bad_value = 0
        good_value = 4
        other_good_value = 320
        too_large = 4096

        Shuffle(self.arr,good_value)
        Shuffle(self.arr_bigger,good_value)
        Shuffle(self.arr_g,good_value)
        Shuffle(self.arr_bigger_g,other_good_value)
        
        with self.assertRaises(ValueError):
            # we cant shuffle 0 partitions, and 1 partition is the entire image
            Shuffle(self.arr,other_bad_value)
            Shuffle(self.arr_bigger,bad_value)
            Shuffle(self.arr_g,other_bad_value)
            Shuffle(self.arr_bigger_g,other_bad_value)
            
            # in this case, the number of partitions is >= the area of the image, which we cannot shuffle
            Shuffle(self.arr, too_large)
            Shuffle(self.arr_g, too_large)
        












