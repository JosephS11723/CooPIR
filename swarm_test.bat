::create the swarm
docker swarm init
::create the registry
docker service create --name registry --publish published=5000,target=5000 registry:2
docker service ls

::push the image to the registry
docker-compose push

::deploy the stack to the swarm (named stackdemo)
docker stack deploy --compose-file docker-compose.yml CooPIR

::check that it is running
docker stack services CooPIR
PAUSE