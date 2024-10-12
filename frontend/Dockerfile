FROM node:20-slim AS builder

WORKDIR /app

COPY . .

RUN npm i

RUN npm run build

FROM nginx:alpine

COPY --from=builder /app/build  /usr/share/nginx/html

COPY conf/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 3000

CMD ["nginx", "-g", "daemon off;"]