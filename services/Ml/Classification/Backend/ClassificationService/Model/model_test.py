import unittest
import logging

import numpy as np
import torch

from ClassificationService.Model import Classifier
from PIL import Image


class FakeModelOutput:
    def __init__(self, logits):
        self.logits = logits

class FakeModelConfig:
    def __init__(self, test_labels):
        self.id2label = test_labels

class FakeModel:
    def __init__(self, config, test_output):
        self.config = config
        self.test_output = test_output
        pass
    def __call__(self, *args, **kwargs):
        return FakeModelOutput(self.test_output)

class FakeProcessor:
    def __init__(self):
        pass
    def __call__(self, *args, **kwargs):
        return {'input': 0.0}

output_logits = torch.tensor([0.0,0.1,0.2,0.7]) # max index is 3
labels = ['a', 'b', 'c','d']

class TestModel(unittest.TestCase):
    def test_invalid_processor(self):

        test_params = {
            'processor' : None,
            'model':  FakeModel(config= FakeModelConfig(labels), test_output=output_logits),
            'logger': logging.getLogger(__name__)

        }

        test_image = Image.new('RGB', (100, 100))

        test_model = Classifier(**test_params)

        with self.assertRaises(ValueError):
              test_model.predict(test_image, lambda x: np.argmax(x))

    def test_invalid_Model(self):
        test_params = {
            'processor': FakeProcessor(),
            'model': None ,
            'logger': logging.getLogger(__name__)

        }

        test_image = Image.new('RGB', (100, 100))

        test_model = Classifier(**test_params)
        with self.assertRaises(ValueError):
            test_model.predict(test_image, lambda x: np.argmax(x))

    def test_model(self):


        test_params = {
            'processor': FakeProcessor(),
            'model': FakeModel(config= FakeModelConfig(labels), test_output=output_logits),
            'logger': logging.getLogger(__name__)

        }

        test_image = Image.new('RGB', (100, 100))

        test_model = Classifier(**test_params)

        pred_class = test_model.predict(test_image, lambda x: np.argmax(x))

        self.assertEqual(pred_class, labels[np.argmax(output_logits)])



if __name__ == '__main__':
    unittest.main()
