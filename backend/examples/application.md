# Application API Examples
This file will provide examples for how to hit the application routes:



## Create a new application
```
curl -X POST http://localhost:8080/app \
-H "Content-Type: application/json" \
-d '{
	"name": "my-app",
	"image": "my-app-image:latest",
	"containerPort": 8080
}'
```
**browse you deployment on url: ingress_ip/application_name**
## update an application
```
curl -X PUT http://localhost:8080/app \
-H "Content-Type: application/json" \
-d '{
    "name": "my-app",
    "image": "my-app-image:updated",
    "containerPort": 8081
}'

```
