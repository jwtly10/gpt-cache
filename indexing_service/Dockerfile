from python:3.12

WORKDIR /app

COPY . /app/

RUN pip install --no-cache-dir -r requirements.txt

EXPOSE 8000

ENV MODULE_NAME=app.main
ENV VARIABLE_NAME=app

CMD ["sh", "-c", "uvicorn $MODULE_NAME:$VARIABLE_NAME --host 0.0.0.0 --port 8000"]
