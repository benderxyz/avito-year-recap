FROM node:20-alpine

WORKDIR /app

COPY . .

RUN if [ -f package.json ]; then npm install; fi

EXPOSE 3000

CMD ["sh", "-c", "if [ -f package.json ]; then npm run start; else npx http-server -p 3000 .; fi"]
