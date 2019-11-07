FROM golang:1.13.0-alpine3.10 as zcli_build

WORKDIR /app

COPY . .

# Build ZCLI, the ZRNT command line interface!
RUN go build -o zcli -tags preset_minimal -v -i .


FROM alpine:latest

# Set SERVICE_KEY_B64 to the base64 encoded service_account.key.json file
ARG SERVICE_KEY_B64

WORKDIR /app

# Copy over the CLI from the build phase
COPY --from=zcli_build /app/zcli /bin/zcli
# Copy over the wrapper executable that will call the cli to run Muskoka tasks
COPY --from=protolambda/muskoka_worker:v0.1.1 /app/muskoka_worker /bin/muskoka_worker


# specify GCP project
ENV GCP_PROJECT muskoka

RUN echo -n $SERVICE_KEY_B64 > key_b64.txt
RUN base64 -d key_b64.txt > service_account.key.json

RUN rm key_b64.txt

# Direct application to find the service account
ENV GOOGLE_APPLICATION_CREDENTIALS /app/service_account.key.json

ENTRYPOINT muskoka_worker --spec-version=v0.8.3 --spec-config=minimal --cli-cmd="zcli transition blocks" --worker-id=worker1 --client-name=zrnt  --results-bucket=muskoka_zrnt  --client-version=v0.8.3

