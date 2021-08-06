# Real time chat

## Description
- Users can use this application to chat to each other
- To chat, user must first enter his/her name
- Users should see all messages history

## Install requirements
- Kubernetes (minikube if local)
- Helm
- golang 1.16
- nodeJS

## Technical used
- WebSocket is used for transferring messages
- Mysql as the db to store room and message data
- Gin-gonic for http api
- Angular 12 as client side.

## Flow
- User creates new room
- By default, user will be assigned to have a random name
- User access to the room, he/she can copy the url and send to friends

## Live Demo

http://35.193.251.224/

## Project structure

- charts: contains Helm templates for Mysql, server and client
- client: contains client code using Angular 12
- internal: internal packages of server
- scripts: contains deployment scripts
- main.go: main function of server

## TODO:

- Apply script to deploy eveything to minikube or GKE.
- Validate if a roomId is in database or not.
- Load more previous messages when user scrolls up
- scroll bar always auto scroll to the bottom when a new message is added