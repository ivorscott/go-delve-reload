# 1. Capture ARGs when provided
ARG backend
ARG version

# 2. Use the node apline image as the base stage of a multi-stage routine
FROM node:12.16.1-alpine as dev

# 3. Set the working directory to /client
WORKDIR /client

# 4. Copy both package.json and package-lock.json into /client in the image's filesystem
COPY package*.json ./

# 5. Install development dependencies
RUN npm install --only=development && npm cache clean --force

# 6. Print npm configuration for debugging purposes
RUN npm config list

# 7. Set the NODE_ENV and PATH Environment variables
ENV NODE_ENV=development
ENV PATH /client/node_modules/.bin:$PATH

# 8. Provide meta data about the port the container must expose
EXPOSE 3000

# 9. Set the working directory to /client/app
WORKDIR /client/app

# 10. Provide the default command for the development container
CMD ["npm", "run", "start"]

# 11. Extend the dev stage and create a new stage called test
FROM dev as test

# 12. Copy both package.json and package-lock.json to /client/app for testing and build
COPY package*.json ./

# 13. Run node_module vulnerability checks
RUN npm audit

# 14. Install production dependencies
RUN npm ci --silent

# 15. Copy the client folder source code into the image's filesystem
COPY . .

# 16. Run unit tests
RUN CI=true npm run test

# 17. Extend the test stage and create a new stage named build-stage
FROM test as build-stage

# 18. Use ARGs in build-stage
ARG backend
ARG version

# 19. Create .env file for static build with production vars
RUN echo $'\n\
    REACT_APP_BACKEND='$backend$'\n\
    REACT_APP_VERSION='$version$'' >> .env

# 20. Build the production static assets
RUN npm run build

# 21. Install aquasecurity's trivy for robust image scanning
FROM aquasec/trivy:0.4.4 as trivy

# 22. Scan the nginx alpine image before production use
RUN trivy nginx:1.17-alpine && \
    echo "No image vulnerabilities" > result

# 23. Extend the nginx apline image and create a new stage named prod
FROM nginx:1.17-alpine as prod

# 24. Copy only the files we want from a few stages into the prod stage
COPY --from=trivy result secure
COPY --from=build-stage /client/app/build /usr/share/nginx/html

# 25. Copy non-root user nginx configuration
# https://hub.docker.com/_/nginx
COPY --from=build-stage /client/app/nginx /etc/nginx/

# 26. Provide meta data about the port the container must expose
EXPOSE 80

# 27. Run nginx
CMD ["nginx", "-g", "daemon off;"]