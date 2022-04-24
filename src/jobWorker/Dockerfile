FROM ubuntu:18.04

WORKDIR /usr/src/app

# we need to user root to mount with sudo
USER root

# install dependencies
RUN apt-get update
RUN apt-get install -y libfuse-dev wget fuse

# install go
RUN wget "https://go.dev/dl/go1.18.1.linux-amd64.tar.gz"
RUN tar -C /usr/local -xzf go1.18.1.linux-amd64.tar.gz
RUN rm go1.18.1.linux-amd64.tar.gz

# chmod go files
RUN chmod -R 777 /usr/local/go/bin

# move code into the container
COPY . .

# build go code
RUN /usr/local/go/bin/go mod vendor
RUN /usr/local/go/bin/go build -o jobWorker

# get weed binary
RUN wget https://github.com/chrislusf/seaweedfs/releases/download/2.99/linux_amd64.tar.gz
RUN tar -xzf linux_amd64.tar.gz
RUN rm linux_amd64.tar.gz

# copy bash script to run

# /bin/bash -c "./weed mount; ./jobWorker"

#CMD ["/bin/bash", "-c",  ""]

#CMD ["/bin/bash", "-c ", ";", "./jobWorker"]

#CMD ["/bin/bash", "-c", "./src/run.sh"]

CMD ["./jobWorker"]
