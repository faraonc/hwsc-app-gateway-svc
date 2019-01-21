#!/bin/bash

# Ensure you have a Docker Hub account and belongs to hwsc organization
docker build -t dev-hwsc-app-gateway-svc .
docker tag dev-hwsc-app-gateway-svc hwsc/dev-hwsc-app-gateway-svc
docker push hwsc/dev-hwsc-app-gateway-svc