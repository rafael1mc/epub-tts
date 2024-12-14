FROM node:latest

WORKDIR /app

RUN npm install @gxl/epub-parser --save

COPY main.js main.js

ENV NODE_NO_WARNINGS=1

ENTRYPOINT ["node", "main.js"]