# Humpback Whale Social Call
[Official HWSC Web Application](https://hwsc-org.github.io/hwsc-app-gateway-svc/)

## Overview
The project is a web application that biologists can use to consolidate and to track data about humpback whales and other marine mammals. These data includes media files and source codes that biologists use for their research.
This is a non-profit project dedicated for educational and research purposes.

## Design & Specification
The user can access the application uses a device with a browser to search for data related to humpback whales. The user can also get data using the web api provided, and data are returned as JSON.
The web application is client-rendering, therefore an added initial time is added to load the JavaScript files and to download other libraries needed by the browser to run and to render the application. This reduces network traffic between the client and the server, therefore enabling the client to work remotely and only communicate back to the server if data is needed. 

## Microservices
Each microservice corresponds to specific functionality. This enables the usage of containers for deployment later on.
### Application Gateway
Services the user as the gateway, and the middleware between services.
- [hwsc-app-gateway-svc](https://github.com/hwsc-org/hwsc-app-gateway-svc)

### Frontend Service
Handles graphical user interface for Vue components.
- [hwsc-frontend](https://github.com/hwsc-org/hwsc-frontend)

### User Service
Provides services to hwsc-app-gateway-svc for managing HWSC users.
- [hwsc-user-svc](https://github.com/hwsc-org/hwsc-user-svc)

### File Transaction Service
Provides services to hwsc-app-gateway-svc for managing files.
- [hwsc-file-transaction-svc](https://github.com/hwsc-org/hwsc-file-transaction-svc)

### Document Service
 Provides services to hwsc-app-gateway-svc for CRUD of documents and file metadata in Azure CosmosDB.
- [hwsc-document-svc](https://github.com/hwsc-org/hwsc-document-svc)

### Tutorial gRPC Service in GoLang
- [hwsc-grpc-sample-svc](https://github.com/hwsc-org/hwsc-grpc-sample-svc)

### API
API for hwsc-app-gateway-svc to consume various hwsc services.
- [hwsc-api-blocks](https://github.com/hwsc-org/hwsc-api-blocks)

### Logger
Utility for logging in services.
- [hwsc-logger](https://github.com/hwsc-org/hwsc-logger)

## Team 
### Owners
- Kerri
- Shima

### Members
- Pouria 
- Conard 
- Pavel 
- Ashley 
- Mandy 
- Alex
- Lisa
- Wai

## Links
- [Q1 Website](https://hwss.azurewebsites.net/#/)
- [Team Google Drive](https://drive.google.com/drive/folders/13vJqlP3PRIZJMuMC0tfnGKSoOrWuMX4W)
- [Udemy](https://www.udemy.com/)
- [Vue JS](https://vuejs.org/)
- [npm js](https://www.npmjs.com/)
- [Bootstrap](https://getbootstrap.com/)
- [MongoDB](https://www.mongodb.com/)
- [Azure Portal](https://portal.azure.com/)
- [Google Map API JS](https://developers.google.com/maps/documentation/javascript/tutorial)
- [gRPC](https://grpc.io/)
- [protocol buffers](https://developers.google.com/protocol-buffers/docs/proto3)
- [gRPC with REST](https://grpc.io/blog/coreos)
