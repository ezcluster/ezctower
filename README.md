
# EZCLUSTER TOWER

## Usage

Execute a command and exit

```
docker run ghcr.io/ezcluster/tower:0.1.0 /bin/bash -c "ls /usr"
```

```
mkdir /tmp/git
docker run -v /tmp/git:/home/tower/git ghcr.io/ezcluster/tower:0.1.0 /bin/bash -c "cd ~/git && git clone https://github.com/ezcluster/ezcluster.git"
```

Interactive launch

```
docker run -it ghcr.io/ezcluster/tower:0.1.0 /bin/bash -l
```

Launch with web access

```
mkdir /tmp/git
docker run -v /tmp/git:/home/tower/git -p 7681:7681 ghcr.io/ezcluster/tower:0.1.0 ttyd  -d 4 -p 7681 /bin/bash -l
```

[GO](http://localhost:7681)

