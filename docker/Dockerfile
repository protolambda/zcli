FROM golang:alpine
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates bash
RUN git clone https://github.com/protolambda/zcli.git
WORKDIR /zcli
RUN go get -u -tags preset_minimal github.com/protolambda/zcli
RUN wget https://raw.githubusercontent.com/ethereum/eth2.0-pm/master/interop/mocked_start/keygen_10000_validators.yaml

ENTRYPOINT ["/bin/sh"]
