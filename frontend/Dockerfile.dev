FROM node:20-alpine

WORKDIR /my-app

COPY ./task-manager-app ./
# RUN cd /my-app
RUN npm install

# クラウドなどへのデプロイを考慮して、CMD ["npm", "start"] としている
# 本番環境用
RUN npm run build
CMD ["npm", "start"]

#　開発環境はdocker-compose.ymlの設定で起動する。

