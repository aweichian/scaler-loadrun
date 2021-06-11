#!/usr/bin/env bash
cd ..
GOOS="linux" GOARCH="amd64" CGO_ENABLED="0" go build main.go
cd docker
tag=`git log --pretty=oneline -1 | awk '{print $1}' | cut -c1-7`
echo $tag
cp ../main .

docker build -t scaler-loadrun:$tag .

# docker tag scaler-loadrun:$tag ic-harbor.baozun.com/ic/scaler-loadrun:$tag
# docker push ic-harbor.baozun.com/ic/scaler-loadrun:$tag

