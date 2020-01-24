FROM node:10.15.0-alpine as base

# create a parent directory for the app
WORKDIR /client

# copy the package.json and lock files
COPY package*.json ./

# install production dependencies & clean up after
RUN npm ci \ 
    && npm cache clean --force

FROM base as dev
ENV NODE_ENV=development
ENV PATH /client/node_modules/.bin:$PATH

# expose default create-react-app port
EXPOSE 3000

# make app directory
RUN mkdir /client/app

RUN npm i --only=development \
    && npm cache clean --force

# patch create-react-app bug preventing self-signed certificate usage
# https://github.com/facebook/create-react-app/issues/8075
COPY patch.js /client/node_modules/react-dev-utils/webpackHotDevClient.js

# log npm config (for debugging)
RUN npm config list

WORKDIR /client/app

CMD ["npm", "run", "start"]

FROM dev as test
COPY . .
RUN npm audit

FROM test as build-stage
RUN npm run build

FROM nginx:1.17-alpine as prod

COPY --from=build-stage /client/app/build /usr/share/nginx/html
COPY --from=build-stage /client/app/nginx /etc/nginx/

EXPOSE 80
HEALTHCHECK CMD [ "wget", "-q", "0.0.0.0:80" ]