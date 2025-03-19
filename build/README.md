## Docker Build

### Build your own Docker image

Run the following command to build your own Docker image:
```
docker buildx build -t argo-curl-plugin -f build/Dockerfile --load .
```
