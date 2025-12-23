# firaq
A tiny educational-purpose project to create containers, written in Go.

It basically is a tiny version of docker, it uses neither [containerd](https://containerd.io/) nor [runc](https://github.com/opencontainers/runc). Only a set of the Linux features.

## Features
firaq supports:
* __Control Groups__ for resource restriction (CPU, Memory, Swap, PIDs)
* __Namespace__ for global system resources isolation (Mount, UTS, Network, IPS, PID)
* __Union File System__ for branches to be overlaid in a single coherent file system. (OverlayFS)


## Install

    go get -u github.com/samama/firaaq
    
## Usage

    Usage:
      firaaq [command]
    
    Available Commands:
      exec        Run a command inside a existing Container.
      help        Help about any command
      images      List local images
      ps          List Containers
      run         Run a command inside a new Container.

## Examples

Run `/bin/sh` in `alpine:latest`

    firaq run alpine /bin/sh
    firaq run alpine # same as above due to alpine default command

Restart Nginx service inside a container with ID: 123456789123

    firaq exec 1234567879123 systemctrl restart nginx
    
List running containers

    firaq ps
    
List local images

    firaq images

## Notice
firaq, obviously, is not a production ready container manager tool. 
