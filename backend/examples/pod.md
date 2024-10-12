
# Pods API Examples
This file will provide examples for how to hit the pod routes:


## List all pods in a namespace

 ```bash curl -X GET http://localhost:8080/pods/my-namespace```
## Create a new pod
```
curl -X POST http://localhost:8080/pod \
-H "Content-Type: application/json" \
-d '{
    "name": "my-pod",
    "namespaceName": "my-namespace",
    "image": "my-image:latest"
}'
```

## Delete a specified pod
```
curl -X DELETE http://localhost:8080/pod \
-H "Content-Type: application/json" \
-d '{"name": "my-pod", "namespaceName": "my-namespace"}'
```

## update a pods image
```
curl -X PUT http://localhost:8080/pod \
-H "Content-Type: application/json" \
-d '{
    "name": "my-pod",
    "namespaceName": "my-namespace",
    "image": "my-image:updated"
}'
```
