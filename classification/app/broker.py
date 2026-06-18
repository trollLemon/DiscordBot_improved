import os

from celery import Celery

broker_url = os.getenv('CELERY_BROKER_URL', 'redis://worker-data:6379/0')
result_backend = os.getenv('CELERY_RESULT_BACKEND', 'redis://worker-data:6379/1')

celery = Celery('Tasks', broker=broker_url, backend=result_backend)

celery.conf.result_expires = int(os.getenv('CELERY_RESULT_EXPIRES', '3600'))

celery.autodiscover_tasks(['app'])
