#!/bin/bash

set -eu

export CARDANO_NODE_SOCKET_PATH="${CARDANO_NODE_SOCKET_PATH:=${HOME}/alonzo-testnet/data/node-bft1/node.sock}"

ADDR="$1"
if [ -z "${ADDR}" ] ; then
  echo "usage: fund-treasury.sh <path-to-treasury-addr-file>"
fi

TESTNET_MAGIC="${TESTNET_MAGIC:=31415}"
EXPECTED_ADA="10000000000000000" # amount of ADA expected
CARDANO_NODE_HOME="${CARDANO_NODE_HOME:=${HOME}/alonzo-testnet}"

TX0=$(cardano-cli query utxo --cardano-mode --testnet-magic "${TESTNET_MAGIC}" --whole-utxo | grep "${EXPECTED_ADA}" | head -1 | awk '{print $1}')

cardano-cli transaction build-raw \
  --fee 0 \
  --alonzo-era \
  --tx-in "${TX0}#0" \
  --tx-out "$(cat "${ADDR}")+${EXPECTED_ADA}" \
  --out-file /tmp/fund-user.tx$$

cardano-cli query protocol-parameters --testnet-magic "${TESTNET_MAGIC}" > /tmp/pparams$$.json

FEE=$(cardano-cli transaction calculate-min-fee \
  --tx-body-file /tmp/fund-user.tx$$ \
  --protocol-params-file /tmp/pparams$$.json \
	--tx-in-count 1 \
	--tx-out-count 1 \
	--witness-count 10 | awk '{print $1}')

ADA=$(echo "${EXPECTED_ADA} - ${FEE}" | bc)

cardano-cli transaction build-raw \
  --fee "${FEE}" \
  --alonzo-era \
  --tx-in "${TX0}#0" \
  --tx-out "$(cat "${ADDR}")+${ADA}" \
  --out-file /tmp/fund-user.tx$$

cardano-cli transaction sign \
  --tx-body-file /tmp/fund-user.tx$$ \
  --signing-key-file "${CARDANO_NODE_HOME}/addresses/pool-owner1.skey" \
  --signing-key-file "${CARDANO_NODE_HOME}/node-bft1/shelley/operator.skey" \
  --signing-key-file "${CARDANO_NODE_HOME}/node-bft2/shelley/operator.skey" \
  --signing-key-file "${CARDANO_NODE_HOME}/node-pool1/owner.skey" \
  --signing-key-file "${CARDANO_NODE_HOME}/node-pool1/shelley/operator.skey" \
  --signing-key-file "${CARDANO_NODE_HOME}/shelley/genesis-keys/genesis2.skey" \
  --signing-key-file "${CARDANO_NODE_HOME}/shelley/genesis-keys/genesis1.skey" \
  --signing-key-file "${CARDANO_NODE_HOME}/shelley/utxo-keys/utxo1.skey" \
  --signing-key-file "${CARDANO_NODE_HOME}/shelley/delegate-keys/delegate1.skey" \
  --signing-key-file "${CARDANO_NODE_HOME}/shelley/delegate-keys/delegate2.skey" \
  --out-file /tmp/fund-user.sign$$

cardano-cli transaction submit \
  --cardano-mode \
  --testnet-magic "${TESTNET_MAGIC}" \
  --tx-file /tmp/fund-user.sign$$

rm -f /tmp/fund-user.*
