FROM golang:1.13.0-alpine3.10 as zcli_build

WORKDIR /app

COPY . .

RUN go build -o zcli -tags preset_minimal -v -i .


FROM alpine:latest

WORKDIR /app

COPY --from=zcli_build /app/zcli /bin/zcli
COPY --from=muskoka_worker /app/muskoka_worker /bin/muskoka_worker

ENTRYPOINT muskoka_worker --spec-version=v0.8.3 --spec-config=minimal --cli-cmd="zcli transition blocks" --worker-id=worker1 --client-name=zrnt  --client-version=v0.8.3

