# dockfs

A fuse userspace filesytem to interact with docker



# Requirements

this utility requires the following software installed locally on machine where it is executed:
* docker engine (in case of local connection)
* fuse package

in additions, user that execute the command must have access to the mountpoint



# How to build

```
go get "github.com/mysinmyc/dockfs/cmd/dockfs"
```



# How to execute

```
mkdir $HOME/dockfs
$GOPATH/bin/dockfs -mountPoint $HOME/dockfs &
```



#Implemented nodes

* {mountPoint}
  * containers: containers	
    * byId	
      * {container id}
        * json:	informations coming from the docker api
        * name:	contaner name
        * command: command
        * image: symbolic link to the image
  * images: images
    * by id
      * {image id}
        * json:	informations coming from the docker api
        * parent: if present, symbolink link to the image parent


	
# Extenal dependencies

this project depends directly on

* [bazil.org/fuse](https://github.com/bazil/fuse)

* [docker/client](https://github.com/docker/docker/tree/master/client)

Before using it please check license compatibility for your use cases



## Known issues

It can happens that docker api version pulled from github is newer than the version of the local docker engine. If it happens, download the right version or try with the following environment variable if the two releases are compatible

```
export DOCKER_API_VERSION={expected version}
```
