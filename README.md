# Calzone

_Delicious App Bundles_

## Overview

Calzone is **super easy**, **super ux-friendly** app management system for bundling and deploying anything. The author (@warent) created this after not quite getting the desired experience from existing solutions.

- docker-compose is too transient
- Standalone docker is too complicated to wire up multiple dependencies
- Using no VM requires tons of annoying manual configurations and dependencies
- microk8s / minikube / etc. requires too much custom, heavy, repetitive configuration
- helm provides a terrible UX

## Requirements

- [Docker](https://docs.docker.com/get-docker/)
- [Sysbox Docker Runtime](https://github.com/nestybox/sysbox)
