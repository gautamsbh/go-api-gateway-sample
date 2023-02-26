## Microservice Auth App

### Build Instruction

* step 1: Run all services defined in docker-compose.yml file command: `docker-compose up -d`
  to rebuild the project: `docker-compose up -d --build`
* step 2: Run DB migration and Test command: `docker-compose run auth_microservice go test -v`

### Test APP

* Profile endpoint */profile*:
  Request:
```
    curl --location --request GET 'localhost:8000/profile' --header 'username: AKS' --data-raw ''
```

Case: 1. Success Case: If username is valid (exists in DB)
   Response:
```
        {
        "id": 1,
        "first_name": "alok",
        "last_name": "sonker",
        "username": "AKS"
        }
```

Case 2. Failure: If username not sent in headers:
   Response: status code: 400
```
        {
        "error": "username not found in request header"
        }
```

Case 3. Failure: If username is not valid (not exists in DB)
   Response: status code: 401
```
        {
        "error": "unable to process request"
        }
```

* Service endpoint */service*:
  Request:
```
        curl --location --request GET 'localhost:8000/service' --data-raw ''
```

Case 1. Success Case: If microservice is alive:
   Response: status code: 200
```
        {
        "message": "user-microservice"
        }
```

Case 2. Failure Case: If microservice is down:
   Response: status code: 500
   ```
        {
        "error": "internal server error"
        }
```
