FROM python:3.9-slim
# TODO пример окружения

RUN pip install --no-cache-dir flask

WORKDIR /app
COPY run_code.py .

CMD ["python", "run_code.py"]
