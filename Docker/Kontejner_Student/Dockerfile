FROM ubuntu:20.04

ENV TZ=Europe/Prague
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get update
RUN apt-get install curl wget sudo -y

RUN useradd -d /home/eidam/ -m -p eidam -s /bin/bash eidam
RUN echo "eidam:eidam" | chpasswd
RUN usermod -aG sudo eidam

COPY checker /bin/checker
COPY checker-navod.txt /home/eidam/README.txt
RUN chmod +x /bin/checker
RUN ln -s /bin/checker /home/eidam/checker

COPY print_me /bin/print_me
RUN chmod +x /bin/print_me
RUN ln -s /bin/print_me /home/eidam/print_me

WORKDIR /home/eidam
USER eidam
CMD /bin/bash
