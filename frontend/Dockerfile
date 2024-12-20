# Build stage
FROM node:20-alpine AS builder
WORKDIR /app
COPY ./task-manager-app .
RUN npm install
RUN npm run build

# Expose port 3000
EXPOSE 3000

# Start server
CMD ["npm", "start"]