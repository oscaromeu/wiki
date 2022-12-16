
FROM node:lts as base
ENV NPM_CONFIG_LOGLEVEL=warn
ENV NPM_CONFIG_COLOR=false
WORKDIR /app
COPY . ./
COPY docusaurus.config-dev.js docusaurus.config.js
RUN yarn install
RUN yarn build

# Deployment step

FROM busybox:latest as deploy

RUN adduser -D static
USER static
WORKDIR /home/static

COPY --from=base app/build/ ./

EXPOSE 3000

CMD ["busybox", "httpd", "-f", "-v", "-p", "3000"]
