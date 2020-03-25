<p align="center">
    <img alt="Kubernets Logo" src="https://cdn.worldvectorlogo.com/logos/kubernets.svg" height="150" />
    <h2 align="center">Pods</h2>
</p>

Get a full YAML of a deployed pod

```
$ kubectl get pods $podname -o yaml
$ kubectl get pods $podname -o json
```

A simple YAML descriptor for a pod

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: koala
  labels:
    app.kubernetes.io/name: mysql
    app.kubernetes.io/instance: wordpress-abcxzy
    app.kubernetes.io/version: "5.7.21"
    app.kubernetes.io/component: database
    app.kubernetes.io/part-of: wordpress
    app.kubernetes.io/managed-by: helm
spec:
  containers:
    -
      image: clivern/koala
      name: front
      ports:
        -
          containerPort: 8080
          protocol: TCP
```

To create a pod from yaml file

```
$ kubectl create -f koala.yaml
```

To get all pods

```
$ kubectl get pods
$ kubectl get pods --show-labels
$ kubectl get pods -L app.kubernetes.io/name,app.kubernetes.io/version

# Add label to pod
$ kubectl label pods ${podName} app.kubernetes.io/author=clivern

# Overwrite a pod label
$ kubectl label pods ${podName} app.kubernetes.io/author=clivern --overwrite
```

To get pod logs

```
$ kubectl logs ${podName}
$ kubectl logs ${podName} -c ${containerName}

# Obtaining the application log of a crashed container "you want to figure out why the previous container terminated"
$ kubectl logs ${podName} --previous
```

To talk to a pod without going through a service, you can forward a local port to a port of the pod

```
# this will forward local port 8000 to port 8080 of a pod
$ kubectl port-forward ${podName} 8000:8080
```

Listing pods using a label selector

```
$ kubectl get pods -l environment=production,tier=frontend
$ kubectl get pods -l 'environment in (production),tier in (frontend)'
$ kubectl get pods -l 'environment in (production, qa)'
$ kubectl get pods -l 'environment,environment notin (frontend)'

// Field selectors let you select Kubernetes resources based on the value of one or more resource fields.
$ kubectl get pods --field-selector status.phase=Running
$ kubectl get pods --field-selector metadata.name=my-service
$ kubectl get pods --field-selector metadata.namespace!=default
$ kubectl get pods --field-selector status.phase=Pending
```

Labeling a node

```
$ kubectl get nodes
$ kubectl label node ${nodeName} gpu=true

# show nodes with gpu is true
$ kubectl get nodes -l gpu=true

# show nodes with gpu label
$ kubectl get nodes -L gpu
```

Scheduling pods to a specific nodes

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  labels:
    app.kubernetes.io/component: database
    app.kubernetes.io/instance: wordpress-abcxzy
    app.kubernetes.io/managed-by: helm
    app.kubernetes.io/name: mysql
    app.kubernetes.io/part-of: wordpress
    app.kubernetes.io/version: "5.7.21"
  name: koala
spec:
  containers:
    -
      image: clivern/koala
      name: front
      ports:
        -
          containerPort: 8080
          protocol: TCP
  nodeSelector:
    gpu: "true"
```

Annontations can contain large data (up to 256 KB in total) unlike labels.

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  annotations:
    app.kubernetes.io/created-by: "{\"kind\": \"..\": \"apiVersion\": \"1.0.0\", \"Ref\": \"....\"}"
  labels:
    app.kubernetes.io/component: database
    app.kubernetes.io/instance: wordpress-abcxzy
    app.kubernetes.io/managed-by: helm
    app.kubernetes.io/name: mysql
    app.kubernetes.io/part-of: wordpress
    app.kubernetes.io/version: "5.7.21"
  name: koala
spec:
  containers:
    -
      image: clivern/koala
      name: front
      ports:
        -
          containerPort: 8080
          protocol: TCP
  nodeSelector:
    gpu: "true"
```

```
$ kubectl annotate pod ${podName} app.kubernetes.io/created-by="foo bar"
```

Namespaces used to separate objects so you can have the same resource names multiple times across different namespaces.

```
$ kubectl get namespaces
```

if you list pods without namespace, it will use the default namespace

```
$ kubectl get pods --namespace kuber-system
```

To create a namespace

```
apiVersion: v1
kind: Namespace
metadata:
  name: custom-namespace
```

```
$ kubectl create -f custom-namespace.yaml

// OR

$ kubectl create namespace custom-namespace
```

To create a resource on a namespace, either add a `namespace: custom-namespace` entry to `metadata`

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: koala
  namespace: custom-namespace
  labels:
    app.kubernetes.io/name: mysql
    app.kubernetes.io/instance: wordpress-abcxzy
    app.kubernetes.io/version: "5.7.21"
    app.kubernetes.io/component: database
    app.kubernetes.io/part-of: wordpress
    app.kubernetes.io/managed-by: helm
spec:
  containers:
    -
      image: clivern/koala
      name: front
      ports:
        -
          containerPort: 8080
          protocol: TCP
```

or specify a namespace when creating the resource

```
$ kubectl create -f koala.yaml -n custom-namespace
```

To delete a pod

```
$ kubectl delete pods ${podName}

# Delete all pods that has creation_method is manual
$ kubectl delete pods -l creation_method=manual
```

Deleting a namespace will delete all the pods attached automatically

```
$ kubectl delete namespace custom-namespace
```

To delete all pods on the current namespace

```
$ kubectl delete pods --all
```

To delete all pods, services and ReplicationController withing the current namespace

```
$ kubectl delete all --all
```


Kubernetes can check if a container is still alive through liveness probes. You can specify a liveness probe for each container in the pod’s specification. Kubernetes will periodi- cally execute the probe and restart the container if the probe fails.

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  labels:
    app.kubernetes.io/component: database
    app.kubernetes.io/instance: wordpress-abcxzy
    app.kubernetes.io/managed-by: helm
    app.kubernetes.io/name: mysql
    app.kubernetes.io/part-of: wordpress
    app.kubernetes.io/version: "5.7.21"
  name: koala
  namespace: custom-namespace
spec:
  containers:
    -
      image: clivern/koala
      livenessProbe:
        httpGet:
          path: /_health
          port: 8080
      name: front
      ports:
        -
          containerPort: 8080
          protocol: TCP
```

The `periodSeconds` field specifies that the kubelet should perform a liveness probe every 5 seconds. The `initialDelaySeconds` field tells the kubelet that it should wait 5 second before performing the first probe.

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  labels:
    app.kubernetes.io/component: database
    app.kubernetes.io/instance: wordpress-abcxzy
    app.kubernetes.io/managed-by: helm
    app.kubernetes.io/name: mysql
    app.kubernetes.io/part-of: wordpress
    app.kubernetes.io/version: "5.7.21"
  name: koala
  namespace: custom-namespace
spec:
  containers:
    -
      image: clivern/koala
      livenessProbe:
        exec:
          command:
            - cat
            - /tmp/healthy
        initialDelaySeconds: 5
        periodSeconds: 5
      name: front
      ports:
        -
          containerPort: 8080
          protocol: TCP
```

to define HTTP headers for `livenessProbe`

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  labels:
    app.kubernetes.io/component: database
    app.kubernetes.io/instance: wordpress-abcxzy
    app.kubernetes.io/managed-by: helm
    app.kubernetes.io/name: mysql
    app.kubernetes.io/part-of: wordpress
    app.kubernetes.io/version: "5.7.21"
  name: koala
  namespace: custom-namespace
spec:
  containers:
    -
      image: clivern/koala
      livenessProbe:
        httpGet:
          path: /_health
          port: 8080
        httpHeaders:
          -
            name: Custom-Header
            value: Awesome
        initialDelaySeconds: 3
        periodSeconds: 3
      name: front
      ports:
        -
          containerPort: 8080
          protocol: TCP
```

a TCP check is quite similar to an HTTP check.

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: goproxy
  labels:
    app: goproxy
spec:
  containers:
  - name: goproxy
    image: k8s.gcr.io/goproxy:0.1
    ports:
    - containerPort: 8080
    readinessProbe:
      tcpSocket:
        port: 8080
      initialDelaySeconds: 5
      periodSeconds: 10
    livenessProbe:
      tcpSocket:
        port: 8080
      initialDelaySeconds: 15
      periodSeconds: 20
```

You can use a named ContainerPort for HTTP or TCP liveness checks:

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  labels:
    app.kubernetes.io/component: database
    app.kubernetes.io/instance: wordpress-abcxzy
    app.kubernetes.io/managed-by: helm
    app.kubernetes.io/name: mysql
    app.kubernetes.io/part-of: wordpress
    app.kubernetes.io/version: "5.7.21"
  name: koala
  namespace: custom-namespace
spec:
  containers:
    -
      image: clivern/koala
      livenessProbe:
        httpGet:
          path: /_health
          port: liveness-port
        httpHeaders:
          -
            name: Custom-Header
            value: Awesome
        initialDelaySeconds: 3
        periodSeconds: 3
      name: front
      ports:
        -
          containerPort: 8080
          name: liveness-port
          protocol: TCP
```

Sometimes, applications are temporarily unable to serve traffic "might need to load large data or configuration files during startup". In such cases, you don’t want to kill the application, but you don’t want to send it requests either. Kubernetes provides `readiness probes` to detect and mitigate these situations. A pod with containers reporting that they are not ready does not receive traffic through Kubernetes Services.

Probes have a number of fields that you can use to more precisely control the behavior of liveness and readiness checks:

- `initialDelaySeconds`: Number of seconds after the container has started before liveness or readiness probes are initiated.
- `periodSeconds`: How often (in seconds) to perform the probe. Default to 10 seconds. Minimum value is 1.
- `timeoutSeconds`: Number of seconds after which the probe times out. Defaults to 1 second. Minimum value is 1.
- `successThreshold`: Minimum consecutive successes for the probe to be considered successful after having failed. Defaults to 1. Must be 1 for liveness. Minimum value is 1.
- `failureThreshold`: When a Pod starts and the probe fails, Kubernetes will try failureThreshold times before giving up. Giving up in case of liveness probe means restarting the Pod. In case of readiness probe the Pod will be marked Unready. Defaults to 3. Minimum value is 1.

HTTP probes have additional fields that can be set on `httpGet`:

- `host`: Host name to connect to, defaults to the pod IP. You probably want to set “Host” in httpHeaders instead.
- `scheme`: Scheme to use for connecting to the host (HTTP or HTTPS). Defaults to HTTP.
- `path`: Path to access on the HTTP server.
- `httpHeaders`: Custom headers to set in the request. HTTP allows repeated headers.
- `port`: Name or number of the port to access on the container. Number must be in the range 1 to 65535.
