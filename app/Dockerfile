FROM node:latest as builder

WORKDIR /build
ADD package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build


FROM nginx:alpine

COPY .docker/site.conf /etc/nginx/conf.d/default.conf
COPY .docker/nginx.conf /etc/nginx/nginx.conf
COPY .docker/run.sh /run.sh
RUN chmod +x /run.sh
COPY --from=builder /build/dist /usr/share/nginx/html/

ENV API_URL=http://localhost:8080 WS_URL=ws://localhost:8080/notifications BUGSNAG_KEY=notset

CMD ["/bin/sh", "/run.sh"]
