FROM ubuntu:20.04

ENV TZ=Europe/Prague
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN apt-get update
RUN apt-get install dnsutils curl wget sudo inetutils-ping binutils -y

RUN useradd -d /home/eidam/ -m -p eidam -s /bin/bash eidam
RUN echo "eidam:eidam" | chpasswd
RUN usermod -aG sudo eidam

COPY bin/kontrolka /bin/kontrolka
COPY task/kontrolka-zadani.txt /home/eidam/README_kontrolka.txt
RUN chmod +x /bin/kontrolka
RUN ln -s /bin/kontrolka /home/eidam/kontrolka

COPY bin/ukladatko /bin/ukladatko
COPY task/ukladatko-zadani.txt /home/eidam/README_ukladatko.txt
RUN chmod +x /bin/ukladatko
RUN ln -s /bin/ukladatko /home/eidam/ukladatko

COPY bin/print_me /bin/print_me
RUN chmod +x /bin/print_me
RUN ln -s /bin/print_me /home/eidam/print_me

COPY bin/neofetch /bin/neofetch
RUN chmod +x /bin/neofetch
RUN ln -s /bin/neofetch /home/eidam/neofetch

WORKDIR /home/eidam
USER eidam
CMD /bin/bash
