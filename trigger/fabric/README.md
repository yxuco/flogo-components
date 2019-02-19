
# 	fabric-invoke Trigger

## Fix flogo-web

To use this trigger with flogo-web UI, you need to fix the flogo-web docker image as follows:

First, start flogo-web, i.e.,
```
docker run -it -p 3303:3303 flogo/flogo-docker eula-accept
```
Then, in a second terminal, find the docker container and fix it as follows:
```
docker exec -it $(docker ps | grep "flogo/flogo-docker" | awk '{print $1}') bash
apk add --no-cache musl-dev
cd /tmp/flogo-web/build/server/local/engines/flogo-web/src/flogo-web
rm -Rf Gopkg.* vendor
export GOPATH=/tmp/flogo-web/build/server/local/engines/flogo-web
dep init
go build -o ../../bin/flogo-web
```
Finally, you can save the image for later use from another terminal, i.e.,
```
docker commit $(docker ps | grep "flogo/flogo-docker" | awk '{print $1}') new/flogo-docker
```