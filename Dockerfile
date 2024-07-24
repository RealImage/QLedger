# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:1.11

# copy Lambda Web Adapter binary to /opt/extensions
COPY --from=public.ecr.aws/awsguru/aws-lambda-adapter:0.8.3 /lambda-adapter /opt/extensions/lambda-adapter

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/RealImage/QLedger

# Build the QLedger command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go install github.com/RealImage/QLedger

# Run the QLedger command by default when the container starts.
ENTRYPOINT /go/bin/QLedger

# Document that the service listens on port 7000.
EXPOSE 7000
