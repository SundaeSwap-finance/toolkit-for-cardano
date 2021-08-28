#!/bin/bash

#---------------------------------------------
# hash the genesis blocks

set -eu

CONFIG="${CONFIG:=/ipc-config/cardano-db-sync-config.json}"
DIR="${IPC_CONFIG:=/ipc-config}"


# wait until configuration file is present
#
until [ -f "${CONFIG}" ] ; do
	>&2 echo "waiting for cardano-db-sync configuration ..."

	if [ -f "${DIR}/byron-genesis.json" ] && [ ! -f "${DIR}/byron-genesis.hash" ] ; then
		cardano-cli byron genesis print-genesis-hash --genesis-json "${DIR}/byron-genesis.json" > "${DIR}/byron-genesis.hash"
		echo "wrote ${DIR}/byron-genesis.hash"
	fi

	if [ -f "${DIR}/shelley-genesis.json" ] && [ ! -f "${DIR}/shelley-genesis.hash" ] ; then
		cardano-cli genesis hash --genesis "${DIR}/shelley-genesis.json" > "${DIR}/shelley-genesis.hash"
		echo "wrote ${DIR}/shelley-genesis.hash"
	fi

	if [ -f "${DIR}/alonzo-genesis.json" ] && [ ! -f "${DIR}/alonzo-genesis.hash" ] ; then
		cardano-cli genesis hash --genesis "${DIR}/alonzo-genesis.json" > "${DIR}/alonzo-genesis.hash"
		echo "wrote ${DIR}/alonzo-genesis.hash"
	fi

	sleep 1
done