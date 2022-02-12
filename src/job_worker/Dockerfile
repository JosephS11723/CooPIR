FROM python:3.9.10-bullseye

# get weed binary
RUN wget https://github.com/chrislusf/seaweedfs/releases/download/2.88/linux_amd64.tar.gz
RUN tar -xzf linux_amd64.tar.gz
RUN mv weed /usr/bin/weed
RUN rm linux_amd64.tar.gz

# move src into container
ADD src /src

# copy bash script to run


CMD ["/bin/bash", "-c", "./src/run.sh"]