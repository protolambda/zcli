FROM golang:1.13.0-alpine3.10 as zcli_build

WORKDIR /app

COPY . .

# Build ZCLI, the ZRNT command line interface!
RUN go build -o zcli -tags preset_minimal -v -i .


FROM alpine:latest

# Set SERVICE_KEY_PATH to find the service account credentials on the machine, without checking them in to git
ARG SERVICE_KEY_PATH=service_account.key.json

WORKDIR /app

# Copy over the CLI from the build phase
COPY --from=zcli_build /app/zcli /bin/zcli
# Copy over the wrapper executable that will call the cli to run Muskoka tasks
COPY --from=protolambda/muskoka_worker:v0.1.1 /app/muskoka_worker /bin/muskoka_worker


# specify GCP project
ENV GCP_PROJECT muskoka

# Load service account
COPY ${SERVICE_KEY_PATH} service_account.key.json
# Direct application to find the service account
ENV GOOGLE_APPLICATION_CREDENTIALS service_account.key.json

ENTRYPOINT muskoka_worker --spec-version=v0.8.3 --spec-config=minimal --cli-cmd="zcli transition blocks" --worker-id=worker1 --client-name=zrnt  --client-version=v0.8.3

