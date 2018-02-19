FROM python:2.7-alpine

RUN mkdir /app/
WORKDIR /app/
COPY . /app/

RUN pip install -r requirements.txt

EXPOSE 8888

CMD [ "python", "app.py" ]
