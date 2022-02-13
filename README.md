# geotracker

`geotracker` is a pet project that represents a geo tracker application.
It's built as a distributes microservice application for academic reasons.

## Usage

### Docker

Run the following command to run a local cluster

`docker-compose -f ./deployments/docker-compose.yml up`

It will listen for requests on `localhost:10000`.

### Without Docker

Run following commands in different terminals:

- `make run_locations`  
  Locations microservice will be available on `localhost:8080`
- `make run_history`  
  History microservice will be available on `localhost:8081`

## Structure

### Components

It consists of two microservices:
- Locations
- History

Local cluster based on docker uses [Envoy][envoy] as a gateway.

### Endpoints

It serves 3 endpoints:
- Set location
- List users in radius
- Get distance

#### Set Location

It sets a user's location by username.

Request:

`PUT` `/v1/users/{username}/location`

`Content-Type: application/json`

```json
{
  "latitude": 0.0,
  "longitude": 0.0
}
```

Response:

`Content-Type: application/json`

```json
{
  "latitude": 0.0,
  "longitude": 0.0
}
```

#### List users in radius

It returns a paginated list of users by radius in meters and geo coordinates.

Request:

`GET` `/v1/users/radius`

Query params`radius`, `latitude`, `longitude` are required.  
Either `page_size` or `page_token` is required as well.

Response:

`Content-Type: application/json`

```json
{
  "users": [
    {
      "id": 1,
      "username": "luke",
      "created_at": "2022-02-13T16:46:43.822482Z",
      "updated_at": "2022-02-13T16:46:43.822482Z"
    }
  ],
  "next_page_token": "opaque_token"
}
```

#### Get distance

Returns distance walked by a user with given `username` in a period of time.

Request:

`GET` `/v1/users/{username}/distance`

Response:

`Content-Type: application/json`

```json
{
  "distance": 100.0
}
```

#### Possible errors

These are possile errors endpoinds can respond with:

##### Internal Error

HTTP Status: `500`

Response body:

```json
{
  "error": {
    "code": 500,
    "message": "internal error",
    "status": "INTERNAL"
  }
}
```

##### Invalid Argument Error

HTTP Status: `400`

Response body:

```json
{
  "error": {
    "code": 400,
    "message": "invalid argument",
    "status": "INVALID_ARGUMENT"
  }
}
```

##### Not Found Error

HTTP Status: `404`

Response body:

```json
{
  "error": {
    "code": 404,
    "message": "not found",
    "status": "NOT_FOUND"
  }
}
```

##### Failed Precondition Error
 
HTTP Status: `422`

Response body:

```json
{
  "error": {
    "code": 422,
    "message": "failed precondition",
    "status": "FAILED_PRECONDITION"
  }
}
```
 
[envoy]: https://www.envoyproxy.io/