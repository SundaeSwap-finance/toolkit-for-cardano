cardano-db-sync
------------------

As of this writing, `sundaeswap/cardano-db-sync` is built from an unmerged branch,
`https://github.com/input-output-hk/cardano-db-sync/pull/766` that resolves an
issue with spinning up `cardano-db-sync` from on a testnet that forks to Shelley
in epoch 0.

```
"TestShelleyHardForkAtEpoch": 0
```

This is a temporary workaround while we wait for the PR to be merged.
