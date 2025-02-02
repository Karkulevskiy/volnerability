FROM python:3.9-slim
# TODO пример окружения

RUN pip install --no-cache-dir flask

CMD ["sleep", "infinity"]