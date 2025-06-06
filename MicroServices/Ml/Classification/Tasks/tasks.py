import logging
import os
import shutil
from typing import Optional
from PIL import Image

import config
from Tasks.task_queue import app
from Model import Classifier

model_config = config.get_model_and_processor(config.ModelType.VITBASE)


model = Classifier(
    processor = model_config['preprocessor'],
    model     = model_config['model'],
    logger    = logging.getLogger(__name__)
)


@app.task()
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