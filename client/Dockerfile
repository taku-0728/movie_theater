#nodeイメージをpullする
FROM node:latest

#working directory
WORKDIR /client/work

#vuecliインストール
RUN npm install -g @vue/cli

# ホストのファイルをコンテナの作業ディレクトリに移行
ADD ./ /client/work