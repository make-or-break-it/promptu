FROM golang:1.18 as base

# --- Build
FROM base as build

WORKDIR /code

COPY . .

RUN go mod download && \
    go mod verify && \
    mkdir /build && \
    go build -v -o /build/server .

# --- Run
FROM base as run

COPY --from=build /build/server /build/

ENTRYPOINT [ "/build/server" ]