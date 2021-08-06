#!/bin/sh

#IFS=$'\n\t'
#set -eou pipefail

IFS=$'\n\t'
set -eou pipefail

if  [ "$#" -ne 2 ] || [ "${1}" = '-h' ] || [ "${1}" = '--help' ] ; then
  cat >&2 <<"EOF"
USAGE:
  server.sh REPOSITORY BASE_DIRECTORY
  if REPOSITORY is set to minikube, it means that we are using local cluster
  then reset docker-env and don't push image to docker registry
EOF
  exit 1
fi

main(){
  REPOSITORY="asia.gcr.io/api-project-695028345372"
  if  [ "$1" != "minikube" ] ; then
    REPOSITORY="$1"
  fi
  BASE_DIRECTORY="$2"

  # build docker images
#  if  [ "$REPOSITORY" = "minikube" ] ; then
#    minikube docker-env
#    eval '$(minikube -p minikube docker-env)'
#  fi

  docker build -t "$REPOSITORY"/chat_demo_client "$BASE_DIRECTORY/client"

  # push built image to container image registry
  if  [ "$REPOSITORY" != "minikube" ] ; then
    docker push "$REPOSITORY"/chat_demo_client
  fi

  # install server by using helm chart
  helm install live-chat "$BASE_DIRECTORY"/charts/client

  sleep 10

  # the deployment will be created as backend-live-chat-server
  # expose the port
  if  [ "$1" != "minikube" ] ; then
    kubectl expose deployment live-chat-client --type=LoadBalancer --name=live-chat-client-expose-service
    printf "Wait until live-chat-server-expose-service is finished generating external IP and Port.\n Copy the and paste it to client/environment.prod.ts"
  elif [ "$1" = "minikube" ] ; then
    kubectl port-forward svc/live-chat-client 4200:80 &
  fi
  kubectl get svc

}
main "${1}" "${2}"

