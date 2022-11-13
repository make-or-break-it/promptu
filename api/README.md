# API

The API is written in Golang.

## Run

```go run cmd/main.go```

### Running using docker

Build container image:
```
docker build -t promptu-api .
```

Run container:
```
docker run -p 8080:8080 promptu-api
```

## Supported endpoints

Post with:
```
curl -i -XPOST  -H "Content-Type: application/json" -d '{"user":"some_user","message":"hey there. i feel great!"}' http://localhost:8080/post
```

Get feed with:
```
curl -i http://localhost:8080/feed
```
