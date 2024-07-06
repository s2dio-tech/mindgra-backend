# build the app.
FROM golang:1.22-rc-alpine AS build-env

# set the working dir.
WORKDIR /app

# copy the go module dependency files.
COPY go.mod go.sum ./

# download the go module dependencies.
RUN go mod download

# copy the app source code.
COPY application ./application
COPY common ./common
COPY datasource ./datasource
COPY docs ./docs
COPY domain ./domain
COPY internal ./internal

# build the go app binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app /app/application/main.go

# set up the container.
FROM golang:1.22-rc-alpine

# set the working dir.
WORKDIR /app

# copy the built app binary from the build-env.
COPY --from=build-env /app/app ./app

# expose the port.
EXPOSE 8080

# command to run the app.
CMD ["./app"]
