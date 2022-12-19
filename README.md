Run the docker image with
```
docker run -p 8080:8080 --name delete --rm \
      -e SMTP_HOST=XXXX \
      -e SMTP_PORT=XXXX \
      -e SMTP_FROM=XXXX \
      -e SMTP_PASSWORD=XXXX \ 
      -e MONGO_USERNAME=XXXX \
      -e MONGO_PASSWORD=XXXX \
      -e MONGO_DATABASE=XXXX \
      -e MONGO_COLLECTION=XXXX \
      test
      /
```