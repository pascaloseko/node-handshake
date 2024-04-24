FROM ubuntu:22.04

WORKDIR /node-app

COPY bitcoin.conf /node-app/.bitcoin/bitcoin.conf

COPY . /node-app/

EXPOSE 8333

CMD ["./bin/bitcoind", "-conf=/node-app/.bitcoin/bitcoin.conf"]