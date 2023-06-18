FROM node:14-alpine as builder

WORKDIR /build
ADD package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build


FROM nginx:alpine

COPY nginx.conf /etc/nginx/nginx.conf
COPY run.sh /run.sh
RUN chmod +x /run.sh
COPY --from=builder /build/dist/app /app

ENV ENABLE_SSL=false

CMD ["/bin/sh", "/run.sh"]
