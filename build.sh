#!/bin/bash
tag="source-go-http-$(date +%s%N)"
echo $tag
docker build . -t $REPOSITORY:${tag}
docker push $REPOSITORY:${tag}
echo "pushed ${tag}"