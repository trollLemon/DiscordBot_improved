import unittest

import numpy as np
from Model.model import Classifier
from PIL import Image

class TestModelOutput:
    def __init__(self, logits):
        self.logits = logits

class TestModelConfig:
    def __init__(self, labels):
        self.id2label = labels

class FakeModel:
    def __init__(self, config, test_output):
        self.config = config
        self.test_output = test_output
        pass
    def __call__(self, *args, **kwargs):
        return TestModelOutput(self.test_output)

class FakeProcessor:
    def __init__(self):
        pass
    def __call__(self, *args, **kwargs):
        return {'input': 0.0}



class TestModel(unittest.TestCase):
    def test_invalid_processor(self):

        test_params = {
            'processor' : None,
            'model': FakeModel(None),

        }

        test_image = Image.new('RGB', (100, 100))

        test_model = Classifier(**test_params)

        with self.assertRaises(ValueError):
              test_model.predict(test_image, lambda x: np.argmax(x))

    def test_invalid_Model(self):
        test_params = {
            'processor': FakeProcessor(),
            'model': None ,

        }

        test_image = Image.new('RGB', (100, 100))

        test_model = Classifier(**test_params)
        with self.assertRaises(ValueError):
            test_model.predict(test_image, lambda x: np.argmax(x))

    def test_model(self):

        output_logits = [0.0,0.1,0.2,0.7] # max index is 3

        labels = ['a', 'b', 'c','d']
        test_params = {
            'processor': FakeProcessor(),
            'model': FakeModel( config= TestModelConfig(labels), test_output=output_logits),

        }

        test_image = Image.new('RGB', (100, 100))

        test_model = Classifier(**test_params)

        pred_class = test_model.predict(test_image, lambda x: np.argmax(x))

        self.assertEqual(pred_class, labels[np.argmax(output_logits)])



if __name__ == '__main__':
    unittest.main()
