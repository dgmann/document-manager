FROM node:lts-alpine as builder

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

ENV API_HOST=localhost ENABLE_SSL=true

CMD ["/bin/sh", "/run.sh"]
