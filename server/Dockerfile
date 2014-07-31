# Scooby
#
# VERSION 0.1

FROM ubuntu:12.04
MAINTAINER tylrtrmbl <taylor@thenewtricks.com>

# make sure the package repository is up to date
RUN echo "deb http://archive.ubuntu.com/ubuntu precise main universe" > /etc/apt/sources.list
RUN apt-get -q -y update
RUN apt-get -q -y upgrade

# Install nginx, Python, and supervisord
RUN apt-get -q -y install nginx python-setuptools python-dev
RUN easy_install supervisor

# Install specific common tools
RUN apt-get -q -y install curl

# Install golang, git, and hg
ENV GOPATH /root/go
RUN curl -L https://go.googlecode.com/files/go1.2.linux-amd64.tar.gz | tar xvz --directory=/usr/local
ENV PATH /usr/local/go/bin:$PATH
RUN apt-get -q -y install git mercurial

# Build scooby
ADD . $GOPATH/src/github.com/teslaworksumn/scooby
RUN cd $GOPATH/src/github.com/teslaworksumn/scooby && go get
RUN go install github.com/teslaworksumn/scooby

# Configure nginx
ADD config/nginx_config /etc/nginx/sites-available/default

# Configure supervisord
ADD config/supervisord_programs.conf /etc/supervisord_programs.conf
RUN echo_supervisord_conf > /etc/supervisord.conf
RUN cat /etc/supervisord_programs.conf >> /etc/supervisord.conf

# Add run script
ADD scripts/run.sh /root/run.sh
RUN chmod a+x /root/run.sh

EXPOSE 80

CMD ["/root/run.sh"]
