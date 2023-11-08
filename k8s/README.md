# K8s

## Installation k3s

### master node

```
echo -e "192.168.1.131\tpimaster" | sudo tee -a /etc/hosts
echo -e "192.168.1.130\tpiworker1" | sudo tee -a /etc/hosts
echo -e "192.168.1.129\tpiworker2" | sudo tee -a /etc/hosts
ssh pi@pimaster
sudo sed -i '$ s/$/ cgroup_memory=1 cgroup_enable=memory/g' /boot/firmware/cmdline.txt
curl -sfL https://get.k3s.io | sh -
sudo kubectl get nodes
sudo cat /var/lib/rancher/k3s/server/node-token
```

### worker node

```
ssh pi@piworker1
curl -sfL http://get.k3s.io | K3S_URL=https://192.168.1.131:6443 \
    K3S_TOKEN=join_token_we_copied_earlier sh -
```

### configure kubeconfig

In the master node, copy content:

```
sudo cat /etc/rancher/k3s/k3s.yaml
```

In local machine, paste content copied before:

```
mkdir ~/.kube
vi ~/.kube/config_k3spi
```

Change:

```
server: https://localhost:6443
```

To be:

```
server: https://pimaster:6443
```

After:

```
export KUBECONFIG=~/.kube/config_k3spi
kubectl get nodes
```

## Install metallb

```
kubectl apply -f https://raw.githubusercontent.com/metallb/metallb/v0.13.12/config/manifests/metallb-native.yaml
```

## Install certificate manager

```
helm repo add jetstack https://charts.jetstack.io
helm repo update
helm install \
  cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version v1.13.2 \
  --set installCRDs=true
```

## Install NATS

```
helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm install nats nats/nats --values nats-values.yaml
```

## Install Prometheus/Grafana

Install

```
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm upgrade --install kube-prometheus-stack prometheus-community/kube-prometheus-stack -n kube-prometheus --create-namespace --values prom-values.yaml
k apply -f prom-ingress.yaml
```

Password to access to Grafana (user is `admin`)

```
kubectl get secret --namespace kube-prometheus kube-prometheus-stack-grafana -o jsonpath="{.data.admin-password}" | base64 --decode ; echo
```

To see all possible chart values:

```
helm show values prometheus-community/kube-prometheus-stack
```

## Install the ingestion service

```
helm upgrade --install -n ingestion-service --create-namespace ingestion-service grow/ingestion-service --version 0.1.1
```

## Install the router config controller

```
./router-config-secret.sh
helm upgrade --install -n router-config --create-namespace router-config-controller grow/router-config-controller --version 0.1.0 --set 'controllerManager.manager.image.repository=ghcr.io/renato0307/grow-router-config-controller' --set 'controllerManager.manager.image.tag=v0.2.1'
```