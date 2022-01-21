FROM python:3.7-slim
ADD calibration.py /
ADD client.py /
ENV WEBSOCKET ws://localhost:8888/ws/calibration
COPY ./requirements.txt /var/www/requirements.txt
RUN pip install -r /var/www/requirements.txt
CMD [ "python", "./client.py" ]