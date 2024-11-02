FROM node:18.19.0-alpine3.18 as build


WORKDIR /app
COPY ../ant-vue ./ant-vue

RUN cd ant-vue && npm install --registry=https://registry.npmmirror.com && npm run build-docker



FROM nginx:stable-alpine3.17


COPY ../ant-vue/nginx.conf /etc/nginx/nginx.conf
WORKDIR /app
COPY --from=build /app/ant-vue/dist /app/iot/project/html

RUN mkdir /var/log/nginx/iot

CMD ["nginx", "-g", "daemon off;"]
