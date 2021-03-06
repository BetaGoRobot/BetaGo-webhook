FROM alpine as builder

COPY * /data/
ENV GOPROXY https://goproxy.io,direct

WORKDIR /data/


RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories && \
    apk update && \
    apk add go && \
    go build -o betaGoHook *.go 

FROM alpine as runner

COPY --from=builder  /data/ /data/ 

WORKDIR /data/


CMD ["./betaGoHook"]