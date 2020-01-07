import os
import sys
import subprocess


def main():
    # Find all of the entity_genesis.json files and node_genesis.json files
    unpacked_entities_path = os.path.abspath(sys.argv[1])

    # Hacky overrides for running locally.
    output_path = os.environ.get("GENESIS_OUTPUT_PATH", "/tmp/genesis.json")
    staking_path = os.environ.get("STAKING_GENESIS_PATH", "/tmp/staking.json")
    oasis_node_path = os.environ.get("OASIS_NODE_PATH", '/tmp/oasis-node')

    genesis_command = [
        oasis_node_path, "genesis", "init",
        "--genesis.file", output_path,
        "--chain.id", "sometest-chain-id",
        "--staking", staking_path,
        "--epochtime.tendermint.interval", "200",
        "--consensus.tendermint.timeout_commit", "5s",
        "--consensus.tendermint.empty_block_interval", "0s",
        "--consensus.tendermint.max_tx_size", "32kb",
        "--consensus.backend", "tendermint"
    ]

    for entity_name in os.listdir(unpacked_entities_path):
        if os.path.isfile(os.path.join(unpacked_entities_path, entity_name)):
            continue
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
