# Deploy Mastro to a K8s Cluster
Mastro is a stateless service that can be easily deployed to a K8s cluster.

In the examples below, we assume a previously deployed mongo database available on the same namespace or any reachable host at `mongo-mongodb:27017`.

For instance, we used the one using a StatefulSet and deployed as a Helm chart provided by bitnami (see [here](https://bitnami.com/stack/mongodb/helm)).

## Catalogue

### Config Map

The config for the catalogue can be defined as a K8s config map, as follows:

```
apiVersion: v1
data:
  catalogue-conf.yaml: |
    type: catalogue
    details:
      port: 8085
    backend:
      name: catalogue-mongo
      type: mongo
      settings:
        database: mastro
        collection: mastro-catalogue
        connection-string: "mongodb://mastro:mastro@mongo-mongodb:27017/mastro"
kind: ConfigMap
metadata:
  name: catalogue-conf
```

Mind that in the example above we specified directly the DB user and password (i.e., `mastro:mastro`).
A K8s secret or one injected by an external vault (e.g. hashicorp) can be used for this purpose.

### Deployment

A deployment can be created to spawn multiple replicas for the catalogue.

The configuration is mounted as volume and its path set using the MASTRO_CONFIG variable.

```
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: mastro-catalogue
  name: mastro-catalogue
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mastro-catalogue
  strategy: {}
  template:
    metadata:
      labels:
        app: mastro-catalogue
    spec:
      containers:
      - image: pilillo/mastro-catalogue:20210306-static
        imagePullPolicy: Always
        name: mastro-catalogue
        resources: {}
        ports:
        - containerPort: 8085
          protocol: TCP
        env:
        - name: MASTRO_CONFIG
          value: /conf/catalogue-conf.yaml
        volumeMounts:
        - mountPath: /conf
          name: catalogue-conf-volume
      securityContext: {}
      volumes:
      - name: catalogue-conf-volume
        configMap:
          defaultMode: 420
          name: catalogue-conf
```

### Service

A service is created with:

```
apiVersion: v1
kind: Service
metadata:
  labels:
    app: mastro-catalogue
  name: mastro-catalogue
spec:
  ports:
  - name: rest-8085
    port: 8085
    protocol: TCP
    targetPort: 8085
  selector:
    app: mastro-catalogue
  type: ClusterIP
```

Mind that the service only exposes the catalogue across the namespace.

You will have to create an ingress or a route (respectively on plain K8s and openshift) to make it reachable from the outside world.

## Feature Store

### Config Map

```
apiVersion: v1
data:
  fs-conf.yaml: |
    type: featurestore
    details:
      port: 8085
    backend:
      name: fs-mongo
      type: mongo
      settings:
        database: mastro
        collection: mastro-featurestore
        connection-string: "mongodb://mastro:mastro@mongo-mongodb:27017/mastro"
kind: ConfigMap
metadata:
  creationTimestamp: null
  name: fs-conf
```

### Deployment

```
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: mastro-featurestore
  name: mastro-featurestore
spec:
  replicas: 1
  selector:
    matchLabels:
      app: mastro-featurestore
  strategy: {}
  template:
    metadata:
      labels:
        app: mastro-featurestore
    spec:
      containers:
      - image: pilillo/mastro-featurestore:20210306-static
        imagePullPolicy: Always
        name: mastro-featurestore
        resources: {}
        ports:
        - containerPort: 8085
          protocol: TCP
        env:
        - name: MASTRO_CONFIG
          value: /conf/fs-conf.yaml
        volumeMounts:
        - mountPath: /conf
          name: fs-conf-volume
      securityContext: {}
      volumes:
      - name: fs-conf-volume
        configMap:
          defaultMode: 420
          name: fs-conf
```

### Service

```
apiVersion: v1
kind: Service
metadata:
  labels:
    app: mastro-featurestore
  name: mastro-featurestore
spec:
  ports:
  - name: rest-8085
    port: 8085
    protocol: TCP
    targetPort: 8085
  selector:
    app: mastro-featurestore
  type: ClusterIP
```
