<p align="center">
    <img alt="Kubernets Logo" src="https://cdn.worldvectorlogo.com/logos/kubernets.svg" height="150" />
    <h2 align="center">Deployments</h2>
</p>


A `Deployment` is a higher-level resource meant for deploying applications and updating them declaratively, instead of doing it through a `ReplicationController` or a `ReplicaSet`, which are both considered lower-level concepts.

#### Creating a Deployment

A `Deployment` is also composed of a label selector, a desired replica count, and a pod template. In addition to that, it also contains a field, which specifies a deployment strategy that defines how an update should be performed when the `Deployment` resource is modified.

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kubia
spec:
  replicas: 3
  template:
    metadata:
      labels:
        app: kubia
      name: kubia
    spec:
      containers:
        -
          image: "luksa/kubia:v1"
          name: nodejs
```

You’re now ready to create the Deployment:

```console
$ kubectl create -f kubia-deployment-v1.yaml --record

deployment "kubia" created
```

Be sure to include the `--record` command-line option when creating it. This records the command in the revision history

Checking a Deployment’s status:

```console
$ kubectl rollout status deployment kubia

deployment kubia successfully rolled out
```

You should see the three pod replicas up and running:

```console
$ kubectl get po
```

The three pods created by the Deployment include an additional numeric value in the middle of their names. The number corresponds to the hashed value of the pod template in the `Deployment` and the `ReplicaSet` managing these pods.

a `Deployment` doesn’t manage pods directly. Instead, it creates `ReplicaSets` and leaves the managing to them, so let’s look at the ReplicaSet created by your `Deployment`:

```console
$ kubectl get replicasets

NAME               DESIRED   CURRENT   AGE
kubia-1506449474   3         3         10s
```


#### Updating a Deployment

Kubernetes will take all the steps necessary to get the actual system state to what’s defined in the resource. Similar to scaling a `ReplicationController` or `ReplicaSet` up or down, all you need to do is reference a new image tag in the Deployment’s pod template and leave it to Kubernetes to transform your system so it matches the new desired state.

The `Recreate` strategy causes all old pods to be deleted before the new ones are created. Use this strategy when your application doesn’t support running multiple versions in parallel and requires the old version to be stopped completely before the new one is started. This strategy does involve a short period of time when your app becomes completely unavailable.

The `RollingUpdate` strategy, on the other hand, removes old pods one by one, while adding new ones at the same time, keeping the application available throughout the whole process, and ensuring there’s no drop in its capacity to handle requests. This is the default strategy. The upper and lower limits for the number of pods above or below the desired replica count are configurable. You should use this strategy only when your app can handle running both the old and new version at the same time.

When you execute this command, you’re updating the kubia Deployment’s pod template so the image used in its nodejs container is changed to `luksa/kubia:v2`

```console
$ kubectl set image deployment kubia nodejs=luksa/kubia:v2

deployment "kubia" image updated
```

You can follow the progress of the rollout with kubectl rollout status:

```console
$ kubectl rollout status deployment kubia
```

To undo the last rollout of a Deployment:

```console
$ kubectl rollout undo deployment kubia

deployment "kubia" rolled back
```

This rolls the Deployment back to the previous revision.

The revision history can be displayed with the kubectl rollout history command:

```console
$ kubectl rollout history deployment kubia

deployments "kubia":
REVISION    CHANGE-CAUSE
2           kubectl set image deployment kubia nodejs=luksa/kubia:v2
3           kubectl set image deployment kubia nodejs=luksa/kubia:v3
```

You can roll back to a specific revision by specifying the revision in the undo command.

```console
$ kubectl rollout undo deployment kubia --to-revision=1
```


#### Controlling the rate of the rollout

Two properties affect how many pods are replaced at once during a Deployment’s rolling update. They are `maxSurge` and `maxUnavailable` and can be set as part of the `rollingUpdate` sub-property of the Deployment’s strategy attribute, as shown in the following listing.

```yaml
spec:
    strategy:
        rollingUpdate:
            maxSurge: 1
            maxUnavailable: 0
        type: RollingUpdate
```

`maxSurge`: Determines how many pod instances you allow to exist above the desired replica count configured on the Deployment. It defaults to 25%, so there can be at most 25% more pod instances than the desired count. If the desired replica count is set to four, there will never be more than five pod instances running at the same time during an update. When converting a percentage to an absolute number, the number is rounded up. Instead of a percentage, the value can also be an absolute value (for example, one or two additional pods can be allowed).

`maxUnavailable`: Determines how many pod instances can be unavailable relative to the desired replica count during the update. It also defaults to 25%, so the number of avail- able pod instances must never fall below 75% of the desired replica count. Here, when converting a percentage to an absolute number, the number is rounded down. If the desired replica count is set to four and the percentage is 25%, only one pod can be unavailable. There will always be at least three pod instances available to serve requests during the whole rollout. As with maxSurge, you can also specify an absolute value instead of a percentage.


You can trigger the rollout by changing the image to `luksa/kubia:v4`, but then immediately (within a few seconds) pause the rollout (`canary release`).

A `canary release` is a technique for minimizing the risk of rolling out a bad version of an application and it affecting all your users. Instead of rolling out the new version to everyone, you replace only one or a small number of old pods with new ones.

```console
$ kubectl set image deployment kubia nodejs=luksa/kubia:v4
deployment "kubia" image updated

$ kubectl rollout pause deployment kubia
deployment "kubia" paused
```

Once you’re confident the new version works as it should, you can resume the deployment to replace all the old pods with new ones:

```console
$ kubectl rollout resume deployment kubia
deployment "kubia" resumed
```

#### Blocking rollouts of bad versions

The `minReadySeconds` property specifies how long a newly created pod should be ready before the pod is treated as available. Until the pod is available, the rollout process will not continue (remember the `maxUnavailable` property?). A pod is ready when readiness probes of all its containers return a success. If a new pod isn’t functioning properly and its readiness probe starts failing before `minReadySeconds` have passed, the rollout of the new version will effectively be blocked.

The fact that the deployment will stuck is a good thing, because if it had continued replacing the old pods with the new ones, you’d end up with a completely non-working service

```yaml
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: kubia
spec:
  minReadySeconds: 10
  replicas: 3
  strategy:
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
    type: RollingUpdate
  template:
    metadata:
      labels:
        app: kubia
      name: kubia
    spec:
      containers:
        -
          image: "luksa/kubia:v3"
          name: nodejs
          readinessProbe:
            httpGet:
              path: /
              port: 8080
            periodSeconds: 1
```

Usually, you’d set `minReadySeconds` to something much higher to make sure pods keep reporting they’re ready after they’ve already started receiving actual traffic.