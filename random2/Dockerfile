FROM --platform=$TARGETPLATFORM ubuntu:20.04
WORKDIR .
COPY bin/random-scheduler /usr/local/bin
RUN apt-get update && apt-get install -y curl
CMD ["random-scheduler"]