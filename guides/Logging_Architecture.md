<p align="center">
    <img alt="Kubernets Logo" src="https://cdn.worldvectorlogo.com/logos/kubernets.svg" height="150" />
    <h2 align="center">Logging Architecture</h2>
</p>

#### Fluentd 

It is an open source data collector for unified logging layer. Fluentd allows you to unify data collection and consumption for a better use and understanding of data.

Basic Setup on `Ubuntu 18.04`

```zsh
$ curl -L https://toolbelt.treasuredata.com/sh/install-ubuntu-bionic-td-agent2.5.sh | sh

$ sudo /etc/init.d/td-agent status

# Test input
$ curl -X POST -d 'json={"json":"message"}' http://x.x.x.x:8888/debug.test

# Output is log by default
$ tail -f /var/log/td-agent/td-agent.log

# systemd service file 
$ cat /lib/systemd/system/td-agent.service 

# default configs
$ cat /etc/td-agent/td-agent.conf 
```
