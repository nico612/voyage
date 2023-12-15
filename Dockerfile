FROM golang:1.20.2 as builder

WORKDIR /app

ARG SERVICE
ARG BUILDPATH

ENV GO111MODULE=on
#ENV GOPROXY=https://goproxy.cndo

COPY . .

RUN sed -i 's#http://deb.debian.org#https://mirrors.163.com#g' /etc/apt/sources.list\
    && apt update && apt install tree protobuf-compiler -y \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28 \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2 \
    && go mod download \
    && go mod tidy

RUN CGO_ENABLED=0 go build -o ${SERVICE} ./${BUILDPATH}/main.go


FROM debian:buster-slim
RUN sed -i 's#http://deb.debian.org#http://mirrors.tuna.tsinghua.edu.cn#g' /etc/apt/sources.list

WORKDIR /app

ARG SERVICE
ENV SERVICE=${SERVICE}

COPY --from=builder /app/configs /app/configs
COPY --from=builder /app/${SERVICE} /app/${SERVICE}
EXPOSE 4560
CMD [ "/bin/sh", "-c", "./${SERVICE}", "--config", "configs/config.docker.yaml"]

#ENTRYPOINT [ "/bin/sh", "-c", "./${SERVICE}", "-e", "configs/config.docker.yaml"]
