# ベースとなるDockerイメージ指定
FROM golang:latest

# コンテナログイン時のディレクトリ指定
WORKDIR /server/work/

# モジュールインポート
RUN go mod init goquery

RUN go get -u github.com/PuerkitoBio/goquery \
    && go get -u github.com/djimenez/iconv-go

# ホストのファイルをコンテナの作業ディレクトリに移行
ADD ./src /server/work