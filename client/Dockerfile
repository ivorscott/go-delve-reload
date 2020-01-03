FROM node:10.15.0-alpine as base
ENV NODE_ENV=production

# create a parent directory for the app
WORKDIR /client

# copy the package.json and lock files
COPY package*.json ./

# install production dependencies & clean up after
RUN npm ci \ 
    && npm cache clean --force

FROM base as dev
ENV NODE_ENV=development

# set node_module executables (*/.bin) inside path 
ENV PATH /client/node_modules/.bin:$PATH

# expose default create-react-app port
EXPOSE 3000

# make app directory and update permissions to use node user
RUN mkdir /client/app && chown -R node:node .

# change user to node user (otherwise we are always root in the container)
USER node
RUN npm i --only=development \
    && npm cache clean --force

# log npm config (for debugging)
RUN npm config list
WORKDIR /client/app
CMD ["npm", "run", "start"]

FROM dev as test
COPY . .
RUN npm audit

FROM test as build-stage
RUN npm run build

FROM nginx:1.15-alpine as prod
EXPOSE 80
COPY --from=build-stage /client/app/build /usr/share/nginx/html
COPY --from=build-stage /client/app/nginx.conf /etc/nginx/conf.d/default.conf