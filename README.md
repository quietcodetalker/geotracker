# geotracker

`geotracker` is a pet project that represents a geo tracker application.
It's built as a distributes microservice application for academic reasons.

## Usage

Load dependencies
```bash
go mod tidy
```

Build images
```bash
make build_history_image
make build_locations_image
```

Run the following command to run a local cluster

```bash
docker-compose -f ./deployments/docker-compose.yml up
```

It will listen for requests on `localhost:10000`.

Run database migrations

```bash
make migrate_locations_up
make migrate_history_up
```

## Structure

It consists of two microservices:
- Locations
- History

Local cluster based on docker uses [Envoy][envoy] as a gateway.

## Documentation

To see documentation open `http://localhost:8080/` in your browser to open [swagger][swagger]
 
[envoy]: https://www.envoyproxy.io/
[swagger]: https://swagger.io/