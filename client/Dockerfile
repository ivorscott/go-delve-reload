# 1. Use the Node image as the base stage of a multi-stage routine
FROM node:10.15.0-alpine as base

# 2. Set the working directory to /client
WORKDIR /client

# 3. Copy both package.json and package-lock.json into /client in the image's filesystem
COPY package*.json ./

# 4. Install only the production node_modules and clean up the cache
RUN npm ci \
    && npm cache clean --force

# 5. Extend the base stage to create a new stage named dev
FROM base as dev

# 6. Set the NODE_ENV and PATH Environment variables
ENV NODE_ENV=development
ENV PATH /client/node_modules/.bin:$PATH

# 7. Provide meta data about the port the container must expose
EXPOSE 3000

# 8. Create a new /app directory in /client
RUN mkdir /client/app

# 9. Install development dependencies
RUN npm i --only=development \
    && npm cache clean --force

# 10. Patch create-react-app bug preventing self-signed certificate usage
# https://github.com/facebook/create-react-app/issues/8075
COPY patch.js /client/node_modules/react-dev-utils/webpackHotDevClient.js

# 11. Print npm config for debugging purposes
RUN npm config list

# 12. Change directory
WORKDIR /client/app

# 13. Provide the default command for the development container
CMD ["npm", "run", "start"]

# 14. Extend the dev stage and create a new stage called test
FROM dev as test

# 15. Copy the remainder of the client folder source code into the image's filesystem
COPY . .

# 16. Run node_module vulnerability checks
RUN npm audit

# 17. Run unit tests
RUN npm test

# 18. Extend the test stage to create a new stage named build-stage
FROM test as build-stage

# 19. Build the production static assets
RUN npm run build

# 20. Extend the Nginx image to create a new stage named prod
FROM nginx:1.15-alpine as prod

# 21. Copy only the files we want from the build stage into the prod stage
COPY --from=build-stage /client/app/build /usr/share/nginx/html

# 22. Copy non-root user nginx configuration
# https://hub.docker.com/_/nginx
COPY --from=build-stage /client/app/nginx /etc/nginx/

# 23. Provide meta data about the port the container must expose
EXPOSE 80

# 24. Define how Docker should test the container to check that it is still working
HEALTHCHECK CMD [ "wget", "-q", "0.0.0.0:80" ]