#!/bin/bash
if [[ $1 == *"k"* ]]; then
  docker container kill $(docker ps -q)
elif [[ $1 == *"rebuild"* ]]; then
  docker build -t server ./Kontejner_SSH/
  docker build -t student ./Kontejner_Student/
else
  docker run -p 2222:1234 -p 1234:1234 -p 4321:1234 -v /var/run/docker.sock:/var/run/docker.sock -itd server
fi