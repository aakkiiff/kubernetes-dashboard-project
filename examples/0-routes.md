# API Endpoints

## Namespaces
- `GET /namespace`: List all namespaces
- `POST /namespace`: Create a new namespace
- `DELETE /namespace`: Delete a specified namespace

## Pods
- `GET /pods/:ns`: List all pods in a namespace
- `POST /pod`: Create a new pod
- `DELETE /pod`: Delete a specified pod
- `PUT /pod`: Update a pod's image

## Applications
- `POST /app`: Create a new application (automatically creates a namespace, deployment, service, and ingress)
- `PUT /app`: Update an existing application

## Nodes
- `GET /node`: List all nodes in the cluster
