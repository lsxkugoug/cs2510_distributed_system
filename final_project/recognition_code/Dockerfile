# syntax=docker/dockerfile:1
FROM gcp_img_recignition

COPY . /final_img_recognition

WORKDIR /final_img_recognition
EXPOSE 9999

ENTRYPOINT ["python3", "server.py" ]



