import torch

from enum import Enum

from transformers import ViTImageProcessor, ViTForImageClassification


class ModelType(Enum):
    VITBASE = "vit-base"


def get_model_and_processor(model_type: ModelType):
    if model_type == ModelType.VITBASE:
        return {
            "model": ViTForImageClassification.from_pretrained('google/vit-base-patch16-224'),
            "preprocessor": ViTImageProcessor.from_pretrained('google/vit-base-patch16-224'),
        }

    raise Exception("Unknown model type")


class Classifier:
    def __init__(self, **kwargs):
        self.processor = kwargs['processor']
        self.model = kwargs['model']
        self.logger = kwargs['logger']

    def predict(self, image, decider):
        if not self.processor or not self.model:
            self.logger.error("No processor or model")
            raise ValueError('processor and model must be not None')

        self.logger.info("Starting Prediction")
        with torch.no_grad():
            inputs = self.processor(images=image, return_tensors="pt")
            output = self.model(**inputs)
            logits = output.logits
            props = logits.softmax(-1)

        self.logger.info("Finished Prediction")
        predicted_class = decider(props)

        return self.model.config.id2label[predicted_class]
