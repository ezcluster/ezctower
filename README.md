
# EZCLUSTER TOWER

## Usage

Execute a command and exit

```
docker run ghcr.io/ezcluster/tower:0.1.0 /bin/bash -c "ls /usr"
```




Interactive launch

```
docker run -it ghcr.io/ezcluster/tower:0.1.0 /bin/bash -l
```


If you need to perform some admin task (i.e. install package)

```
docker run -it --user 0 ghcr.io/ezcluster/tower:0.1.0 /bin/bash
```



