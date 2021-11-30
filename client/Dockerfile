#nodeイメージをpullする
FROM node:14.18.1

#working directory
WORKDIR /work

# カレントディレクトリのファイル全てをコンテナ内にコピー
COPY . .

# 依存モジュールのインストール
RUN yarn install