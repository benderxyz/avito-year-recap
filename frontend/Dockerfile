FROM node:20-alpine

WORKDIR /app/frontend

COPY frontend/package.json ./
RUN npm install

COPY frontend/index.html frontend/vite.config.ts frontend/tsconfig.json frontend/tsconfig.node.json ./
COPY frontend/public ./public
COPY frontend/src ./src
COPY recap-engine/packages /app/recap-engine/packages

EXPOSE 3000

CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
