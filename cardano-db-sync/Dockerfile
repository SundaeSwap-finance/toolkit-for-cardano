FROM sundaeswap/cardano-db-sync:fix-fork-at-0

ADD entrypoint.sh /usr/sbin/entrypoint.sh

ENTRYPOINT [ \
  "/usr/sbin/entrypoint.sh", "cardano-db-sync-extended", \
  "--config",      "/ipc-config/cardano-db-sync-config.json", \
  "--socket-path", "/ipc-node/node.sock", \
  "--state-dir",   "/data", \
  "--schema-dir",  "/nix/store/z2vk0sddkm94dvinvpvg1dff4yr4r1qf-schema" \
]
