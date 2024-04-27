FROM node:20-alpine as builder
ARG SERVICE
WORKDIR /build
RUN --mount=type=cache,target=/root/.npm \
    --mount=type=bind,source=apps/${SERVICE}/package.json,target=package.json \
    --mount=type=bind,source=apps/${SERVICE}/package-lock.json,target=package-lock.json \
    npm ci

COPY apps/${SERVICE} .
RUN npm run build


FROM nginx:alpine
ARG SERVICE
COPY --link apps/${SERVICE}/nginx.conf /etc/nginx/nginx.conf
COPY --link --chmod=+x apps/${SERVICE}/run.sh /run.sh
COPY --link --from=builder /build/dist/app /app

ENV ENABLE_SSL=false

CMD ["/bin/sh", "/run.sh"]
