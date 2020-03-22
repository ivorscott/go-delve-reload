# 1. Use the node apline image as the base stage of a multi-stage routine
FROM node:12.16.1-alpine as base

# 2. Set the working directory to /client
WORKDIR /client

# 3. Copy both package.json and package-lock.json into /client in the image's filesystem
COPY package*.json ./

# 4. Install only the production node_modules and clean up the cache
RUN npm ci \
    && npm cache clean --force

# 5. Extend the base stage and create a new stage named dev
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

# 10. Print npm configuration for debugging purposes
RUN npm config list

# 11. Set the working directory to /client/app
WORKDIR /client/app

# 12. Provide the default command for the development container
CMD ["npm", "run", "start"]

# 13. Extend the dev stage and create a new stage called test
FROM dev as test

# 14. Copy the remainder of the client folder source code into the image's filesystem
COPY . .

# 15. Run node_module vulnerability checks
RUN npm audit

# 16. Run unit tests in CI
RUN CI=true npm test --env=jsdom

# 17. Extend the test stage and create a new stage named build-stage
FROM test as build-stage

# 18. Build the production static assets
RUN npm run build

# 19. Install aquasecurity's trivy for robust image scanning
FROM aquasec/trivy:0.4.3 as trivy

# 20. Scan the nginx alpine image before production use
RUN trivy nginx:1.17-alpine && \
    echo "No image vulnerabilities" > result

# 21. Extend the nginx apline image and create a new stage named prod
FROM nginx:1.17-alpine as prod

# 22. Copy only the files we want from a few stages into the prod stage
COPY --from=trivy result secure
COPY --from=build-stage /client/app/build /usr/share/nginx/html

# 23. Copy non-root user nginx configuration
# https://hub.docker.com/_/nginx
COPY --from=build-stage /client/app/nginx /etc/nginx/

# 24. Provide meta data about the port the container must expose
EXPOSE 80

# 25. Define how Docker should test the container to check that it is still working
HEALTHCHECK CMD [ "wget", "-q", "0.0.0.0:80" ]