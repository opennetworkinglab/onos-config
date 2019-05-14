# Deploying onos-config

One of the goals of the onos-config project is to provide simple deployment options
that integrate with modern technologies. Deployment configurations can be found in
the `/deployments` folder in this repository.

## Deploying on Kubernetes with Helm

[Helm] is a package manager for [Kubernetes] that allows projects to provide a
collection of templates for all the resources needed to deploy on k8s. ONOS Config
provides a Helm chart for deploying a cluster for development and testing. In the
future, this chart will be extended for production use.

### Resources

The Helm chart provides resources for deploying the config service and accessing
it over the network, both inside and outside the k8s cluster:
* `Deployment` - Provides a template for ONOS Config pods
* `ConfigMap` - Provides test configurations for the application
* `Service` - Exposes ONOS Config to other applications on the network
* `Secret` - Provides TLS certificates for end-to-end encryption
* `Ingress` - Optionally provides support for external load balancing

### Local Deployment Setup

To deploy the Helm chart locally, install [Minikube] and [Helm]. On OSX, this can be done
using [Brew]:

```bash
> brew install minikube
> brew install helm
```

On Linux, users have [additional options](https://kubernetes.io/docs/setup/minikube/#additional-links)
for installing local k8s clusters.

Once Minikube has been installed, start it with `minikube start`. If deploying an
[ingress] to access the service from outside the cluster, be sure to give enough
memory to the VM to run [NGINX].

```bash
> minikube start --memory=4096 --disk-size=50g --cpus=4
```

Once Minikube has been started, set your Docker environment to the Minikube Docker
daemon and build the ONOS Config image:

```bash
> eval $(minikube docker-env)
> make
```

Helm requires a special pod - called Tiller - to be running inside the k8s cluster for deployment
management. Deploy the Tiller pod to enable Helm for your cluster:

```bash
> helm init
```

For ingress, later versions of Minikube ship with [NGINX], so Minikube users simply
need to enable the ingress addon:

```bash
> minikube addons enable ingress
```

### Installing the Chart

To install the chart, simply run `helm install deployments/helm/onos-config` from]
the root directory of this project:

```bash
> helm install deployments/helm/onos-config
```

You can optionally enable [ingress] by overriding `ingress.enabled`. Note that you
must have an ingress controller installed/enabled:

```bash
> helm install \
    -n onos-config \
    --set ingress.enabled=true \
    deployments/helm/onos-config
```

The ingress controller uses the self-signed certificates that ship with the chart
to provide end-to-end routing, load balancing, and encryption, making the onos-config
services accessible from outside the k8s cluster. The default certificates expect the
service to be reached through the `config.onosproject.org` domain. Thus, to connect
to the service through the ingress, you must configure `/etc/hosts` to add the
load balancer's IP:

```
192.168.99.102 config.onosproject.org
```

The IP address of the ingress differs depending on the environment. In Minikube,
the ingress can be reached through the Minikube IP:

```bash
LBIP=$(minikube ip)
```

In clustered environments, the ingress IP can be retrieved from the ingress
metadata:

```bash
> kubectl get ingress
NAME                                      HOSTS                    ADDRESS     PORTS     AGE
onos-config-onos-config-manager-ingress   config.onosproject.org   10.0.2.15   80, 443   76m
```

Once you've determined the ingress IP, use the Helm chart certificates to connect
to the service through the load balancer:

```bash
> go run github.com/onosproject/onos-config/cmd/diags/changes \
    -address=config.onosproject.org:443 \
    -keyPath=deployments/helm/onos-config/files/certs/tls.key \
    -certPath=deployments/helm/onos-config/files/certs/tls.crt
```

The ingress routes requests based on the host header and redirects the HTTP/2
traffic to provide end-to-end encryption.

### Deploying the device simulator

onos-config provides a [device simulator](../tools/test/devicesim/gnmi_user_manual.md)
for end-to-end testing. As with the onos-config app, a [Helm] chart is provided for
deployment in [Kubernetes]. Each chart instance deploys a single simulated device
`Pod` and a `Service` through which the simulator can be accessed. The onos-config chart can
then be configured to connect to the devices in k8s.

To deploy a device, install the `deployments/helm/device-simulator` chart:

```bash
> helm install -n device-1 deployments/helm/device-simulator
NAME:   device-1
LAST DEPLOYED: Sun May 12 01:16:41 2019
NAMESPACE: default
STATUS: DEPLOYED

RESOURCES:
==> v1/ConfigMap
NAME                              DATA  AGE
device-1-device-simulator-config  1     1s

==> v1/Service
NAME                       TYPE       CLUSTER-IP     EXTERNAL-IP  PORT(S)    AGE
device-1-device-simulator  ClusterIP  10.110.252.69  <none>       10161/TCP  1s

==> v1/Pod
NAME                       READY  STATUS             RESTARTS  AGE
device-1-device-simulator  0/1    ContainerCreating  0         1s
```

onos-config pods can be connected to the device through the `Service` that's
created by the chart. This is done by simply passing a list of `devices` to the
config manager when deploying the Helm chart:

```bash
> helm install \
    -n onos-config \
    --set ingress.enabled=true \
    --set devices='{device-1-device-simulator}' \
    deployments/helm/onos-config
```

To deploy onos-config with multiple simulators, simply install the simulator
chart _n_ times to create _n_ devices:

```bash
> helm install -n device-1 deployments/helm/device-simulator
> helm install -n device-2 deployments/helm/device-simulator
> helm install -n device-3 deployments/helm/device-simulator
```

Then pass a comma-separated list of device services to the onos-config chart:

```bash
> helm install \
    -n onos-config \
    --set ingress.enabled=true \
    --set devices='{device-1-device-simulator,device-2-device-simulator,device-3-device-simulator}' \
    deployments/helm/onos-config
```

[Brew]: https://brew.sh/
[Helm]: https://helm.sh/
[Kubernetes]: https://kubernetes.io/
[k8s]: https://kubernetes.io/
[Minikube]: https://kubernetes.io/docs/setup/minikube/
[NGINX]: https://www.nginx.com/
[ingress]: https://kubernetes.io/docs/concepts/services-networking/ingress/
