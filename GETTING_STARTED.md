# hwsc-app-gateway-svc

## Purpose
 Provides services to hwsc-frontend from the cluster

## Contract
The proto file and compiled proto buffers are located in [hwsc-api-blocks](https://github.com/hwsc-org/hwsc-api-blocks/tree/master/int/hwsc-app-gateway-svc/proto).

### GetStatus
- Gets the current status of the cluster
### CreateUser
- Creates a user
- Returns the user with password field set to empty string
### DeleteUser
- Deletes a user
- Returns the deleted user (TODO decide if we really need to return this to chrome)
### UpdateUser
- Updates a user
- Returns the updated user
### AuthenticateUser
- Looks through users and perform email and password match
- Returns matched user
### ListUsers
- Retrieves all the users
- Returns a collection of users
### GetUser
- Retrieves a user given UUID
- Returns found user
### ShareDocument
- Shares a user's document to another user
### CreateDocument
- Creates a document
- Returns the document
### ListUserDocumentCollection
- Retrieves all documents for a specific user with the given UUID
- Returns a collection of documents
### UpdateDocument
- (completely) Updates a document using DUID
- Returns the updated document
### DeleteDocument
- Deletes a document using DUID
- Returns the deleted document
### AddFile
- Adds a new file
- Returns the updated document
### DeleteFile
- Deletes a a file
- Returns the updated document
### ListDistinctFieldValues
- Retrieves all the unique fields values required for the front-end drop-down filter
- Returns the query transaction
### QueryDocument
- Queries the document service with the given query parameters
- Returns a collection of documents

## Prerequisites
- GoLang version [go 1.11.4](https://golang.org/dl/)
- GoLang Dependency Management [dep](https://github.com/golang/dep)
- Go Source Code Linter [golint](https://github.com/golang/lint)
- Docker
- [Optional] If a new proto file and compiled proto buffer exists in [hwsc-api-blocks](https://github.com/hwsc-org/hwsc-api-blocks/tree/master/int/hwsc-app-gateway-svc/proto), update dependency ``$dep ensure -update``

## How to Run without Docker Container
1. Install dependencies and generate vendor folder ``$ dep ensure -v``
2. Update ENV variables
3. Run main ``$ go run main.go``

## How to Run with Docker Container
1. Install dependencies and generate vendor folder ``$ dep ensure -v``
2. ``$ generate_container.sh``
3. Find your image ``$ docker images``
4. Acquire ``env.list`` configuration
5. ``$ docker run --env-file ./env.list -it -p 50055:50055 <imagename>``

## How to Unit Test
1. ``$ cd service``
2. For command-line summary, ``$ go test -cover -v``
3. For comprehensive summary, ``$ bash unit_test.sh``

