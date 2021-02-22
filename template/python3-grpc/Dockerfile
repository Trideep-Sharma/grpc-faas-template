FROM python:3.7-slim

RUN pip install --upgrade pip
RUN pip install grpcio
RUN pip install google-cloud-storage

WORKDIR /home/app/

COPY index.py           .
COPY requirements.txt   .
RUN pip install -r requirements.txt

RUN mkdir -p function
WORKDIR /home/app/function/
COPY function/requirements.txt	.
RUN pip install -r requirements.txt

WORKDIR /home/app/
COPY function   function

ENV fprocess="python index.py"
CMD ["python", "index.py"]
