FROM node:lts-alpine as builder

WORKDIR /build
ADD package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build-low-memory


FROM nginx:alpine

COPY nginx.conf /etc/nginx/nginx.conf
COPY run.sh /run.sh
RUN chmod +x /run.sh
COPY --from=builder /build/dist/app /app

ENV API_URL=http://localhost:8080 WS_URL=ws://localhost:8080/notifications

CMD ["/bin/sh", "/run.sh"]
