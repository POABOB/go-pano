FROM python:3.9.15-slim-bullseye
LABEL mantainer="POABOB <zxc752166@gmail.com>"

RUN apt-get -qq update \
    && apt-get -qq install -y --no-install-recommends \
    build-essential \
    git \
    wget \
    unzip \
    pkg-config \
    libswscale-dev \
    ffmpeg \
    libsm6 \
    libxext6 \
    libtbb2 \
    libtbb-dev \
    libjpeg-dev \
    libpng-dev \
    libtiff-dev \
    libopenjp2-7-dev \
    libavformat-dev \
    libpq-dev \
    python3-pip\
    python3-opencv \
    && rm -rf /opt/build/* \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get -qq autoremove \
    && apt-get -qq clean


RUN python3 -m pip install opencv-python-headless  

# requirements
COPY . /python/
RUN python3 -m pip install -r /python/requirements.txt --quiet
# RUN touch /usr/local/lib/python3.9/dist-packages/yolov5/requirements.txt


# 暴露端口
# PYTHON HTTP || GRPC
WORKDIR /python/
EXPOSE 5000
CMD ["python3", "server_http.py"]