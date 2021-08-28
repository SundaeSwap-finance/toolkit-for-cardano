#!/bin/bash

set -eu

DIR="${IPC_CONFIG:=/ipc-config}"
ROOT="${TESTNET_ROOT:=/testnet}"
TESTNET_MAGIC="${TEST_NAME_MAGIC:=31415}"

START_TIME="$(jq -r .startTime ${ROOT}/byron/genesis.json)"

cat "${ROOT}/byron/genesis.json" \
	| jq ".protocolConsts.protocolMagic = ${TESTNET_MAGIC}"	\
	> "${DIR}/byron-genesis.json"
echo "wrote ${DIR}/byron-genesis.json"

cat "${ROOT}/shelley/genesis.json" \
	| jq ".networkMagic = ${TESTNET_MAGIC}"	\
	| jq ".systemStart = \"$(date -d @${START_TIME} "+%Y-%m-%dT%H:%M:%S.000000000Z")\"" \
	> "${DIR}/shelley-genesis.json"
echo "wrote ${DIR}/shelley-genesis.json"

cat "${ROOT}/shelley/genesis.alonzo.json" \
	> "${DIR}/alonzo-genesis.json"
echo "wrote ${DIR}/alonzo-genesis.json"

until [ -f "${DIR}/byron-genesis.hash" ] ; do
	echo "waiting for ${DIR}/byron-genesis.hash ..."
	sleep 1
done
until [ -f "${DIR}/shelley-genesis.hash" ] ; do
	echo "waiting for ${DIR}/shelley-genesis.hash ..."
	sleep 1
done
until [ -f "${DIR}/alonzo-genesis.hash" ] ; do
	echo "waiting for ${DIR}/alonzo-genesis.hash ..."
	sleep 1
done

cat "${ROOT}/configuration.yaml" \
	| yq -o json eval \
	| jq ".ByronGenesisFile   = \"${DIR}/byron-genesis.json\"" \
	| jq ".ByronGenesisHash   = \"$(cat "${DIR}/byron-genesis.hash")\"" \
	| jq ".ShelleyGenesisFile = \"${DIR}/shelley-genesis.json\"" \
	| jq ".ShelleyGenesisHash = \"$(cat "${DIR}/shelley-genesis.hash")\"" \
	| jq ".AlonzoGenesisFile  = \"${DIR}/alonzo-genesis.json\"" \
	| jq ".AlonzoGenesisHash  = \"$(cat "${DIR}/alonzo-genesis.hash")\"" \
	| jq ".SocketPath         = \"${CARDANO_NODE_SOCKET_PATH:=/ipc-node/node.sock}\"" \
	| jq . \
	> ${DIR}/cardano-node.json
echo "wrote ${DIR}/cardano-node.json"

cp "/etc/cardano-db-sync/cardano-db-sync-config.json" "${DIR}/cardano-db-sync-config.json"
echo "wrote ${DIR}/cardano-db-sync-config.json"


sleep 60000

