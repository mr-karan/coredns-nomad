FROM golang:1.21-alpine
RUN apk update && apk add make git
WORKDIR /app

RUN git clone https://github.com/coredns/coredns
RUN echo "nomad:https://github.com/mr-karan/coredns-nomad" >> coredns/plugin.cfg
WORKDIR /app/coredns
RUN go mod download

RUN make 
RUN mv coredns /coredns

FROM scratch
WORKDIR /
COPY --from=0 /coredns /

EXPOSE 53

ENTRYPOINT ["/coredns"]