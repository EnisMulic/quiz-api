# quiz-api

## Dependencies

* Go - https://golang.org/
* MongoDB - https://www.mongodb.com/

## How to run

### Clone the repository and navigate to the application

```
git clone https://github.com/EnisMulic/quiz-api.git
cd quiz-api/quiz-api
```

### Create a `.env` file and set the environment variables

```
cat example.env > .env
```

### Run the api

```
go run main.go
```

### Making changes to the swagger documentation

After modifying the files according to https://goswagger.io/ run
```
make swagger
```