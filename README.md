# Real time chat

## Description
- Users can use this application to chat to each other
- To chat, user must first enter his/her name
- Users should see all messages history


## Functional Requirements
- WebSocket to announce new messages
- Mysql is used to store user's messages since we need consistence data
- Gin-gonic is used as a http server
- Angular 12 as client side.

## Flow
- User creates new room
- By default, user will be assigned to have a random name
- User access to the room, he/she can copy the url and send to others

## Install requirements
- Kubernetes (minikube if local)
- Helm
- golang 1.16
- nodeJS

## Live Demo

http://35.193.251.224/

## TODO:
- Apply script to deploy eveything to minikube or GKE.
- Validate if a roomId is in database or not.