
# EZCLUSTER TOWER

## Usage

Execute a command and exit

```
docker run ghcr.io/ezcluster/tower:0.1.0 /bin/bash -c "ls /usr"
```

```
rm -rf /tmp/git && mkdir -p /tmp/git && chmod 777 /tmp/git && \
docker run -v /tmp/git:/home/tower/git ghcr.io/ezcluster/tower:0.1.0 /bin/bash -c "cd ~/git && git clone https://github.com/ezcluster/ezcluster.git"
```

Interactive launch

```
docker run -it ghcr.io/ezcluster/tower:0.1.0 /bin/bash -l
```

Launch with web access

```
rm -rf /tmp/git && mkdir -p /tmp/git && chmod 777 /tmp/git && \
docker run -v /tmp/git:/home/tower/git -p 7681:7681 ghcr.io/ezcluster/tower:0.1.0 ttyd  -d 4 -p 7681 /bin/bash -l
```

[GO](http://localhost:7681)


Launch in a VM (builder2.ops.scw01) targeting scw01/project33/kspray1

```
rm -rf /tmp/git && mkdir -p /tmp/git && chmod 777 /tmp/git && \
    docker run -v /tmp/git:/home/tower/git -p 7681:7681 \
    --env EZCT_WORKDIR=/home/tower/git \
    --env EZCT_REPO=https://github.com/KubeDP/osiac.git \
    --env EZCT_BRANCH=test2 \
    --env EZCT_PATH=iac/scw01/33-project33/kspray1 \
    --env EZCT_GIT_USERNAME=Xxxxxxxxx \
    --env EZCT_GIT_TOKEN=github_pat_xxxxxxxxxxxxxxxxxxxxxxxxxxxxx \
    ghcr.io/ezcluster/tower:0.1.0 ttyd  -d 4 -p 7681 /bin/bash -l
```



[GO](http://builder2.ops.scw01:7681)
