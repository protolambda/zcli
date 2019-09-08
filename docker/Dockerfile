FROM ubuntu:18.04 as built
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update &&\
    apt-get install -y apt-utils expect git git-extras software-properties-common \
    wget curl build-essential golang-go make

RUN git clone https://github.com/protolambda/zcli.git
WORKDIR /zcli

RUN go get -u -tags preset_minimal github.com/protolambda/zcli
RUN wget https://raw.githubusercontent.com/ethereum/eth2.0-pm/master/interop/mocked_start/keygen_10000_validators.yaml

RUN ln -s /root/go/bin/zcli /usr/bin/zcli

ENTRYPOINT ["/bin/bash"]
