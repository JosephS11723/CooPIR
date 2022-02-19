docker stop api_test
docker rm api_test
docker rmi api_test
docker build -t api_test .
docker run --name api_test -p 8080:8080 api_test