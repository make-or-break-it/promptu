FROM node:19 as base

# --- Build
FROM base as build

COPY . .

ARG PUBLIC_PROMPTU_API_HOST="http://localhost:8080"

RUN npm ci --only=production && \
    npm run build

# --- Run
FROM base as run

WORKDIR /code

COPY --from=build /dist /dist

# For more production grade proxying, see nginx:
# https://sveltesociety.dev/recipes/publishing-and-deploying/dockerize-a-svelte-app
RUN npm install -g serve

ENTRYPOINT [ "serve", "-s", "/dist" ]