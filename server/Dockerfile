# ベースとなるDockerイメージ指定
FROM golang:latest

# コンテナログイン時のディレクトリ指定
WORKDIR /server/work/

COPY ./src/go.mod /server/work
COPY ./src/go.sum /server/work

RUN go mod download
COPY ./src /server/work

RUN apt-get update && apt-get install -y unzip

#Chrome
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add && \
wget http://dl.google.com/linux/deb/pool/main/g/google-chrome-unstable/google-chrome-unstable_97.0.4692.20-1_amd64.deb && \
apt-get install -y -f ./google-chrome-unstable_97.0.4692.20-1_amd64.deb

# ChromeDriver
ADD https://chromedriver.storage.googleapis.com/97.0.4692.20/chromedriver_linux64.zip /opt/chrome/
RUN cd /opt/chrome/ && \
unzip chromedriver_linux64.zip

ENV PATH /opt/chrome:$PATH
