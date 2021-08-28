#!/bin/bash

#---------------------------------------------
# generates any number of Shelley addresses
#
# Usage: generate-address.sh treasury
#

set -eu

TESTNET_MAGIC="${TESTNET_MAGIC:=42}"

for ADDR in $*; do

  # Payment address keys
  cardano-cli address key-gen \
      --verification-key-file "${ADDR}.vkey" \
      --signing-key-file      "${ADDR}.skey"

  # Stake address keys
  cardano-cli stake-address key-gen \
      --verification-key-file "${ADDR}-stake.vkey" \
      --signing-key-file      "${ADDR}-stake.skey"

  # Payment addresses
  cardano-cli address build \
      --payment-verification-key-file "${ADDR}.vkey" \
      --stake-verification-key-file "${ADDR}-stake.vkey" \
      --testnet-magic "${TESTNET_MAGIC}" \
      --out-file "${ADDR}.addr"

  # Stake addresses
  cardano-cli stake-address build \
      --stake-verification-key-file "${ADDR}-stake.vkey" \
      --testnet-magic "${TESTNET_MAGIC}" \
      --out-file "${ADDR}-stake.addr"

  # Stake addresses registration certs
  cardano-cli stake-address registration-certificate \
      --stake-verification-key-file "${ADDR}-stake.vkey" \
      --out-file "${ADDR}-stake.reg.cert"

done
