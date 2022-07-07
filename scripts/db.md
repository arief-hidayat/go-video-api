## Local Docker Setup
starting postgresql
```
docker run --name postgresql \
-e POSTGRESQL_USERNAME=vod123 \
-e POSTGRESQL_PASSWORD=password123 \
-e POSTGRESQL_DATABASE=vod \
-p 5432:5432 \
bitnami/postgresql:14.4.0
```

running the client
```
docker run -it --rm \
    bitnami/postgresql:14.4.0 psql -h host.docker.internal -U vod123 -d vod
```

then, execute sql statements.

