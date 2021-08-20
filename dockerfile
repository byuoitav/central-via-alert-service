FROM gcr.io/distroless/static
MAINTAINER Clinton Reeder <clinton_reeder@byu.edu>

ARG NAME

COPY ${NAME} /central-via-alert-service

ENTRYPOINT ["/central-via-alert-service"]
