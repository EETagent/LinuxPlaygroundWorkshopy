FROM ubuntu:20.04

ENV TZ=Europe/Prague
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get upgrade
RUN apt-get update
RUN apt-get install -y figlet
RUN apt-get -y install openssh-server docker.io sudo

RUN mkdir -p /run/sshd
RUN rm /etc/legal
RUN rm /etc/ssh/sshd_config
COPY sshd_config /etc/ssh/
RUN service ssh reload

RUN chmod -x /etc/update-motd.d/*
COPY varovani /etc/update-motd.d/
RUN chmod +x /etc/update-motd.d/varovani
COPY banner /opt/
COPY sprava.sh /opt/
RUN chmod +x /opt/sprava.sh

RUN rm /etc/sudoers
COPY sudoers /etc/

RUN useradd -d /home/madeta/ -m -p ssps -s /opt/sprava.sh madeta
RUN echo "madeta:ssps" | chpasswd
WORKDIR /home/madeta
RUN chown -R root:root /home/madeta

COPY entrypoint.sh /root/
RUN chmod +x /root/entrypoint.sh 

ENTRYPOINT /root/entrypoint.sh

#Spouštět s -v /var/run/docker.sock:/var/run/docker.sock
