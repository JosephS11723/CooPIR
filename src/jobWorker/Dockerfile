FROM golang:alpine

WORKDIR /usr/src/app

# move code into the container
COPY . .

# build go code
RUN go mod vendor
RUN go build -o jobWorker

# get weed binary
RUN wget https://github.com/chrislusf/seaweedfs/releases/download/2.99/linux_amd64.tar.gz
RUN tar -xzf linux_amd64.tar.gz
RUN mv weed /usr/bin/weed
RUN rm linux_amd64.tar.gz

# we need to user root to mount with sudo
USER root

# copy bash script to run

# /bin/bash -c "./weed mount; ./jobWorker"

#CMD ["/bin/bash", "-c",  ""]

#CMD ["/bin/bash", "-c ", ";", "./jobWorker"]

#CMD ["/bin/bash", "-c", "./src/run.sh"]

CMD ["./jobWorker"]
