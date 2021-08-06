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
  if [ "${1}" != "minikube" ] ; then
    exit 1
  fi
  "$2"/scripts/server.sh "${1}" "${2}"
  "$2"/scripts/client.sh "${1}" "${2}"
}
main "${1}" "${2}"

