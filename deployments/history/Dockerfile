FROM golang as build

COPY . /src

WORKDIR /src

RUN CGO_ENABLED=0 GOOS=linux go build -o history ./cmd/history


FROM scratch

COPY --from=build /src/history .

EXPOSE 8080
EXPOSE 50051

CMD ["/history"]