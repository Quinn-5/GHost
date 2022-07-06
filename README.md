# GHost

Game server hosting service

## What is this?

This is a personal project of mine, that I made when I wanted to learn how big game hosting services operate, and I wanted to make a service that would make running a game server that easy, but let me run it on my own hardware.

## How does it work?

GHost is built on [Kubernetes](https://k8s.io), which allows you to network together many computers to run applications on, unlike similar solutions like docker and podman, which are built to run on a single machine

The servers themselves are run as container images that are deployed to a Kubernetes cluster using API calls. Using the Kubernetes Go client, we can customize these API calls based on user input to define cpu, memory, and storage limits, choose what kind of server to run, and lots of other things that I just haven't gotten around to implementing yet.

## How can I run this myself?

While this is by no means production ready, setting up a development environment is quite simple. All you need is a Kubernetes cluster, and a computer that can run Go

0. You need to have some type of storageclass set up on your cluster. By default this repo is set up to use longhorn, which is pretty plug and play in my experience. To set it up, there are two commands to run on your cluster

To check the environment is set up properly:

```curl -sSfL https://raw.githubusercontent.com/longhorn/longhorn/v1.3.0/scripts/environment_check.sh | bash```

Install any missing dependencies this finds, then install longhorn on the cluster:

```kubectl apply -f https://raw.githubusercontent.com/longhorn/longhorn/v1.3.0/deploy/longhorn.yaml```


1. Download the Go programming language. I develop on the latest version
2. Make sure you have your cluster's kubeconfig on your development machine. on linux, it should be at ``~/.kube/config``. I don't use Windows so I'm not sure, but I'd imagine it's something similar.
3. After you've cloned the repo, navigate to it in a terminal or IDE, and run ```go run .```

## Development Updates and TODO

I'm very close at this point to this being a drop-in replacement of the python version. I only need to find a way to grab the IP and port from the cluster, and also have it not be a really janky solution I'll have to fix later. As soon as I can fix that, I have some other things in mind I'll work on after, not necessarily in this order:

- Server view to view and delete all existing servers, (yes, it's still impossible to delete a server with this)
- Management interface to edit comfigurations and possibly (???) get console access to the pods running on the server
- User authentication and management, so multiple people can use this, and it can be run on an unsecure network
- Admin console to manage users, resources, app configurations

Oh my gosh, I have a lot of work to do

---
This is a complete rewrite of a python project of the same name, archived at https://github.com/Quinn-5/GHost-python
