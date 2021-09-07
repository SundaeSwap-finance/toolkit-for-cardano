#!/bin/bash

#---------------------------------------------
# register the stake pool for delegators
#
# Usage: generate-address.sh treasury
#

set -eu
set -x

TESTNET_MAGIC="${TESTNET_MAGIC:=42}"
CARDANO_NODE_HOME="${CARDANO_NODE_HOME:=${HOME}/alonzo-testnet}"
# Note: user1 in ~/alonzo-testnet/addresses and in ~/sundaeswap/data are not the same
# The latter is used for the treasury
TREASURY_ADDR_LOC="${TREASURY_ADDR_LOC:=${HOME}/sundaeswap/data/user1}"
TREASURY_ADDR=$(cat "${TREASURY_ADDR_LOC}.addr")
TREASURY_KEY="${TREASURY_ADDR_LOC}.skey"

OWNER_ADDR=$(cat "${CARDANO_NODE_HOME}/addresses/pool-owner1.addr")
OWNER_KEY="${CARDANO_NODE_HOME}/addresses/pool-owner1.skey"
OWNER_STAKE_VKEY="${CARDANO_NODE_HOME}/addresses/pool-owner1-stake.vkey"
OWNER_STAKE_SKEY="${CARDANO_NODE_HOME}/addresses/pool-owner1-stake.skey"

PROTOCOL_PARAMS=$(mktemp -p /tmp pparams.XXXX.json)
cardano-cli query protocol-parameters --testnet-magic "${TESTNET_MAGIC}" > ${PROTOCOL_PARAMS}

# To register a stakepool, we need to:
#  - Register the owner stake address
#  - Include the stake-pool registration cert
#  - Include a delegation cert to cover the pledge

# The stake-pool registration cert is pre-built at
# CARDANO_NODE_HOME/node-pool1/registration.cert

# The owner stake registration cert is pre-built at
# CARDANO_NODE_HOME/addresses/pool-owner1-stake.reg.cert

# We need to register the pool owner stake address
# And to do that, the owner account needs funds
# Usually we'd do gql, but I don't think that's running yet
TX_RAW=$(mktemp -p /tmp fund-pool-owner.XXXX.tx)
TX_SIGNED=$(mktemp -p /tmp fund-pool-owner.XXXX.signed)
TREASURY_UTXO=$(cardano-cli query utxo --testnet-magic "${TESTNET_MAGIC}" --address ${TREASURY_ADDR} | head -3 | tail -1 | awk '{print $1"#"$2}')
AMT="10000000"

cardano-cli transaction build \
  --testnet-magic "${TESTNET_MAGIC}" \
  --alonzo-era \
  --tx-in "${TREASURY_UTXO}" \
  --tx-out "${OWNER_ADDR}+${AMT}" \
  --change-address "${TREASURY_ADDR}" \
  --out-file "${TX_RAW}"

UTXO="$(cardano-cli transaction txid --tx-body-file ${TX_RAW})#1"

cardano-cli transaction sign \
  --testnet-magic "${TESTNET_MAGIC}" \
  --tx-body-file "${TX_RAW}" \
  --signing-key-file "${TREASURY_KEY}" \
  --out-file "${TX_SIGNED}"

cardano-cli transaction submit \
  --testnet-magic "${TESTNET_MAGIC}" \
  --tx-file "${TX_SIGNED}"

rm -f "${TX_RAW}"
rm -f "${TX_SIGNED}"

# Wait for the transaction to settle
sleep 3s

# Next, register the stake address
TX_RAW=$(mktemp -p /tmp register-stake.XXXX.tx)
TX_SIGNED=$(mktemp -p /tmp register-stake.XXXX.signed)

# Estimate the fees, since `transaction build` is broken with certs
cardano-cli transaction build-raw \
  --alonzo-era \
  --tx-in "${UTXO}" \
  --tx-out "${OWNER_ADDR}+${AMT}" \
  --fee 0 \
  --certificate-file "${CARDANO_NODE_HOME}/addresses/pool-owner1-stake.reg.cert" \
  --out-file "${TX_RAW}"

FEE=$(
  cardano-cli transaction calculate-min-fee \
    --tx-body-file "${TX_RAW}" \
    --protocol-params-file "${PROTOCOL_PARAMS}" \
    --tx-in-count 1 \
    --tx-out-count 1 \
    --witness-count 2 | awk '{print $1}'
)

AMT=$(expr ${AMT} - ${FEE})

cardano-cli transaction build-raw \
  --alonzo-era \
  --tx-in "${UTXO}" \
  --tx-out "${OWNER_ADDR}+${AMT}" \
  --fee "${FEE}" \
  --certificate-file "${CARDANO_NODE_HOME}/addresses/pool-owner1-stake.reg.cert" \
  --out-file "${TX_RAW}"

UTXO="$(cardano-cli transaction txid --tx-body-file ${TX_RAW})#0"

cardano-cli transaction sign \
  --testnet-magic "${TESTNET_MAGIC}" \
  --tx-body-file "${TX_RAW}" \
  --signing-key-file "${OWNER_KEY}" \
  --signing-key-file "${CARDANO_NODE_HOME}/addresses/pool-owner1-stake.skey" \
  --out-file "${TX_SIGNED}"

cardano-cli transaction submit \
  --testnet-magic "${TESTNET_MAGIC}" \
  --tx-file "${TX_SIGNED}"

rm -f "${TX_RAW}"
rm -f "${TX_SIGNED}"

# Wait for the transaction to settle
sleep 3s

# Now, register the pool with a pledge
CERT_FILE=$(mktemp -p /tmp pool-pledge.XXXX.cert)
TX_RAW=$(mktemp -p /tmp register-pool.XXXX.tx)
TX_SIGNED=$(mktemp -p /tmp register-pool.XXXX.signed)

cardano-cli stake-address delegation-certificate \
  --stake-verification-key-file "${OWNER_STAKE_VKEY}" \
  --cold-verification-key-file "${CARDANO_NODE_HOME}/node-pool1/shelley/operator.vkey" \
  --out-file "${CERT_FILE}"

# Estimate the fees, since `transaction build` is broken with certs
cardano-cli transaction build-raw \
  --alonzo-era \
  --tx-in "${UTXO}" \
  --tx-out "${OWNER_ADDR}+${AMT}" \
  --fee 0 \
  --certificate-file "${CARDANO_NODE_HOME}/node-pool1/registration.cert" \
  --certificate-file "${CERT_FILE}" \
  --out-file "${TX_RAW}"

FEE=$(
  cardano-cli transaction calculate-min-fee \
    --tx-body-file "${TX_RAW}" \
    --protocol-params-file "${PROTOCOL_PARAMS}" \
    --tx-in-count 1 \
    --tx-out-count 1 \
    --witness-count 3 | awk '{print $1}'
)

AMT=$(expr ${AMT} - ${FEE})

cardano-cli transaction build-raw \
  --alonzo-era \
  --tx-in "${UTXO}" \
  --tx-out "${OWNER_ADDR}+${AMT}" \
  --fee ${FEE} \
  --certificate-file "${CARDANO_NODE_HOME}/node-pool1/registration.cert" \
  --certificate-file "${CERT_FILE}" \
  --out-file "${TX_RAW}"

cardano-cli transaction sign \
  --testnet-magic "${TESTNET_MAGIC}" \
  --tx-body-file "${TX_RAW}" \
  --signing-key-file "${OWNER_KEY}" \
  --signing-key-file "${OWNER_STAKE_SKEY}" \
  --signing-key-file "${CARDANO_NODE_HOME}/node-pool1/shelley/operator.skey" \
  --out-file "${TX_SIGNED}"

cardano-cli transaction submit \
  --testnet-magic "${TESTNET_MAGIC}" \
  --tx-file "${TX_SIGNED}"

rm -f "${TX_RAW}"
rm -f "${TX_SIGNED}"
rm -f "${PROTOCOL_PARAMS}"
# The pool should now be ready for delegation