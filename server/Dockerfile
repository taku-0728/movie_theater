# ベースとなるDockerイメージ指定
FROM golang:latest

# コンテナログイン時のディレクトリ指定
WORKDIR /server/work/

# モジュールインポート
RUN go mod init scraping

RUN go get -u github.com/PuerkitoBio/goquery \
    && go get -u github.com/djimenez/iconv-go \
    && go get -u github.com/labstack/echo \
    && go get -u github.com/labstack/echo/middleware \
    && go get github.com/sclevine/agouti

RUN apt-get update && apt-get install -y unzip

#Chrome
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add && \
wget http://dl.google.com/linux/deb/pool/main/g/google-chrome-unstable/google-chrome-unstable_97.0.4692.20-1_amd64.deb && \
apt-get install -y -f ./google-chrome-unstable_97.0.4692.20-1_amd64.deb

# ChromeDriver
ADD https://chromedriver.storage.googleapis.com/97.0.4692.36/chromedriver_linux64.zip /opt/chrome/
RUN cd /opt/chrome/ && \
unzip chromedriver_linux64.zip

ENV PATH /opt/chrome:$PATH

# ホストのファイルをコンテナの作業ディレクトリに移行
ADD ./src /server/work