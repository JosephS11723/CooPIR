docker volume prune -f
docker rmi $(docker images --filter "dangling=true" -q --no-trunc)