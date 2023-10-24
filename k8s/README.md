# K8s

## Installation k3s

### master node

```
echo -e "192.168.1.131\tpimaster" | sudo tee -a /etc/hosts
echo -e "192.168.1.130\tpiworker1" | sudo tee -a /etc/hosts
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

## Install NATS

```
helm repo add nats https://nats-io.github.io/k8s/helm/charts/
helm install nats nats/nats --values nats-values.yaml
```