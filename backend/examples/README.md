# API Endpoints

## Namespaces - [click here for request examples](./namespace.md)
- `GET /namespace`: List all namespaces
- `POST /namespace`: Create a new namespace
- `DELETE /namespace`: Delete a specified namespace

## Pods - [click here for request examples](./pod.md)
- `GET /pods/:ns`: List all pods in a namespace
- `POST /pod`: Create a new pod
- `DELETE /pod`: Delete a specified pod
- `PUT /pod`: Update a pod's image

## Applications - [click here for request examples](./application.md)
- `POST /app`: Create a new application (automatically creates a namespace, deployment, service, and ingress)
**browse you deployment on url: ingress_ip/application_name**
- `PUT /app`: Update an existing application

## Nodes
- `GET /node`: List all nodes in the cluster
