# Deploy to GKE

1. Deploy server first with 2 parameters:
- Your docker repository (asia.gcr.io) combine with your projectId (api-project-695028345372)
- Path lead to server's Dockerfile

```
./server.sh asia.gcr.io/api-project-695028345372 ./
```

The script will deploy backend to kubernetes cluster and expose it to an external IP address with port 5000
After the script finishes, wait until `expose` service finishes generating the IP address, copy it and paste to field `apiUrl` within `client/src/environments/environment.prod.ts`

2. Deploy client with 2 parameters
- Your docker repository (asia.gcr.io) combine with your projectId (api-project-695028345372)
- Path lead to server's Dockerfile

```
./client.sh asia.gcr.io/api-project-695028345372 ./
```

Client is also deployed the same way as server.
Wait until `expose` service finishes, we can now use the URL to access the application