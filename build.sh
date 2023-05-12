docker run \
    --rm \
    -w /tgrobot \
    -v $(pwd):/tgrobot\
    -v $(go env GOPATH):/go \
    -e GOCACHE=/tmp \
    -u 1000:1000 \
    golang:1.18.10-buster go build
