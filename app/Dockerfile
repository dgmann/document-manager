FROM node:latest as builder

WORKDIR /build
ADD package*.json ./
RUN npm install

COPY . .
RUN npm run build


FROM nginx:alpine

COPY --from=builder /build/dist /usr/share/nginx/html/
