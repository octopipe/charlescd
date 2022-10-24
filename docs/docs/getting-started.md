# Getting Started

## Requirements

- You have a cluster that uses Kubernetes v1.22 or newer
- Installed [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) command-line tool.

## Install Charles CD

1.Install core components of CharlesCD by running the command:

```bash
kubectl create namespace charlescd
kubectl apply -f https://github.com/octopipe/charlescd/releases/download/charlescd-0.0.1/install.yaml -n charlescd
```

!!! not "Component by component"

    If you want to install compont by component, go to the [release page]() and see all installation possibilities


## Install a network layer

The following tabs expand to show instructions for installing a networking layer.
Follow the procedure for the networking layer of your choice:

=== "Gate - The CharlesCD component (WIP)"

    WIP

=== "Istio"

    Install a properly configured Istio by following the [Istio installation](https://istio.io/latest/docs/setup/install/)


## Verify the installation

```bash
kubectl get pods -n charlescd
```

WIP