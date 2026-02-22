# Member Service

App to manage manager creation and authentication.

## Build

From the repository root, run:

```sh
docker build \
  -f services/member-service/Dockerfile \
  -t <registry>/member-service:latest \
  .

docker push <registry>/member-service:latest
```
