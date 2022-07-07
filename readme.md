## What is this?
Simple app to search video (postgreSQL full text search)

### Prepare Your DB

see [DB setup docs](./scripts/db.md)

### Run App

Example to connect the default docker setup
```
DB_NO_SSL=true DB_PWD=password123 go run ./server/cmd/main.go
```
or if you want to run on docker
```
docker build -t api .
docker run -e DB_NO_SSL=true -e DB_PWD=password123 -p 3000:3000 api 
```

### Test API
```
curl "http://127.0.0.1:8000/videos?q=sea"
```
or can use tools like [vegeta](https://github.com/tsenart/vegeta)
```
  echo "GET http://127.0.0.1:8000/videos?q=sea" | vegeta attack -rate 100/1s -timeout 2s -duration=30s | tee results.bin | vegeta report
  vegeta report -type=json results.bin > metrics.json
  cat results.bin | vegeta plot > plot.html
  cat results.bin | vegeta report -type="hist[0,100ms,200ms,300ms]"
```

