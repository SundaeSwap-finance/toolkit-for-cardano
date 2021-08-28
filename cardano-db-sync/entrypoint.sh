#!/bin/bash

set -eu

CONFIG="${CONFIG:=/ipc-config/cardano-db-sync-config.json}"
DIR="/ipc-config"


#-- inline functions -----------------------------------------


# wait until postgresql is responding
#
function get_env() {
	KEY="$1"
	DEFAULT="$2"

	FILE="$(env | grep -E "^${KEY}=" | awk -F= '{print $2}')"
	if [ ! -z "${FILE}" ] && [ -f "${FILE}" ] ; then
		cat "${FILE}"
		return
	fi

	echo "${DEFAULT}"
}


#-- begin script ---------------------------------------------


# wait until configuration file is present
#
until [ -f "${CONFIG}" ] ; do
	>&2 echo "waiting for cardano-db-sync configuration ..."
	sleep 1
done
  

# wait until postgres is available
#
POSTGRES_DB="${POSTGRES_DB:=$(get_env POSTGRES_DB_FILE 'cardano_toolkit')}"
POSTGRES_HOST="${POSTGRES_HOST:=$(get_env POSTGRES_HOST_FILE '127.0.0.1')}"
POSTGRES_PORT="${POSTGRES_PORT:=$(get_env POSTGRES_PORT_FILE '5432')}"
POSTGRES_PASSWORD="${POSTGRES_PASSWORD:=$(get_env POSTGRES_PASSWORD_FILE 'password')}"
POSTGRES_USER="${POSTGRES_USER:=$(get_env POSTGRES_USER_FILE 'postgres')}"

until PGPASSWORD="${POSTGRES_PASSWORD}" psql -h "${POSTGRES_HOST}" -p "${POSTGRES_PORT}" -U "${POSTGRES_USER}" -c '\q'; do
	>&2 echo "waiting for postgres ..."
  sleep 1
done


# generate PGPASSFILE
#
export PGPASSFILE="${PGPASSFILE:=/etc/cardano-db-sync/.pgpass}"
mkdir -p "$(dirname "${PGPASSFILE}")"

cat <<EOF > "${PGPASSFILE}"
${POSTGRES_HOST}:${POSTGRES_PORT}:${POSTGRES_DB}:${POSTGRES_USER}:${POSTGRES_PASSWORD}
EOF
chmod 0600 "${PGPASSFILE}"


exec "$@"

