# Kubernetes API Project

## Overview

This project is a RESTful API built with Go and Gin, designed to interact with Kubernetes clusters. It allows users to perform CRUD operations on Kubernetes resources, including Pods, Namespaces, and Applications. When creating an application, the API automatically provisions a Namespace, Deployment, Service, and Ingress, exposing the application at a specified path.

## Features


- **Automatic Resource Creation**:
	- **namespace**
	- **pods**
	- **service**
	- **ingress**
	- **application(custom resource):** When creating an application, the API automatically provisions a full application stack, including a Namespace, Deployment, Service (ClusterIP), and Ingress, which exposes the application at the path `ip/app_name`.
	- comming soon!
	-  comming soon!


## Requirements

- Go 1.23.0 was used
- Kubernetes cluster access - (kubeadm,eks,kind,minikube compatibility)
- NGINX Ingress Controller (for ingress functionality) 
	- https://github.com/kubernetes/ingress-nginx/tree/main/charts/ingress-nginx

## Getting Started

### Installation

1. **install**
  ```
   git clone https://github.com/aakkiiff/kubernetes-crud-api.git
   cd kubernetes-crud-api
   go mod tidy
```

2. **Run the application**
`go run main.go`
3. **Access the API**: The API will be running on `http://localhost:8080`. Use tools like Postman or curl to interact with the endpoints.

## Examples
- [click here for examples](./examples/README.md)