FROM quay.io/0voice/node:lts-alpine AS frontend
RUN npm install pnpm -g
COPY . /app
WORKDIR /app
RUN pnpm bootstrap
RUN pnpm run build-only

FROM quay.io/0voice/nginx:1.25.4 as web
COPY --from=frontend /app/dist/ /usr/share/nginx/html/
# 指定入口程序
ENTRYPOINT ["/docker-entrypoint.sh"]
# 指定容器的启动命令或者入口程序的参数
CMD ["nginx", "-g", "daemon off;"]

