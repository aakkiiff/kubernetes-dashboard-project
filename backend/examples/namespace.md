# Namespaces API Examples
This file will provide examples for how to hit the namespace routes:

## List all namespaces
```
curl -X GET http://localhost:8080/namespace
```
## Create a new namespace
```
curl -X POST http://localhost:8080/namespace \
-H "Content-Type: application/json" \
-d '{"name": "my-new-namespace"}'
```

## Delete a specified namespace
```
curl -X DELETE http://localhost:8080/namespace \
-H "Content-Type: application/json" \
-d '{"name": "my-new-namespace"}'

```