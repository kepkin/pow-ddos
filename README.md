Intro
=====

This is a one day coding challenge for a job interview.


Task
====


Test task for Server Engineer

Design and implement “Word of Wisdom” tcp server.  
 • TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.  
 • The choice of the POW algorithm should be explained.  
 • After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.  
 • Docker file should be provided both for the server and for the client that solves the POW challenge


Solution
========

As we want to protect from DDOS attacks, cpu bound POW are less favourable. Alternatives are:

 - memory bound POWs
    - MBound
    - Cuckoo cycle

 - Guided tour puzzle protocol


I couldn't find any GO implementation for memory bound POWs that I could reuse at the moment of writting. Guided tour puzzle is pretty simple protocol wihtout any complex math, so it was an ideal choice for one day project.


What is missing
---------------

As it's a one day project, I missed the following features which I would love to add:

 - get secrets from shared storage
 - auto rotation of secrets
 - the protocol may look somewhat inconsistent, my main approach was:
    - for server request you should pass hash in Header to not interfere with other API parameters
    - guide server doesn't answer anything from hashes, so it recieves and answers in body


How to run
==========

For testing purposes docker compose file is available with following configuration

 - number of Guides servers 2
 - the tour length is hard coded to 5


 docker-compose up
 docker-compose run client ./example-client --server http://server:8080 --guides http://guide1:8080 http://guide2:8080
