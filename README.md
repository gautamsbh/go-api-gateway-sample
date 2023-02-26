## Microservice Auth App

### Build Instruction

* step 1: Run all services defined in docker-compose.yml file command: `docker-compose up -d`
  to rebuild the project: `docker-compose up -d --build`
* step 2: Run DB migration and Test command: `docker-compose run auth_microservice go test -v`

### Test APP

* /profile endpoint:
  Request:
  ```shell

curl --location --request GET 'localhost:8000/profile' \
--header 'username: AKS' \
--data-raw ''```

1. Success Case: If username is valid (exists in DB)
   Response:
   ```json

{
"first_name": "alok",
"id": 1,
"last_name": "sonker",
"username": "AKS"
}```

2. Failure: If username not sent in headers:
   Response: status code: 400
   ```json

{
"error": "username not found in request header"
}```

3. Failure: If username is not valid (not exists in DB)
   Response: status code: 401
   ```json

{
"error": "unable to process request"
}```

* /service endpoint:
  Request:
  ```shell

curl --location --request GET 'localhost:8000/service' \
--data-raw ''```

1. Success Case: If microservice is alive:
   Response: status code: 200
   ```json

{
"message": "user-microservice"
}```

2. Failure Case: If microservice is down:
   Response: status code: 500
   ````json

{
"error": "internal server error"
}````
