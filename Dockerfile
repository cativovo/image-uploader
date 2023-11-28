FROM node:21-alpine3.17 as build1

WORKDIR /app

COPY package-lock.json  package.json  tailwind.config.js ./
COPY web web

RUN npm i

RUN npm run tailwind:prod

COPY public public
COPY go.mod go.sum main.go ./

FROM golang:1.21.4-alpine3.18 as build2

WORKDIR /app

COPY --from=build1 /app/go.mod /app/go.sum ./
COPY --from=build1 /app/main.go .
COPY --from=build1 /app/public public
COPY --from=build1 /app/web/templates web/templates

RUN go mod download && go mod verify

RUN go build

FROM golang:1.21.4-alpine3.18

WORKDIR /app

COPY --from=build2 /app/web web
COPY --from=build2 /app/public public
COPY --from=build2 /app/image-uploader .

CMD [ "/app/image-uploader" ]
