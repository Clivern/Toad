<p align="center">
    <img alt="Kubernets Logo" src="https://cdn.worldvectorlogo.com/logos/kubernets.svg" height="150" />
    <h2 align="center">Volumes</h2>
</p>

A wide variety of volume types is available. Several are generic, while others are specific to the actual storage technologies used underneath

- `emptyDir`: A simple empty directory used for storing transient data.
- `hostPath`: Used for mounting directories from the worker node’s filesystem into the pod.
- `gitRepo`: A volume initialized by checking out the contents of a Git repository.
- `nfs`: An NFS share mounted into the pod.
- `gcePersistentDisk` (Google Compute Engine Persistent Disk), `awsElastic-BlockStore` (Amazon Web Services Elastic Block Store Volume), `azureDisk` (Microsoft Azure Disk Volume): Used for mounting cloud provider-specific storage.
- `cinder`, `cephfs`, `iscsi`, `flocker`, `glusterfs`, `quobyte`, `rbd`, `flexVolume`, `vsphere-Volume`, `photonPersistentDisk`, `scaleIO`: Used for mounting other types of network storage.
- `configMap`, `secret`, `downwardAPI`: Special types of volumes used to expose certain Kubernetes resources and cluster information to the pod.
- `persistentVolumeClaim`: A way to use a pre- or dynamically provisioned persistent storage.


#### `emptyDir` Volume:

The volume starts out as an empty directory. The app running inside the pod can then write any files it needs to it. Because the volume’s lifetime is tied to that of the pod, the volume’s contents are lost when the pod is deleted.

An `emptyDir` volume is especially useful for sharing files between containers running in the same pod. But it can also be used by a single container for when a con- tainer needs to write data to disk temporarily, such as when performing a sort operation on a large dataset, which can’t fit into the available memory.

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: fortune
spec:
  containers:
    -
      image: luksa/fortune
      name: html-generator
      volumeMounts:
        -
          mountPath: /var/htdocs # The volume called html is mounted at /var/htdocs in the container.
          name: html
    -
      image: "nginx:alpine"
      name: web-server
      ports:
        -
          containerPort: 80
          protocol: TCP
      volumeMounts:
        -
          mountPath: /usr/share/nginx/html # The volume called html is mounted at /usr/share/nginx/html in the container as read-only.
          name: html
          readOnly: true
  volumes: # A single emptyDir volume called html that’s mounted in the two containers above
    -
      emptyDir: {}
      name: html
```

The `emptyDir` you used as the volume was created on the actual disk of the worker node hosting your pod, so its performance depends on the type of the node’s disks. But you can tell Kubernetes to create the `emptyDir` on a `tmpfs` filesystem (in memory instead of on disk). To do this, set the `emptyDir`’s `medium` to `Memory` like this:

```yaml
volumes:
  -
    emptyDir:
      medium: Memory
    name: html
```

An `emptyDir` volume is the simplest type of volume, but other types build upon it.


#### `gitRepo` Volume:

A `gitRepo` volume is basically an `emptyDir` volume that gets populated by cloning a Git repository and checking out a specific revision when the pod is starting up (but before its containers are created).

After the `gitRepo` volume is created, it isn’t kept in sync with the repo it’s referencing. The files in the volume will not be updated when you push additional commits to the Git repository.

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: gitrepo-volume-pod
spec:
  containers:
    -
      image: "nginx:alpine"
      name: web-server
      ports:
        -
          containerPort: 80
          protocol: TCP
      volumeMounts:
        -
          mountPath: /usr/share/nginx/html
          name: html
          readOnly: true
  volumes:
    -
      gitRepo:
        directory: "." # You want the repo to be cloned into the root dir of the volume.
        repository: "https://github.com/luksa/kubia-website-example.git"
        revision: master
      name: html
```

the `gitRepo` volume simple and not add any support for cloning private repositories through the SSH protocol, because that would require adding additional config options to the `gitRepo` volume.

If you want to clone a private Git repo into your container, you should use a git-sync sidecar container or a similar method instead of a `gitRepo` volume.


#### `hostPath` Volume:

`hostPath` volumes are the first type of persistent storage we’re introducing, because both the `gitRepo` and `emptyDir` volumes contents get deleted when a pod is torn down, whereas a `hostPath` volume’s contents don’t. If a pod is deleted and the next pod uses a `hostPath` volume pointing to the same path on the host, the new pod will see whatever was left behind by the previous pod, but only if it’s scheduled to the same node as the first pod.


#### `nfs` Volume:

If your cluster is running on your own set of servers, you have a vast array of other supported options for mounting external storage inside your volume. For example, to mount a simple `NFS` share, you only need to specify the `NFS` server and the path exported by the server, as shown in the following listing.

```yaml
volumes:
  -
    name: mongodb-data
    nfs:
      path: /some/path
      server: "1.2.3.4"
```


#### `PersistentVolumes` and `PersistentVolumeClaims`:

To enable apps to request storage in a Kubernetes cluster without having to deal with infrastructure specifics, two new resources were introduced. They are `PersistentVolumes` and `PersistentVolumeClaims`.

- Cluster admin sets up some type of network storage (NFS export or similar.
- Admin then creates a `PersistentVolume` (`PV`) by posting a `PV` descriptor to the Kubernetes API.
- User creates a `PersistentVolumeClaim` (`PVC`).
- Kubernetes finds a `PV` of adequate size and access mode and binds the `PVC` to the `PV`.
- User creates a pod with a volume referencing the `PVC`.

First let's create a `PersistentVolume` backed by the GCE Persistent Disk

```yaml
---
apiVersion: v1
kind: PersistentVolume
metadata:
  name: mongodb-pv
spec:
  accessModes: # It can either be mounted by a single client for reading and writing or by multiple clients for reading only.
    - ReadWriteOnce
    - ReadOnlyMany
  capacity:
    storage: 1Gi # # Defining the PersistentVolume’s size
  gcePersistentDisk:
    fsType: ext4
    pdName: mongodb
  persistentVolumeReclaimPolicy: Retain # After the claim is released, the PersistentVolume should be retained (not erased or deleted).
```

After you create the `PersistentVolume` with the kubectl create command, it should be ready to be claimed. See if it is by listing all `PersistentVolumes`:

```
$ kubectl get pv

NAME         CAPACITY   RECLAIMPOLICY   ACCESSMODES   STATUS      CLAIM
mongodb-pv   1Gi        Retain          RWO,ROX       Available
```

As expected, the `PersistentVolume` is shown as Available, because you haven’t yet created the `PersistentVolumeClaim`.

*Note: `PersistentVolumes` don’t belong to any namespace. They’re cluster-level resources like nodes but `PersistentVolumeClaim` belong to namespaces*

Say you need to deploy a pod that requires persistent storage. You’ll use the `PersistentVolume` you created earlier. But you can’t use it directly in the pod. You need to claim it first.

```yaml
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mongodb-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
  storageClassName: ""
```

As soon as you create the claim, Kubernetes finds the appropriate `PersistentVolume` and binds it to the claim. The `PersistentVolume`’s capacity must be large enough to accommodate what the claim requests. Additionally, the volume’s access modes must include the access modes requested by the claim. In your case, the claim requests 1 GiB of storage and a `ReadWriteOnce` access mode. The `PersistentVolume` you created earlier matches those two requirements so it is bound to your claim. You can see this by inspecting the claim.

List all `PersistentVolumeClaims` to see the state of your `PVC`:

```
$ kubectl get pvc

NAME          STATUS    VOLUME       CAPACITY   ACCESSMODES   AGE
mongodb-pvc   Bound     mongodb-pv   1Gi        RWO,ROX       3s
```

The claim is shown as `Bound` to `PersistentVolume` `mongodb-pv`. Note the abbreviations used for the access modes:
- `RWO` — `ReadWriteOnce`: Only a single node can mount the volume for reading and writing.
- `ROX` — `ReadOnlyMany`: Multiple nodes can mount the volume for reading.
- `RWX` — `ReadWriteMany`: Multiple nodes can mount the volume for both reading and writing.

*Note: `RWO`, `ROX`, and `RWX` pertain to the number of worker nodes that can use the volume at the same time, not to the number of pods!*

You can also see that the `PersistentVolume` is now `Bound` and no longer Available by inspecting it with `kubectl get`:

```
$ kubectl get pv

NAME         CAPACITY   ACCESSMODES   STATUS   CLAIM                 AGE
mongodb-pv   1Gi        RWO,ROX       Bound    default/mongodb-pvc   1m
```

The `PersistentVolume` shows it’s bound to claim `default/mongodb-pvc`. The `default` part is the namespace the claim resides in.

To use it inside a pod, you need to reference the `PersistentVolumeClaim` by name inside the pod’s volume

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: mongodb
spec:
  containers:
    -
      image: mongo
      name: mongodb
      ports:
        -
          containerPort: 27017
          protocol: TCP
      volumeMounts:
        -
          mountPath: /data/db
          name: mongodb-data
  volumes:
    -
      name: mongodb-data
      persistentVolumeClaim:
        claimName: mongodb-pvc # Referencing the PersistentVolumeClaim by name in the pod volume
```

What if you delete the pod and the claim

```
$ kubectl delete pod mongodb
pod "mongodb" deleted

$ kubectl delete pvc mongodb-pvc
persistentvolumeclaim "mongodb-pvc" deleted
```

What if you create the PersistentVolumeClaim again? Will it be bound to the `PersistentVolume` or not? After you create the claim, what does kubectl get pvc show?

```
$ kubectl get pvc

NAME           STATUS    VOLUME       CAPACITY   ACCESSMODES   AGE
mongodb-pvc    Pending                                         13s
```

The claim’s status is shown as Pending

```
$ kubectl get pv

NAME        CAPACITY  ACCESSMODES  STATUS    CLAIM               REASON AGE
mongodb-pv  1Gi       RWO,ROX      Released  default/mongodb-pvc        5m
```

Because you’ve already used the volume, it may contain data and shouldn’t be bound to a completely new claim without giving the cluster admin a chance to clean it up. Without this, a new pod using the same `PersistentVolume` could read the data stored there by the previous pod.

You told Kubernetes you wanted your `PersistentVolume` to behave like this when you created it by setting its `persistentVolumeReclaimPolicy` to `Retain`.

Two other possible reclaim policies exist: `Recycle` and `Delete`. The first one deletes the volume’s contents and makes the volume available to be claimed again. This way, the `PersistentVolume` can be reused multiple times by different `PersistentVolumeClaims` and different pods.

The `Delete` policy, on the other hand, deletes the underlying storage.


#### Dynamic provisioning of `PersistentVolumes`:

`PersistentVolumes` and `PersistentVolumeClaims` makes it easy to obtain persistent storage without the developer having to deal with the actual storage technology used underneath. But this still requires a cluster administrator to provision the actual storage up front. Luckily, Kubernetes can also perform this job automatically through dynamic provisioning of `PersistentVolumes`.

The cluster admin, instead of creating `PersistentVolumes`, can deploy a `PersistentVolume` provisioner and define one or more `StorageClass` objects to let users choose what type of `PersistentVolume` they want. The users can refer to the `StorageClass` in their `PersistentVolumeClaims` and the provisioner will take that into account when provisioning the persistent storage.

*Note: Similar to `PersistentVolumes`, `StorageClass` resources aren’t namespaced.*

Before a user can create a `PersistentVolumeClaim`, which will result in a new `PersistentVolume` being provisioned, an admin needs to create one or more `StorageClass` resources.

```yaml
---
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: fast
parameters:
  type: pd-ssd
  zone: europe-west1-b
provisioner: kubernetes.io/gce-pd
```

The `StorageClass` resource specifies which provisioner should be used for provisioning the `PersistentVolume` when a `PersistentVolumeClaim` requests this `StorageClass`. The parameters defined in the `StorageClass` definition are passed to the provisioner and are specific to each provisioner plugin.

You can modify your mongodb-pvc to use dynamic provisioning.

```yaml
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mongodb-pvc
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
  storageClassName: fast
```

The cluster admin can create multiple storage classes with different performance or other characteristics. The developer then decides which one is most appropriate for each claim they create.

When you created your custom storage class called fast, you didn’t check if any existing storage classes were already defined in your cluster.

```
$ kubectl get sc

NAME                TYPE
fast                kubernetes.io/gce-pd
standard (default)  kubernetes.io/gce-pd
```

We’re using `sc` as shorthand for `storageclass`.

*Note: Specifying an empty string as the storage class name (`storageClassName: ""`) ensures the PVC binds to a pre-provisioned PV instead of dynamically provisioning a new one.*