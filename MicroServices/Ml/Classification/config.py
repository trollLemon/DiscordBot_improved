from enum import Enum

from transformers import ViTImageProcessor, ViTForImageClassification


class ModelType(Enum):
    VITBASE = "vit-base"

def get_model_and_processor(model_type : ModelType):

    if model_type == ModelType.VITBASE:
        return {
            "model":  ViTForImageClassification.from_pretrained('google/vit-base-patch16-224'),
            "preprocessor": ViTImageProcessor.from_pretrained('google/vit-base-patch16-224'),
        }

    else:
        raise Exception("Unknown model type")
