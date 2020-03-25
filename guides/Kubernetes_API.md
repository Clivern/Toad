<p align="center">
    <img alt="Kubernets Logo" src="https://cdn.worldvectorlogo.com/logos/kubernets.svg" height="150" />
    <h2 align="center">Kubernetes API</h2>
</p>


#### Passing metadata through the Downward API

The Downward API enables you to expose the pod’s own metadata to the processes running inside that pod. Currently, it allows you to pass the following information to your containers:

- The pod’s name
- The pod’s IP address
- The namespace the pod belongs to
- The name of the node the pod is running on
- The name of the service account the pod is running under
- The CPU and memory requests for each container
- The CPU and memory limits for each container
- The pod’s labels
- The pod’s annotations

First, let’s look at how you can pass the pod’s and container’s metadata to the container through environment variables. You’ll create a simple single-container pod from the following listing’s manifest.

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: downward
spec:
  containers:
    -
      command:
        - sleep
        - "9999999"
      env:
        -
          name: POD_NAME
          valueFrom:
            fieldRef:
              fieldPath: metadata.name
        -
          name: POD_NAMESPACE
          valueFrom:
            fieldRef:
              fieldPath: metadata.namespace
        -
          name: POD_IP
          valueFrom:
            fieldRef:
              fieldPath: status.podIP
        -
          name: NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
        -
          name: SERVICE_ACCOUNT
          valueFrom:
            fieldRef:
              fieldPath: spec.serviceAccountName
        -
          name: CONTAINER_CPU_REQUEST_MILLICORES
          valueFrom:
            resourceFieldRef:
              divisor: 1m
              resource: requests.cpu
        -
          name: CONTAINER_MEMORY_LIMIT_KIBIBYTES
          valueFrom:
            resourceFieldRef:
              divisor: 1Ki
              resource: limits.memory
      image: busybox
      name: main
      resources:
        limits:
          cpu: 100m
          memory: 4Mi
        requests:
          cpu: 15m
          memory: 100Ki
```


#### Passing metadata through files in a downwardAPI volume

If you prefer to expose the metadata through files instead of environment variables, you can define a downwardAPI volume and mount it into your container. You must use a downwardAPI volume for exposing the pod’s labels or its annotations

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  annotations:
    key1: value1
    key2: |
        multi
        line
        value
  labels:
    foo: bar
  name: downward
spec:
  containers:
    -
      command:
        - sleep
        - "9999999"
      image: busybox
      name: main
      resources:
        limits:
          cpu: 100m
          memory: 4Mi
        requests:
          cpu: 15m
          memory: 100Ki
      volumeMounts:
        -
          mountPath: /etc/downward
          name: downward
  volumes:
    -
      downwardAPI:
        items:
          -
            fieldRef:
              fieldPath: metadata.name
            path: podName
          -
            fieldRef:
              fieldPath: metadata.namespace
            path: podNamespace
          -
            fieldRef:
              fieldPath: metadata.labels
            path: labels
          -
            fieldRef:
              fieldPath: metadata.annotations
            path: annotations
          -
            path: containerCpuRequestMilliCores
            resourceFieldRef:
              containerName: main
              divisor: 1m
              resource: requests.cpu
          -
            path: containerMemoryLimitBytes
            resourceFieldRef:
              containerName: main
              divisor: 1
              resource: limits.memory
      name: downward
```

```console
$ kubectl exec downward ls -lL /etc/downward

.....  4 May 25 10:23 annotations
.....  2 May 25 10:23 containerCpuRequestMilliCores
.....  7 May 25 10:23 containerMemoryLimitBytes
.....  9 May 25 10:23 labels
.....  8 May 25 10:23 podName
.....  7 May 25 10:23 podNamespace
```

```console
$ kubectl exec downward cat /etc/downward/labels

foo="bar"
```

```console
$ kubectl exec downward cat /etc/downward/annotations

key1="value1"
key2="multi\nline\nvalue\n"
kubernetes.io/config.seen="2016-11-28T14:27:45.664924282Z"
kubernetes.io/config.source="api"
```

You may remember that labels and annotations can be modified while a pod is running. As you might expect, when they change, Kubernetes updates the files holding them, allowing the pod to always see up-to-date data. This also explains why labels and annotations can’t be exposed through environment variables. Because environment variable values can’t be updated afterward, if the labels or annotations of a pod were exposed through environment variables, there’s no way to expose the new values after they’re modified.


#### Talking to the Kubernetes API server

The `kubectl proxy` command runs a proxy server that accepts HTTP connections on your local machine and proxies them to the API server while taking care of authentication, so you don’t need to pass the authentication token in every request. It also makes sure you’re talking to the actual API server and not a man in the middle.

```console
$ kubectl proxy

Starting to serve on 127.0.0.1:8001
```

```console
$ curl http://localhost:8001

{
  "paths": [
    "/api",
    "/api/v1",
    "/apis",
    "/apis/apps",
    "/apis/apps/v1beta1",
    ...
    "/apis/batch",
    "/apis/batch/v1",
    "/apis/batch/v2alpha1",
```

Now, let’s see how to talk to it from within a pod, where you (usually) don’t have `kubectl`. Therefore, to talk to the API server from inside a pod, you need to take care of three things:

- Find the location of the API server.
- Make sure you’re talking to the API server and not something impersonating it.
- Authenticate with the server; otherwise it won’t let you see or do anything.

You need to use a container image that contains the curl binary

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: curl
spec:
  containers:
  - name: main
    image: tutum/curl
    command: ["sleep", "9999999"]
```

After creating the pod, run kubectl exec to run a bash shell inside its container:

```console
$ kubectl exec -it curl bash
root@curl:/#
```

You can get both the `IP address` and the port of the API server by looking up the `KUBERNETES_SERVICE_HOST` and `KUBERNETES_SERVICE_PORT` variables (inside the container):

```console
root@curl:/# env | grep KUBERNETES_SERVICE

KUBERNETES_SERVICE_PORT=443
KUBERNETES_SERVICE_HOST=10.0.0.1
KUBERNETES_SERVICE_PORT_HTTPS=443
```

In the previous chapter, while discussing Secrets, we looked at an automatically created Secret called `default-token-xyz`, which is mounted into each container at `/var/run/secrets/kubernetes.io/serviceaccount/`. Let’s see the contents of that Secret again, by listing files in that directory:

```console
root@curl:/#ls/var/run/secrets/kubernetes.io/serviceaccount/

ca.crt    namespace    token
```

You may also remember that each service also gets a DNS entry, so you don’t even need to look up the environment variables, but instead simply point curl to `https://kubernetes`.

```console
root@curl:/# export CURL_CA_BUNDLE=/var/run/secrets/kubernetes.io/serviceaccount/ca.crt
```

```console
root@curl:/# TOKEN=$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)
```

```console
root@curl:/# curl -H "Authorization: Bearer $TOKEN" https://kubernetes

    {
    "paths": [
        "/api",
        "/api/v1",
        "/apis",
        "/apis/apps",
        "/apis/apps/v1beta1",
        "/apis/authorization.k8s.io",
        ...
    "/ui/",
    "/version" ]
    }
```

```console
root@curl:/# NS=$(cat /var/run/secrets/kubernetes.io/serviceaccount/namespace)

root@curl:/# curl -H "Authorization: Bearer $TOKEN" https://kubernetes/api/v1/namespaces/$NS/pods

{
  "kind": "PodList",
  "apiVersion": "v1",
```
