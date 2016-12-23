# dockfs

A fuse userspace filesystem to interact with docker



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
        * stdout: container standard output
        * stderr: container standard error
    * byState
      * {state}
        * {container id}: symlink to container in the state described by the parent node
    * byName
      * {container name}: symlink to the container
  * images: images
    * byId
      * {image id}
        * json:	informations coming from the docker api
        * parent: if present, symbolic link to the image parent
        * containers
          * {container id}: symlink to container started from the image
    * byTag
      * {repository} or {source/repository}
        * {tag}: symlink to theimage 
  * networks: docker networks
    * byId
      * {network id}
        * json: informations coming from the docker api
        * name: network name
        * containers: 
          * {container id}: symlink to container connected to the docker network
	
# External dependencies

This project depends directly on the following projects, I thanks to the authors

* [bazil.org/fuse](https://github.com/bazil/fuse)

* [docker/client](https://github.com/docker/docker/tree/master/client)

Before using it please check license compatibility for your use cases

