# syntax=docker/dockerfile:1

FROM golang:1.18-buster

WORKDIR /go/src/vendingmachine


COPY /go.mod ./
COPY /go.sum ./
COPY /*.go ./
COPY /config/*.go ./config/
COPY /controllers/*.go ./controllers/
COPY /domains/*.go ./domains/
COPY /services/*.go ./services/
COPY /repositories/*.go ./repositories/

ENV SELLER_CLIENT_ID=vendingmachine-app
ENV SELLER_CLIENT_SECRET=t3mp0r4ry-v41u3
ENV REALM=vendingmachine
ENV CLOAK_HOST=http://localhost:8080
ENV REALM_ADMIN_USER=vending-machine-admin
ENV REALM_ADMIN_PASSWORD=dummy-password
ENV KEYCLOAK_ADMIN_USER=admin
ENV KEYCLOAK_ADMIN_PASSWORD=admin

RUN go mod download \
&& mkdir -p /tmp/test-coverage/ \
&& go test -timeout 30s -coverprofile=/tmp/test-coverage/go-code-cover vendingmachine/services \
&& go build -o vendingmachine \
&& echo "vendingmachine built successfully"

EXPOSE 8099

CMD [ "./vendingmachine" ]