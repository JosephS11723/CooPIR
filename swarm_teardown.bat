::bring down the stack
docker stack rm CooPIR

::bring down the registry
docker service rm registry

::close the swarm
docker swarm leave --force
PAUSE