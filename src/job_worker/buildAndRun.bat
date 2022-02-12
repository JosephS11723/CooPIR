docker stop job_worker
docker rm job_worker
docker rmi job_worker
docker build -t job_worker .
docker run --name job_worker job_worker