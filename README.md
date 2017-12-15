# MariaDB Broker

This is an implementation of a Service Broker that uses Helm to provision
instances of [MariaDB](https://kubeapps.com/charts/stable/mariadb). This is a
**proof-of-concept** for the [Kubernetes Service
Catalog](https://github.com/kubernetes-incubator/service-catalog), and should not
be used in production.

[Learn more about the Kubernetes Service Catalog](https://medium.com/@prydonius/service-catalog-in-kubernetes-78c0736e3910)

**NOTE**: Works on Kubernetes v1.8.3 + Service Catalog v0.1.2.

## An existing MariaDB Server

Suppose you already have a MariaDB Server outside the Kubernetes cluster.

```
docker run -d --name=mariadb-server -p 3306:3306 -e MYSQL_ROOT_PASSWORD=passw0rd mariadb:10.1.16
```

## Installing the Broker

```
docker run -d --name=mariadb-broker -p 8005:8005 -e MARIADB_HOST=10.10.25.250 -e MARIADB_PORT=3306 -e MARIADB_USER=root -e MARIADB_PASS=passw0rd siji/mariadb-broker:1.0.0
```

Or using Helm chart:

```
helm install --name=mariadb-broker charts/mariadb-broker
```


## Create MariaDB broker in Kubernetes cluster

Please update the Broker URL in `examples/mariadb-broker.yaml`.

```
kubectl apply -f examples/1.mariadb-broker.yaml
```


## Create a MariaDB instance for JPress blog system

```
kubectl apply -f examples/2.jpress/2.1.jpress-instance.yaml

kubectl apply -f examples/2.jpress/2.2.jpress-binding.yaml

kubectl apply -f examples/2.jpress/2.3.jpress-blog-system.yaml
```


## Create a MariaDB instance for Wordpress blog system

```
kubectl apply -f examples/3.wodpress/3.1.wordpress-instance.yaml

kubectl apply -f examples/3.wodpress/3.2.wordpress-binding.yaml

kubectl apply -f examples/3.wodpress/3.3.wordpress-blog-system.yaml
```
