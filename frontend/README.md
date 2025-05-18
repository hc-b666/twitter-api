### Twitter Frontend

Build docker image
```bash
docker build -t twitter-client .
```

Run docker container
```bash
docker run -d -p 5173:80 twitter-client
```
