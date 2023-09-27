---
# Distributed System: Bridging eCommerce Giants to Merchants 🌐

Welcome to Distributed System! 🎉

Distributed System is an innovative Free Open Source project designed to seamlessly integrate the vast ecosystems of Amazon and eBay with online merchant platforms. This pioneering RESTful API allows merchants to effortlessly manage orders from these prominent sales channels and synchronize product uploads, opening up new horizons in e-commerce connectivity.

Dive into a robust and scalable architecture that embraces modern technology stack, deploying microservices in a distributed environment orchestrated by Kubernetes and balanced by NGINX. Each component is meticulously containerized using Docker, ensuring isolated, consistent, and efficient environments for each service.

🚀 **Features:**
- **Multi-Platform Integration:** Connects Amazon and eBay with online merchant websites to manage orders and upload products.
- **Microservices Architecture:** Ensures modularity, scalability, and maintainability of the application.
- **Distributed Deployment:** Utilizes NGINX load balancer for optimal distribution of network or application traffic across multiple servers.
- **Containerization:** Leverages Docker for encapsulating the application and its dependencies into a single object.
- **Orchestration:** Managed through Kubernetes, enabling automated deployment, scaling, and management of containerized applications.
- **Messaging System:** Implements Apache Kafka for building real-time data pipelines and streaming apps.
- **Cloud Compatibility:** Tested on Linode and DigitalOcean Kubernetes clusters and has plans for testing on Azure, GCP, and AWS to ensure broad applicability and versatility.

🔬 **Testing Environments:**
- Linode Kubernetes Clusters
- DigitalOcean Kubernetes Clusters
- *Upcoming:* Azure, GCP, AWS

Whether you are an e-commerce enthusiast, a developer with a penchant for distributed systems, or simply curious about integrating diverse e-commerce platforms, Distributed System invites you to explore, contribute, and innovate in the world of e-commerce connectivity. Let's build a more connected, efficient, and versatile e-commerce ecosystem together! 🌟

Stay tuned and feel free to dive into the code, raise issues, submit features, and be a part of this exciting journey in reshaping the e-commerce landscape!
---

## Overview

## Prepare Docker containers and upload it on DockerHub
1. Follow the instructions given in https://docs.docker.com/engine/install/ to install docker in your PC/Laptop.

3. Command to build the docker image. Make sure you are in the same folder/directory of the files
```console
docker image build -t rfernandohub/ob:{tag_name} .
```

3. Push your docker image. You need to have an account in https://hub.docker.com/
```console
docker image push rfernandohub/ob:{tag_name}
```

You can use the docker image I have created for your deployments: https://hub.docker.com/repository/docker/rfernandohub/ob/general

## Configure Kubernetes Cluster

1. Create Kubernetes cluster in Linode and Digital Ocean

2. Install Kubectl
```console
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"

chmod +x kubectl

sudo mv ./kubectl /usr/local/bin/kubectl

kubectl version --client
```

3. Apply kubeconfig.yaml file provided by the cloud providers on your PC/Laptop to control Kubernetes 

```console
nano do-kubeconfig.yaml
nano linode-kubeconfig.yaml

export KUBECONFIG=do-kubeconfig.yaml
export KUBECONFIG=linode-kubeconfig.yaml
```

4. To create a service in your Kubernetes cluster, create a YAML file and apply the file using the following command.
```console
kubectl apply -f rfernandogo_service_do.yaml
```

5. How to see the current pods in your Kubernetes Cluster 

```console
kubectl get pods
```

6. How to see the current services in your Kubernetes Cluster 
```console
kubectl get services
```

7. How to details about the nodes in your Kubernetes Cluster
```console
kubectl get nodes -o wide
```

8. If you need to edit the deployment, you can use the following command
```console
kubectl edit deployment web-deploy
```

## Install NGINX Load Balancer

1. Log into the ubuntu sever
```console
ssh root@xx.xxx.xxx.xxx
```

2. Update all packages and install NGINX
```console
sudo apt update -y
sudo apt install nginx -y
```

3. Edit the default NGINX file
```console
cd /etc/nginx/sites-available
sudo nano default
```

4. Remove the existing code and add the following.

```yaml
upstream myproxy {
    ip_hash;
    server aaa.bb.ccc.dd:8080;
    server eee.ff.ggg.hhh:8080;
}

server {
    listen 80;
    server_name localhost;
    root /var/www/html;

    location / {
        proxy_pass http://myproxy;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header Host $host;
        proxy_redirect off;
    }
}
```

5. Reload the NGINX server
```console
sudo systemctl reload nginx
```


## Install Prometheus and Grafana for monitoring the services


```console
nano {service-provider}-kubeconfig.yaml 
export KUBECONFIG={service-provider}-kubeconfig.yaml
```

helm repo add prometheus-community https://prometheus-community.github.io/helm-charts

1. Follow the instructions to intall HELM in your Kubernetes Cluster https://helm.sh/docs/intro/install/. You need to make sure your device(PC/Laptop) is connected to the cluster via Kubectl

```console
curl https://baltocdn.com/helm/signing.asc | gpg --dearmor | sudo tee /usr/share/keyrings/helm.gpg > /dev/null

sudo apt-get install apt-transport-https --yes

echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/helm.gpg] https://baltocdn.com/helm/stable/debian/ all main" | sudo tee /etc/apt/sources.list.d/helm-stable-debian.list

sudo apt-get update

sudo apt-get install helm
```

2. Add the prometheus-community Repo
```console
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
```

3. Install Prometheus 
```console
helm install prometheus prometheus-community/prometheus
```

4. Expose prometheus server as a NodePort to the public
```console
kubectl expose service prometheus-server --type=NodePort --target-port=9090 --name=prometheus-server-ext
```

5. Install Grafna
```console
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
helm install grafana grafana/grafana
```

5. To reveal Grafna username password run the following
```console
kubectl get secret --namespace default grafana -o yaml
```

6. To decode the username password you can run the following 
```console
echo "{username_string}" | openssl base64 -d ; echo
# Here for example username_string is "YWRtaW4=" and its equals to "admin"
echo "YWRtaW4=" | openssl base64 -d ; echo

# Do the same for password 
echo "{password_string}" | openssl base64 -d ; echo
```

7. Expose Grafana server as a NodePort to the public
```console
kubectl expose service grafana --type=NodePort --target-port=3000 --name=grafana-ext
```

Normal Prometheus Dashboard:
http://172.104.128.242:31161/graph?g0.expr=node_memory_Active_bytes&g0.tab=0&g0.stacked=0&g0.show_exemplars=0&g0.range_input=1h

Grafna Dashboard:
http://172.104.128.242:30276/d/k8s_views_ns/kubernetes-views-namespaces?orgId=1&refresh=30s


## How to configure APache Kafka

1. Install Java
```console
sudo apt install default-jre
```


2. Install Scala with cs setup (recommended)
```console
sudo apt-get install scala
```


3. Install Kafka
Make sure it will be same as scala version 2.11.x in your previous step
```console
wget https://archive.apache.org/dist/kafka/2.4.1/kafka_2.11-2.4.1.tgz
```


4. Run the following commands in order to start all services in the correct order:
```console
# Start the ZooKeeper service
bin/zookeeper-server-start.sh config/zookeeper.properties

# Start the Kafka broker service
bin/kafka-server-start.sh config/server.properties
```

5. To create a topic, Open another terminal session and run:
```console
bin/kafka-topics.sh --create --topic openbaypro --zookeeper localhost:2181 --partitions 1 --replication-factor 1
```

6. To see the topic list:
```console
bin/kafka-topics.sh --list --zookeeper localhost:2181
```

7. To write events to the topic:
```console
bin/kafka-console-producer.sh --broker-list localhost:9092 --topic openbaypro
```

8. To read the events published through producer:
```console
bin/kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic openbaypro --from-beginning
```

To read and produce events, we can use the Golang sripts in the Kafka folder mentioned in this repo.
