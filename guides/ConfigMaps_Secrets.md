<p align="center">
    <img alt="Kubernets Logo" src="https://cdn.worldvectorlogo.com/logos/kubernets.svg" height="150" />
    <h2 align="center">ConfigMaps and Secrets</h2>
</p>


Specifying environment variables in a container definition:

```yaml
---
kind: Pod
spec:
  containers:
    -
      env:
        -
          name: INTERVAL
          value: "30"
      image: "luksa/fortune:env"
      name: html-generator
```

Referring to other environment variables in a variable’s value:

```yaml
env:
  -
    name: FIRST_VAR
    value: foo
  -
    name: SECOND_VAR
    value: $(FIRST_VAR)bar
```


#### Creating a `ConfigMap`:

To start with the simplest example, you’ll first create a map with a single key and use it to fill the `INTERVAL` environment variable from your previous example. You’ll create the `ConfigMap` with the special `kubectl create configmap` command instead of posting a YAML with the generic `kubectl create -f` command.

```
$ kubectl create configmap fortune-config --from-literal=sleep-interval=25

configmap "fortune-config" created
```

This creates a `ConfigMap` called `fortune-config` with the single-entry `sleep-interval=25`

`ConfigMaps` usually contain more than one entry. To create a `ConfigMap` with multiple literal entries, you add multiple `--from-literal`arguments:

```
$ kubectl create configmap myconfigmap --from-literal=foo=bar --from-literal=bar=baz --from-literal=one=two
```

Let’s inspect the YAML descriptor of the `ConfigMap` you created by using the `kubectl get` command

```
kubectl get configmap fortune-config -o yaml
```

You could easily have written this YAML yourself and posted it to the Kubernetes API with the well-known

```yaml
# fortune-config.yaml
---
apiVersion: v1
data:
  sleep-interval: "25"
kind: ConfigMap
metadata:
  name: fortune-config
```

```
$ kubectl create -f fortune-config.yaml
```

`ConfigMaps` can also store coarse-grained config data, such as complete config files.

To do this, the `kubectl create configmap` command also supports reading files from disk and storing them as individual entries in the `ConfigMap`.

```
$ kubectl create configmap my-config --from-file=config-file.conf
```

This command will store the file’s contents under the key `customkey`. As with literals, you can add multiple files by using the `--from-file` argument multiple times.

```
$ kubectl create configmap my-config --from-file=customkey=config-file.conf
```

Instead of importing each file individually, you can even import all files from a file directory:

```
$ kubectl create configmap my-config --from-file=/path/to/dir
```

In this case, kubectl will create an individual map entry for each file in the specified directory, but only for files whose name is a valid `ConfigMap` key.

When creating ConfigMaps, you can use a combination of all the options mentioned here

```
$ kubectl create configmap my-config \
    --from-file=foo.json \
    --from-file=bar=foobar.conf \
    --from-file=config-opts/ \ # A whole directory
    --from-literal=some=thing
```

Passing a `ConfigMap` entry to a container as an environment variable

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: fortune-env-from-configmap
spec:
  containers:
    -
      env:
        -
          name: INTERVAL
          valueFrom:
            configMapKeyRef:
              key: sleep-interval
              name: fortune-config
      image: "luksa/fortune:env"
```

You defined an environment variable called `INTERVAL` and set its value to whatever is stored in the `fortune-config` `ConfigMap` under the key `sleep-interval`. When the process running in container reads the `INTERVAL` environment variable, it will see the value 25.

Passing all entries of a `ConfigMap` as environment variables at once

Kubernetes version 1.6 provides a way to expose all entries of a `ConfigMap` as environment variables.

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: fortune-env-from-configmap
spec:
  containers:
    -
      envFrom:
        -
          configMapRef:
            name: my-config-map
          prefix: CONFIG_
      image: some-image
```

As you can see, you can also specify a prefix for the environment variables (`CONFIG_` in this case). This results in the following two environment variables being present inside the container: `CONFIG_FOO`, `CONFIG_BAR` ... etc

Passing a `ConfigMap` entry as a command-line argument

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: fortune-args-from-configmap
spec:
  containers:
    -
      args:
        - $(INTERVAL)
      env:
        -
          name: INTERVAL
          valueFrom:
            configMapKeyRef:
              key: sleep-interval
              name: fortune-config
      image: "luksa/fortune:args"
```

You defined the environment variable exactly as you did before, but then you used the `$(ENV_VARIABLE_NAME)` syntax to have Kubernetes inject the value of the variable into the argument.

Using a `configMap` volume to expose `ConfigMap` entries as files

A `configMap` volume will expose each entry of the `ConfigMap` as a file. The process running in the container can obtain the entry’s value by reading the contents of the file. Although this method is mostly meant for passing large config files to the container, nothing prevents you from passing short single values this way.

```conf
# configmap-files/my-nginx-config.conf
server {
  listen            80;
  server_name       www.kubia-example.com;

  gzip on;
  gzip_types text/plain application/xml;
  location / {
    root   /usr/share/nginx/html;
    index  index.html index.htm;
   }
}
```

```bash
$ echo "25" > configmap-files/sleep-interval
```

Now create a ConfigMap from all the files in the directory like this:

```
$ kubectl create configmap fortune-config --from-file=configmap-files

configmap "fortune-config" created
```

The following listing shows what the YAML of this ConfigMap looks like.

```
$ kubectl get configmap fortune-config -o yaml

apiVersion: v1
data:
  my-nginx-config.conf: |
    server {
      listen              80;
      server_name         www.kubia-example.com;
      gzip on;
      gzip_types text/plain application/xml;
      location / {
        root   /usr/share/nginx/html;
        index  index.html index.htm;
      }
    }
  sleep-interval: |
    25
kind: ConfigMap
```

The `ConfigMap` contains two entries, with keys corresponding to the actual names of the files they were created from. You’ll now use the `ConfigMap` in both of your pod’s containers.

Passing `ConfigMap` entries to a pod as files in a volume

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: fortune-configmap-volume
spec:
  containers:
    -
      image: "nginx:alpine"
      name: web-server
      volumeMounts:
        -
          mountPath: /etc/nginx/conf.d
          name: config
          readOnly: true
  volumes:
    -
      configMap:
        name: fortune-config
      name: config
```

The web server should now be configured to compress the responses it sends.

```yaml
---
volumes:
  -
    configMap:
      items: # Selecting which entries to include in the volume by listing them
        -
          key: my-nginx-config.conf # You want the entry under this key included.
          path: gzip.conf # The entry’s value should be stored in this file.
      name: fortune-config
    name: config
```

When specifying individual entries, you need to set the filename for each individual entry, along with the entry’s key. If you run the pod from the previous listing, the `/etc/nginx/conf.d` directory is kept nice and clean, because it only contains the `gzip.conf` file and nothing else.

In both this and in your previous example, you mounted the volume as a directory, which means you’ve hidden any files that are stored in the `/etc/nginx/conf.d` directory in the container image itself.

You’re now wondering how to add individual files from a `ConfigMap` into an existing directory without hiding existing files stored in it. An additional `subPath` property on the `volumeMount` allows you to mount either a single file or a single directory from the volume instead of mounting the whole volume.

```yaml
---
spec:
  containers:
    -
      image: some/image
      volumeMounts:
        -
          mountPath: /etc/someconfig.conf
          name: myvolume
          subPath: myconfig.conf
```

#### Introducing Secrets:

Now, you’ll create your own little Secret.

```
$ echo bar > foo
```

```
$ kubectl create secret generic fortune-https --from-file=foo

secret "fortune-https" created
```

```
$ kubectl get secret fortune-https -o yaml

apiVersion: v1
data:
  foo: YmFyCg==
kind: Secret
```

When you expose the Secret to a container through a secret volume, the value of the Secret entry is decoded and written to the file in its actual form (regardless if it’s plain text or binary). The same is also true when exposing the Secret entry through an environment variable. In both cases, the app doesn’t need to decode it, but can read the file’s contents or look up the environment variable value and use it directly.

Exposing a Secret’s entry as an environment variable

```yaml
---
env:
  -
    name: FOO_SECRET
    valueFrom:
      secretKeyRef:
        key: foo
        name: fortune-https
```


Mounting The `fortune-https` Secret In a Pod

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: fortune-https
spec:
  containers:
    -
      image: "app:image"
      name: app
      ports:
        -
          containerPort: 80
      volumeMounts:
        -
          mountPath: /etc/app
          name: config
          readOnly: true
volumes:
  -
    name: config
    secret:
      secretName: fortune-https
```

The secret volume uses an in-memory filesystem (`tmpfs`) for the Secret files. Because tmpfs is used, the sensitive data stored in the Secret is never written to disk, where it could be compromised.

Creating a Secret holding the credentials for authenticating with a Docker registry

```
$ kubectl create secret docker-registry mydockerhubsecret \
  --docker-username=myusername --docker-password=mypassword \
  --docker-email=my.email@provider.com
```

```yaml
---
apiVersion: v1
kind: Pod
metadata:
  name: private-pod
spec:
  containers:
    -
      image: "username/private:tag"
      name: main
  imagePullSecrets:
    -
      name: mydockerhubsecret
```