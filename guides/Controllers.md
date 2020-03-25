<p align="center">
    <img alt="Kubernets Logo" src="https://cdn.worldvectorlogo.com/logos/kubernets.svg" height="150" />
    <h2 align="center">Controllers</h2>
</p>


### ReplicationController

A `ReplicationController`’s job is to make sure that an exact number of pods always matches its label selector. If it doesn’t, the `ReplicationController` takes the appropriate action to reconcile the actual with the desired number.

```yaml
---
apiVersion: v1
kind: ReplicationController
metadata:
  name: koala
spec:
  replicas: 3
  selector:
    app: koala
template:
  metadata:
    labels:
      app: koala
  spec:
    containers:
      -
        image: clivern/koala
        name: koala
        ports:
          -
            containerPort: 8080
```

No need to specify a pod selector when defining a `ReplicationController`. Let kubernetes extract it from the pod template. This will keep your YAML shorter and simpler.

```yaml
---
apiVersion: v1
kind: ReplicationController
metadata:
  name: koala
spec:
  replicas: 3
template:
  spec:
    containers:
      -
        image: clivern/koala
        name: koala
        ports:
          -
            containerPort: 8080
```

```
$ kubectl create -f koala-rc.yaml
```

Now, let’s see what information the `kubectl get` command shows for replication controllers

```
$ kubectl get rc
```

```
$ kubectl describe rc koala
```

if you overrite one of the pods label, replication controller will spin another pod to reach the desired state.

```
# changing the labels of a managed pod
$ kubectl label pod ${podName} app=foo --overwrite
$ kubectl get pods --show-labels
```

Editing the `ReplicationController` definition or scaling up or scaling down

```
$ kubectl edit rc koala
```

Delete replication controller without the managed pods

```
$ kubectl delete rc koala --cascade=false
```

Delete replication controller with the managed pods

```
$ kubectl delete rc koala
```

Initially `ReplicationControllers` were the only Kubernetes component for replicating pods and rescheduling them when nodes failed. Later, a similar resource called a `ReplicaSet` was introduced. It’s a new generation of `ReplicationController` and replaces it completely (`ReplicationControllers` will eventually be deprecated).


### ReplicaSet

We’ll rewrite the `ReplicationController` into a `ReplicaSet`

```yaml
---
apiVersion: apps/v1beta2
kind: ReplicaSet
metadata:
  name: koala
spec:
  replicas: 3
  selector:
    matchLabels:
      app: koala
template:
  metadata:
    labels:
      app: koala
  spec:
    containers:
      -
        image: clivern/koala
        name: koala
        ports:
          -
            containerPort: 8080
```

The main improvements of `ReplicaSets` over `ReplicationControllers` are their more expressive label selectors. You intentionally used the simpler `matchLabels` selector in the first `ReplicaSet` example to see that `ReplicaSets` are no different from `ReplicationControllers`.

Now, you’ll rewrite the selector to use the more powerful `matchExpressions` property, as shown in the following listing.

```yaml
selector:
  matchExpressions:
    -
      key: app
      operator: In
      values:
        - koala
```

You can add additional expressions to the selector. As in the example, each expression must contain a key, an operator, and possibly (depending on the operator) a list of values. You’ll see four valid operators:

- `In` Label’s value must match one of the specified values.
- `NotIn` Label’s value must not match any of the specified values.
- `Exists` Pod must include a label with the specified key (the value isn’t important). When using this operator, you shouldn’t specify the values field.
- `DoesNotExist` Pod must not include a label with the specified key. The values property must not be specified.

If you specify multiple expressions, all those expressions must evaluate to true for the selector to match a pod. If you specify both matchLabels and matchExpressions, all the labels must match and all the expressions must evaluate to true for the pod to match the selector.


You can examine the `ReplicaSet` with `kubectl get` and `kubectl describe`

```
$ kubectl get rs
```

You can delete the `ReplicaSet` the same way you’d delete a `ReplicationController`

```
$ kubectl delete rs koala
```

Deleting the `ReplicaSet` should delete all the pods. List the pods to confirm that’s the case.


### DaemonSet

A `DaemonSet` ensures that all (or some) Nodes run a copy of a Pod. As nodes are added to the cluster, Pods are added to them. As nodes are removed from the cluster, those Pods are garbage collected. Deleting a `DaemonSet` will clean up the Pods it created.

Some typical uses of a `DaemonSet` are:

- Running a cluster storage daemon, such as `glusterd`, `ceph`, on each node.
- Running a logs collection daemon on every node, such as `fluentd` or `logstash`.
- Running a node monitoring daemon on every node, such as `Prometheus` Node Exporter, `Sysdig` Agent.

You can describe a `DaemonSet` in a YAML file like the following:

```yaml
---
apiVersion: apps/v1beta2
kind: DaemonSet
metadata:
  name: ssd-monitor
spec:
  selector:
    matchLabels:
      app: ssd-monitor
  template:
    metadata:
      labels:
        app: ssd-monitor
    spec:
      containers:
        -
          image: luksa/ssd-monitor
          name: main
      nodeSelector:
        disk: ssd
```

Let's create the `DaemonSet` from the YAML file

```
$ kubectl create -f ssd-monitor-daemonset.yaml
```

Let’s see the created `DaemonSet`

```
$ kubectl get ds
```

Those zeroes look strange. Didn’t the `DaemonSet` deploy any pods? List the pods:

```
$ kubectl get po

No resources found.
```

```
$ kubectl get node

NAME       STATUS    AGE       VERSION
minikube   Ready     4d        v1.6.0
```

Now, add the `disk=ssd` label to one of your nodes like this:

```
$ kubectl label node minikube disk=ssd
```

The `DaemonSet` should have created one pod now. Let’s see:

```
$ kubectl get po

NAME                READY     STATUS    RESTARTS   AGE
ssd-monitor-hgxwq   1/1       Running   0          35s
```

What happens if you change the node’s label?

```
$ kubectl label node minikube disk=hdd --overwrite

node "minikube" labeled
```

Let’s see if the change has any effect on the pod that was running on that node:

```
$ kubectl get po

NAME                READY     STATUS        RESTARTS   AGE
ssd-monitor-hgxwq   1/1       Terminating   0          4m
```

The pod is being terminated.


### The Job Resource

Kubernetes allows you to run a pod whose container isn’t restarted when the process running inside finishes successfully. Once it does, the pod is considered complete.

```yaml
---
apiVersion: batch/v1
kind: Job
metadata:
  name: batch-job
spec:
  template:
    metadata:
      labels:
        app: batch-job
    spec:
      containers:
        -
          image: luksa/batch-job
          name: main
      restartPolicy: OnFailure
```

In a pod’s specification, you can specify what Kubernetes should do when the processes running in the container finish. This is done through the `restartPolicy`, which defaults to Always. Job pods can’t use the default policy, because they’re not meant to run indefinitely. Therefore, you need to explicitly set the restart policy to either `OnFailure` or `Never`.

After you create this Job with the `kubectl create` command, you should see it start up a pod immediately:

```
$ kubectl get jobs

NAME        DESIRED   SUCCESSFUL   AGE
batch-job   1         0            2s
```

```
$ kubectl get po
NAME              READY     STATUS    RESTARTS   AGE
batch-job-28qf4   1/1       Running   0          4s
```

```
$ kubectl get po -a

NAME              READY     STATUS      RESTARTS   AGE
batch-job-28qf4   0/1       Completed   0          2m
```

The reason the pod isn’t deleted when it completes is to allow you to examine its logs; for example:

```
$ kubectl logs batch-job-28qf4
Fri Apr 29 09:58:22 UTC 2016 Batch job starting
Fri Apr 29 10:00:22 UTC 2016 Finished succesfully
```

If you need a Job to run more than once, you set `completions` to how many times you want the Job’s pod to run. The following listing shows an example.

```yaml
---
apiVersion: batch/v1
kind: Job
metadata:
  name: multi-completion-batch-job
spec:
  completions: 5
  template: ~
```

Instead of running single Job pods one after the other, you can also make the Job run multiple pods in parallel. You specify how many pods are allowed to run in parallel with the `parallelism` Job spec property

```yaml
---
apiVersion: batch/v1
kind: Job
metadata:
  name: multi-completion-batch-job
spec:
  completions: 5
  parallelism: 2
  template: ~
```

By setting `parallelism` to 2, the Job creates two pods and runs them in parallel:

```
$ kubectl get po

NAME                               READY   STATUS     RESTARTS   AGE
multi-completion-batch-job-lmmnk   1/1     Running    0          21s
multi-completion-batch-job-qx4nq   1/1     Running    0          21s
```


### The CronJob Resource

Job resources will be created from the CronJob resource at approximately the scheduled time. The Job then creates the pods.

```yaml
---
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: batch-job-every-fifteen-minutes
spec:
  schedule: "0,15,30,45 * * * *"
  jobTemplate:
    metadata:
      labels:
        app: periodic-batch-job
    spec:
      containers:
        -
          image: luksa/batch-job
          name: main
      restartPolicy: OnFailure
```

It may happen that the Job or pod is created and run relatively late. In that case, you can specify a deadline by specifying the `startingDeadlineSeconds` field

```yaml
---
apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: batch-job-every-fifteen-minutes
spec:
  schedule: "0,15,30,45 * * * *"
  startingDeadlineSeconds: 15
  jobTemplate:
    metadata:
      labels:
        app: periodic-batch-job
    spec:
      containers:
        -
          image: luksa/batch-job
          name: main
      restartPolicy: OnFailure
```