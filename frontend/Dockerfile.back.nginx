# Build stage
FROM node:20-alpine AS builder
WORKDIR /app
COPY ./task-manager-app .
RUN npm install
RUN npm run build

# Serve stage
FROM nginx:alpine

# NGINX設定ファイルをコンテナにコピー
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Remove default nginx static resources and copy from the builder
RUN rm -rf /usr/share/nginx/html/*

# Copy built Next.js application from the builder stage
COPY --from=builder /app/out /usr/share/nginx/html

# Expose port 3000
EXPOSE 3000

# Start Nginx server
CMD ["nginx", "-g", "daemon off;"]