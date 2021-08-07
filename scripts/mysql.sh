#!/bin/sh

IFS=$'\n\t'
set -eou pipefail

if  [ "$#" -ne 1 ] || [ "${1}" = '-h' ] || [ "${1}" = '--help' ] ; then
  cat >&2 <<"EOF"
USAGE:
  mysql.sh BASE_DIRECTORY
EOF
  exit 1
fi

main(){
  BASE_DIRECTORY="$1"
  # install mysql
  kubectl create secret generic kube-dev-mysql-password --from-literal=username=kube_dev_mysql --from-literal=password=123456
  helm install -f "$BASE_DIRECTORY"/charts/mysql/mysql.yaml mysql bitnami/mysql
}
main "${1}"

