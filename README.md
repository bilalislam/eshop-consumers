# eshop-consumers

#### Clean Architecture In Golang
----------------------------
Note : When docker build and you are getting error about not build the docker context
then you should build from one upper step . So that to changes docker context.

False
------
into docker file where exists here :

docker build -t basket-deleted-consumer .

True
------
docker build -f tools/docker/basketdeleted/Dockerfile -t basket-deleted-consumer .


get token for private repo in dockerfile:
https://stackoverflow.com/questions/47617432/docker-change-gitconfig-with-token-for-private-repo-access