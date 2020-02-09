# base image
FROM node:12.2.0-alpine

# set working directory
WORKDIR /app

ENV PATH /app/node_modules/.bin:$PATH

COPY package*.json ./
COPY vue.config.js ./

# install project dependencies
RUN npm install

# start app
CMD ["npm", "run", "serve"]
