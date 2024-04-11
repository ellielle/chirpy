# Chirpy

Chirpy is a web server written in Go, mimicking a Twitter-like API. Built for [Boot.dev's](https://www.boot.dev) Web Servers project. It allows the creation and authentication of users, who can then make "Chirps", posts of 140 characters or less.

## Goal

This project was meant to get my feet wet with building back-ends in Go, and specifically to familiarize myself with how key concepts work in the language when handling web requests. It uses a basic `.json` file as a database, which will be created if it doesn't exist.

The more important concepts explored include:

- [x] Routing with `net/http.ServeMux` in Go 1.22
- [x] Architectural differences in web applications
- [x] Handling JSON requests and responses in Go
- [x] How to use various storage options with Go
- [x] Authentication with JWTs, refreshing and revoking
- [x] Authorization
- [x] Webhooks (with a fake payment processor 'Polka')
- [x] How painful manual documentation of an API is

The project was a bit harder than I anticipated, but I am significantly more confident with Go after comepleting it.

## ⚙️ Installation

Clone the repo with:

```bash
$ git clone git@github.com:ellielle/chirpy.git
```

Chirpy makes use of [godotenv](https://github.com/joho/godotenv) to provide environment variables to the server. You will need to create a `.env` file with 2 variables: `JWT_SECRET` and `POLKA_API_KEY`. The values can be whatever you want!

To build and run the server, use the following command. The `--debug` flag deletes the `database.json` file on load.

```bash
go build -o out && ./out --debug
```

## Usage

### POST /api/users - Create User

Request Body:

```json
{
  "email": "test@test.com",
  "password": "verysecure"
}
```

Response Body:

```json
{
  "id": 1,
  "email": "test@test.com",
  "is_chirpy_red": false
}
```

### PUT /api/users - Update User

Request Body:

```json
{
  "email": "test@test.com",
  "password": "verysecure213"
}
```

Response Body:

```json
{
  "id": 1,
  "email": "test@test.com",
  "is_chirpy_red": false
}
```

### POST /api/login - Login User

Request Body:

```json
{
  "email": "test@test.com",
  "password": "verysecure"
}
```

Response Body:

```json
{
  "id": 1,
  "email": "test@test.com",
  "is_chirpy_red": false,
  "token": "access token",
  "refresh_token": "refresh token"
}
```

### POST /api/refresh - Refresh JWT

Request Header: `"Authentication": "Bearer <refresh_token>"`

Response Body:

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXNzIiwic3ViIjoiMSIsImV4cCI6MTcxMDUyMzAwNywiaWF0IjoxNzEwNTE5NDA3fQ.NWiUcV7irFnTuHQe8qp6UIBw0qv1hgfabpbYixmBrHY"
}
```

### POST /api/revoke - Revoke JWT

Request Header: `"Authentication": "Bearer <refresh_token>"`

Response Body:

```
"OK"
```

### POST /api/chirps - Create Chirp

Request Header: `"Authentication": "Bearer <access_token>"`

Request Body:

```json
{
  "body": "chirp chirp"
}
```

### GET /api/chirps - Get all Chirps

This endpoint takes two optional query parameters:

- `?sort=` - 'asc' or 'desc'. Defaults to 'asc'
- `?author_id=` - ID of author to get chirps from

Response Body:

```json
[
  {
    "id": 1,
    "body": "chirp chirp",
    "author_id": 1
  }
]
```

### GET /api/chirps/{chirpID} - Get a single Chirp by ID

Response Body:

```json
{
  "id": 4,
  "body": "chirp chirp birp",
  "author_id": 2
}
```

### DELETE /api/chirps/{chirpID} - Delete a Chirp

Request Header: `"Authentication": "Bearer <access_token>"`

Response Body:

```
"OK"
```

### POST /api/polka/webhooks - Endpoint to receive events from "Polka"

Request Header: "Authentication": "ApiKey <polka_api_key>"
Request Body:

```json
{
  "event": "user.upgraded",
  "data": {
    "user_id": 1
  }
}
```

Response Body:

```
"OK"
```

### GET /admin/metrics - Simple middleware that counts hits to the fileserver

Resopnse Header:

```
200 OK
```

Response Body:

```html
<!doctype html>
<html>
  <head>
    <meta
      name="generator"
      content="HTML Tidy for HTML5 for Linux version 5.6.0"
    />
    <title></title>
  </head>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited 0 times!</p>
  </body>
</html>
```

### GET /admin/reset - Reset the counter

Response Header:

```
200 OK
```

Response Body:

```
OK
```

### GET /api/healthz - Health check endpoint

Response Header:

```
200 OK
```

Response Body:

```
OK
```
