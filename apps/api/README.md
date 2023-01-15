# API

The API is written in Golang.

## Run

```
go run cmd/main.go
```

### Running using docker

Build container image:
```
docker build -t promptu-api .
```

Run container:
```
docker run -p 8080:8080 promptu-api
```

To interact with the api from your host, you need to set these environment variables (you can create a `.env` file in this folder with them).
```
HOST=0.0.0.0
PORT=8080
```

## Supported endpoints

Post with:
```
curl -i -XPOST  -H "Content-Type: application/json" -d '{"user":"some_user","message":"Test message from localhost", "utcNotificationTime":"2022-12-09T02:17:54.760Z"}' http://localhost:8080/post
```

Get feed with:
```
curl -i http://localhost:8080/feed
```
