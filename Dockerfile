FROM node:21-alpine3.17 as build1

WORKDIR /app

COPY public/image.svg public/image.svg
COPY package-lock.json  package.json  tailwind.config.js ./
RUN npm i

COPY web web
RUN npm run tailwind:prod

FROM golang:1.21.4-alpine3.18 as build2

WORKDIR /app

COPY --from=build1 /app/public public
COPY --from=build1 /app/web/templates web/templates

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY main.go  .
RUN go build

FROM golang:1.21.4-alpine3.18

WORKDIR /app

COPY --from=build2 /app/public public
COPY --from=build2 /app/web web
COPY --from=build2 /app/image-uploader .

CMD [ "/app/image-uploader" ]
