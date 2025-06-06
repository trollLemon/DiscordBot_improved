import os
from celery import Celery

broker_url = os.getenv('CELERY_BROKER_URL', 'redis://worker-data:6379/0')
result_backend = os.getenv('CELERY_RESULT_BACKEND', 'redis://worker-data:6379/1')

app = Celery('Tasks', broker=broker_url, backend=result_backend)

app.autodiscover_tasks(['Tasks.tasks'])