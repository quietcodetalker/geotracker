FROM golang as build

COPY . /src

WORKDIR /src

RUN CGO_ENABLED=0 GOOS=linux go build -o locations ./cmd/locations


FROM scratch

COPY --from=build /src/locations .

EXPOSE 8080
EXPOSE 50051

CMD ["/locations"]