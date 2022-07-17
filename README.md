# GHost

Game server hosting service

## What is this?

This is a personal project of mine, that I made when I wanted to learn how big game hosting services operate, and I wanted a service that would make running a game server that easy, but let me run it on my own hardware.

## How does it work?

GHost is built on [Kubernetes](https://k8s.io), which allows you to network together many computers to run applications on, unlike similar solutions like docker and podman, which are built to run on a single machine

The servers themselves are run as container images that are deployed to a Kubernetes cluster using API calls. Using the Kubernetes Go client, we can customize these API calls based on user input to define CPU, memory, and storage limits, choose what kind of server to run, and lots of other things that I just haven't gotten around to implementing yet.

## How can I run this myself?

While this is by no means production ready, setting up a development environment is quite simple. All you need is a Kubernetes cluster, and a computer that can run Go

Before you start, you need to make sure you have a storageclass set up on your cluster. The default option for GHost is [Longhorn](https://longhorn.io/docs/latest/deploy/install/install-with-kubectl/). Check the docs for full instructions, but the quick setup method is as follows:

```sh
curl -sSfL https://raw.githubusercontent.com/longhorn/longhorn/v1.3.0/scripts/environment_check.sh | bash
```

This checks for any missing dependencies on your cluster. Fix any issues the script finds, then install longhorn on the cluster:

```sh
kubectl apply -f https://raw.githubusercontent.com/longhorn/longhorn/v1.3.0/deploy/longhorn.yaml
```
After that, setting up the dev environment itself is quite simple.

1. Download the Go programming language. I develop on the latest version (1.18)
2. Make sure you have your cluster's kubeconfig on your development machine. on Linux, it should be at ``~/.kube/config``. I don't use Windows so I'm not sure, but I'd imagine it's something similar.
3. After you've cloned the repo, navigate to it in a terminal or IDE, and run ```go run .```

## Development Updates and TODO

This version of GHost is officially a direct upgrade from the original version, so I've been able to begin developing new features. Right now I am adding the management interface. To start with, that will probably look like a WebSocket connection to a Kubernetes exec call to the pod, but I will probably eventually add more customized settings for each game with its respective configuration files and options.

After some form of management is added and functional on the existing servers, I will add actual authentication and some user account management.

From there, we'll see where it goes. Some things I'm thinking about: 
- Admin console
- Begin using configuration files instead of magic constants for some options
- Make use of an actual database for server information instead of repeated and redundant API calls

---
This is a complete rewrite of a python project of the same name, archived at https://github.com/Quinn-5/GHost-python
