FROM golang as build

ENV PROJECT json-formatter

WORKDIR /build

ADD . ./

RUN GOOS=js GOARCH=wasm go build -o bin/main.wasm $PROJECT/src/wasm

RUN CGO_ENABLED=0 GOOS=linux go build -o json-formatter src/main.go

RUN cat /etc/passwd | grep nobody > passwd.nobody


FROM scratch

WORKDIR /var/www

COPY --from=build /build/passwd.nobody /etc/passwd
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/local/go/misc/wasm/wasm_exec.js ./

COPY --from=build /build/bin/* ./

ADD src/public/*.html ./
ADD src/public/style.css ./
ADD src/public/load_wasm.js ./

COPY --from=build /build/json-formatter ./

USER nobody
EXPOSE 3000

CMD ["./json-formatter", "/var/www"]