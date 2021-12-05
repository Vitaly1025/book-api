# BookAPI
This is API for working with books in shop, where you can set price for books, get books and delete books.
### Features
This project implements a clean multi layered architecture that helps to separate the layers among themselves. Each layer has its own area of responsibility and is independent of each other.
   - ```pkg/handler```  - describes logic for receiving request and sending response
   - ```pkg/service```  - describes some buisness logic
   - ```pkg/repository```  - describes logic for working with database 
   - ```pkg/model```  - describes models
   - ```pkg/middleware```  - describes helper logic for processing errors

REST API implements OpenAPI specification. It means that Swagger UI can be used for making request to the API. In section Running you can find instruction how to browse endpoints.
### Preparation

1. Setup golang. See: https://go.dev/doc/install
2. Setup Docker. See: https://docs.docker.com/get-docker/
3. Clone the source down to your machine 

    ```git clone https://github.com/Vitaly1025/book-api.git```
4. Download all GO modules. Run the next command from the local source's directory.

   ```go mod download```
5. Install a database migration tool

   ```go get -tags 'postgres' -u github.com/golang-migrate/migrate/v4/cmd/migrate/```
## Running
### Manual mode
1. Firstly, you need to run database. The simplest way is Docker. In my case i using POSTGRES.  Use the next command:

    ```docker run --name=book-db -e POSTGRES_PASSWORD='1234qwerty' -p 5436:5432 -d --rm postgres```
2. Next, you need migrate the database and create database's resources (tables and the initial data). 

   ```migrate -path ./schema -database 'postgres://postgres:1234qwerty@db:5432/postgres?sslmode=disable' up```
3. Run the bookstore server:

   ```go run ./cmd/main.go```
   If you see the message ```API is running on: 4000``` it means that the server is ready to serve your requests. You can configure port setting via ``/config/config.yml`` or environment variables, which will be have more high priority
4. To see list of available endpoints, run the next address via a browser: 
   ```http://localhost:4000/swagger/index.html```
  
### ```docker-compose``` mode
1. Run this command for building and starting your compose

    ``docker-compose up -d --build``
2. To see list of available endpoints, run the next address via a browser: 

   ```http://localhost:4000/swagger/index.html```
## Testing
App have testing on two layers. It's recieving request and working with database. It's because service layer haven't logic for testing on my opinion
1. Run the next command from the local source's directory:

   ```go test .\pkg\repository && go test .\pkg\service```
