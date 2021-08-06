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
  REPOSITORY=""
  if  [ "$1" != "minikube" ] ; then
    REPOSITORY="$1"
  fi
  BASE_DIRECTORY="$2"
  # install mysql
  kubectl create secret generic kube-dev-mysql-password --from-literal=username=kube_dev_mysql --from-literal=password=123456
  helm install -f "$BASE_DIRECTORY"/charts/mysql/mysql.yaml mysql bitnami/mysql

  # build docker images
#  if  [ "$REPOSITORY" = "minikube" ] ; then
#    minikube docker-env
#    eval '$(minikube -p minikube docker-env)'
#  fi

  docker build -t "$REPOSITORY"chat_demo_server .

  # push built image to container image registry
  if  [ "$1" != "minikube" ] ; then
    docker push "$REPOSITORY"chat_demo_server
  fi

  # install server by using helm chart
  helm install backend -f "$BASE_DIRECTORY"/charts/server.dev.yaml "$BASE_DIRECTORY"/charts/server

  sleep 10

  # the deployment will be created as backend-live-chat-server
  # expose the port
  if  [ "$1" != "minikube" ] ; then
    kubectl expose deployment backend-live-chat-server --type=LoadBalancer --name=live-chat-server-expose-service
    printf "Wait until live-chat-server-expose-service is finished generating external IP and Port.\n Copy the and paste it to client/environment.prod.ts"
  elif [ "$1" = "minikube" ] ; then
    kubectl port-forward svc/backend-live-chat-server 5000:5000 &
  fi
  kubectl get svc

}
main "${1}" "${2}"

