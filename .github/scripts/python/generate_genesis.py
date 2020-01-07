import os
import sys
import subprocess


def main():
    # Find all of the entity_genesis.json files and node_genesis.json files
    unpacked_entities_path = os.path.abspath(sys.argv[1])

    genesis_command = [
        "oasis-node",
        "--genesis.file", "/tmp/genesis.json",
        "--chain.id", "sometest-chain-id",
        "--staking", "/tmp/staking.json",
        "--epochtime.tendermint.interval", "200",
        "--epochtime.tendermint.timeout_commit", "5s",
        "--consensus.tendermint.empty_block_interval", "0s",
        "--consensus.tendermint.max_tx_size", "32kb",
        "--consensus.tendermint.backend", "tendermint"
    ]

    for entity_name in os.listdir(unpacked_entities_path):
        genesis_command.extend([
            "--entity", os.path.join(unpacked_entities_path,
                                     entity_name, "entity/entity_genesis.json"),
            "--node", os.path.join(unpacked_entities_path,
                                   entity_name, "node/node_genesis.json"),
        ])

    # Run genesis command
    subprocess.check_call(genesis_command)


if __name__ == '__main__':
    main()
