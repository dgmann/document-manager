 # syntax=docker/dockerfile:1.10
FROM node:20-alpine AS builder
ARG SERVICE
WORKDIR /build
RUN --mount=type=bind,source=${SERVICE}/package.json,target=package.json \
    --mount=type=bind,source=${SERVICE}/package-lock.json,target=package-lock.json \
    npm ci

RUN --mount=type=bind,source=${SERVICE}/package.json,target=package.json \
    --mount=type=bind,source=${SERVICE}/package-lock.json,target=package-lock.json \
    --mount=type=bind,source=${SERVICE}/src,target=src \
    --mount=type=bind,source=${SERVICE}/angular.json,target=angular.json \
    --mount=type=bind,source=${SERVICE}/tsconfig.json,target=tsconfig.json \
    --mount=type=bind,source=${SERVICE}/tsconfig.app.json,target=tsconfig.app.json \
    npm run build


FROM nginx:alpine
ARG SERVICE
COPY --link ${SERVICE}/nginx.conf /etc/nginx/nginx.conf
COPY --link --chmod=700 ${SERVICE}/run.sh /run.sh
COPY --link --from=builder /build/dist/app /app

ENV ENABLE_SSL=false

CMD ["/bin/sh", "/run.sh"]
