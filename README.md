# firaaq
A tiny educational-purpose project to create containers, written in Go.

It basically is a tiny version of docker, it uses neither [containerd](https://containerd.io/) nor [runc](https://github.com/opencontainers/runc). Only a set of the Linux features.

## Features
Vessel supports:
* __Control Groups__ for resource restriction (CPU, Memory, Swap, PIDs)
* __Namespace__ for global system resources isolation (Mount, UTS, Network, IPS, PID)
* __Union File System__ for branches to be overlaid in a single coherent file system. (OverlayFS)

## Read more
Here is the list of blog posts about container implementation:

1. [Build Containers From Scratch in Go (Part 1: Namespaces)](https://alijosie.medium.com/build-containers-from-scratch-in-go-part-1-namespaces-c07d2291038b)
2. To be continued...

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

    firaaq run alpine /bin/sh
    firaaq run alpine # same as above due to alpine default command

Restart Nginx service inside a container with ID: 123456789123

    firaaq exec 1234567879123 systemctrl restart nginx
    
List running containers

    firaaq ps
    
List local images

    firaaq images

## Notice
firaaq, obviously, is not a production ready container manager tool. 
