#!/bin/sh

#IFS=$'\n\t'
#set -eou pipefail

IFS=$'\n\t'
set -eou pipefail

if  [ "$#" -ne 3 ] || [ "${1}" = '-h' ] || [ "${1}" = '--help' ] ; then
  cat >&2 <<"EOF"
USAGE:
  server.sh REPOSITORY BASE_DIRECTORY IS_UPDATE_IMAGE
  if REPOSITORY is set to minikube, it means that we are using local cluster
  then reset docker-env and don't push image to docker registry
EOF
  exit 1
fi

main(){
  REPOSITORY="$1"
  BASE_DIRECTORY="$2"
  IS_UPDATE_IMAGE="${3}"
  docker build -t "$REPOSITORY"/chat_demo_client "$BASE_DIRECTORY/client/"
  docker push "$REPOSITORY"/chat_demo_client
  # wait a bit to let gcloud update the image
  # TODO: find a way to get the image status rather than this
  i=1
  while [ "$i" -ne 10 ]
  do
      echo "waiting..."
      i=$((i + 1))
      sleep 1
  done
  if [ "$IS_UPDATE_IMAGE" = "true" ]; then
    kubectl scale deployment live-chat-client --replicas=0 -n default
    sleep 2
    kubectl scale deployment live-chat-client --replicas=1 -n default
  else
    # install server by using helm chart
    helm install live-chat "$BASE_DIRECTORY"/charts/client
    sleep 5
    # the deployment will be created as backend-live-chat-server
    # expose the port
    kubectl expose deployment live-chat-client --type=LoadBalancer --name=live-chat-client-expose-service
    printf "Wait until live-chat-server-expose-service is finished generating external IP and Port.\n Copy the and paste it to client/environment.prod.ts"
    kubectl get svc
  fi
}
main "${1}" "${2}" "${3}"

