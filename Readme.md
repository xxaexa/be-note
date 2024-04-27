# Backend Notes
fe = https://github.com/xxaexa/notelify-fe

## API Spec

### User API

#### Register User

Request :

- Method : POST
- Endpoint : `api/v1/auth/register`
- Header :
  - Content-Type : application/json
  - Accept : application/json
- Body :

```json
{
	"email": "string",
  "password": "string",
}
```

Response :

- Status Code: 201 Created
- Body:

```json
{
  "data": bool,
  "message": "string"
}
```

#### Login User

Request :

- Method : GET
- Endpoint : `api/v1/auth/login`
  - Header :
  - Accept : application/json

Response :

- Status Code : 200 OK
- Body:

```json
{
  "data": {
    "token": "string",
    "user": {
      "email": "string",
      "username": "string"
    }
  },
  "message": "string"
}
```



### Note API

#### Create Note

Request :

- Method : POST
- Endpoint : `/api/v1/notes/`
- Header :
  - Content-Type : application/json
  - Accept : application/json
- Body :

```json
{
  "title": "string",
  "description": "string",
  "user_id": int
}
```

Response :

- Status Code: 201 Created
- Body:

```json
{
  "data": bool,
  "message": "string"
}
```

#### Get All Note

Request :

- Method : GET
- Endpoint : `/api/v1/notes/`
  - Header :
  - Accept : application/json

Response :

- Status Code : 200 OK
- Body:

```json
{
	"data": [
		{
      "id": int,
      "title": "string",
      "description": "string",
      "user_id": int
    }
		{
      "id": int,
      "title": "string",
      "description": "string",
      "user_id": int
    }
	],
  "message": "string",
}
```

#### Get Note By Id

Request :

- Method : GET
- Endpoint : `/api/v1/notes/:id`
- Header :
  - Accept : application/json

Response :

- Status Code: 200 OK
- Body :

```json
{
  "data": {
    "id": int,
    "title": "string",
    "description": "string",
    "user_id": int
  },
  "message": "string"
}
```

#### Update Note

Request :

- Method : PUT
- Endpoint : `/api/v1/notes/:id`
- Header :
  - Content-Type : application/json
  - Accept : application/json
- Body :

```json
{
  "title": "string",
  "description": "string",
  "user_id": int
}
```

Response :

- Status Code: 200 OK
- Body :

```json
{
  "data": int,
  "message": "string"
}
```

#### Delete Note

Request :

- Method : DELETE
- Endpoint : `/api/v1/notes/:id`
- Header :
  - Accept : application/json
- Body :

Response :

- Status : 200 OK
- Body :

```json
{
  "data": int,
  "message": "string"
}
```