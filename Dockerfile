FROM alpine
WORKDIR /build
COPY /build/rpc-service /usr/bin/
CMD ["rpc-service"]