FROM node:18.19.0-alpine3.18 as build


WORKDIR /app
COPY ../ant-react ./ant-react

RUN cd ant-react && yarn install --registry=https://registry.npmmirror.com && yarn run build



FROM nginx:stable-alpine3.17


COPY ../ant-react/nginx.conf /etc/nginx/nginx.conf
WORKDIR /app
COPY --from=build /app/ant-react/dist /app/iot/project/html

RUN mkdir /var/log/nginx/iot

CMD ["nginx", "-g", "daemon off;"]
