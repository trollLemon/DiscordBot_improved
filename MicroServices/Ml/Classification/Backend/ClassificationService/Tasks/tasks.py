import logging
import os
import shutil
from typing import Optional
from PIL import Image

from Backend.ClassificationService import model_config

from Broker.broker import celery
from Backend.ClassificationService.Model import Classifier


model_config = model_config.get_model_and_processor(model_config.ModelType.VITBASE)

model = Classifier(
    processor = model_config['preprocessor'],
    model     = model_config['model'],
    logger    = logging.getLogger(__name__)
)


@celery.task(name="tasks.classification")
def classify(image_path)-> Optional[str]:
    image = Image.open(image_path)
    image = image.convert('RGB')
    directory = os.path.dirname(image_path)
    shutil.rmtree(directory)
    try:
        pred_class = model.predict(image, lambda x: x.argmax(-1).item())
        return pred_class
    except Exception as e:
        raise RuntimeError(f"Image classification failed: {str(e)}") from e