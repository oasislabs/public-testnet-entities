import os
import sys
import subprocess


def add_entities_from_directory(genesis_command, entity_dir, add_nodes):
    for entity_name in os.listdir(unpacked_entities_path):
        if os.path.isfile(os.path.join(unpacked_entities_path, entity_name)):
            continue
        genesis_command.append([
            '--entity', os.path.join(unpacked_entities_path,
                                     entity_name, 'entity/entity_genesis.json'),
        ])

        if add_nodes:
            genesis_command.append(
                '--node', os.path.join(unpacked_entities_path,
                                       entity_name, 'node/node_genesis.json'),
            )


def main():
    # Find all of the entity_genesis.json files and node_genesis.json files
    unpacked_entities_path = os.path.abspath(sys.argv[1])

    # Hacky overrides for running locally.
    output_path = os.environ.get('GENESIS_OUTPUT_PATH', '/tmp/genesis.json')
    staking_path = os.environ.get('STAKING_GENESIS_PATH', '/tmp/staking.json')
    oasis_node_path = os.environ.get('OASIS_NODE_PATH', '/tmp/oasis-node')
    oasis_test_only = os.environ.get('OASIS_TEST_ONLY', 'false') == 'true'

    genesis_command = [
        oasis_node_path, 'genesis', 'init',
        '--genesis.file', output_path,
        '--chain.id', 'the-quest-test-2019-01-09',
        '--staking', staking_path,
        '--epochtime.tendermint.interval', '200',
        '--consensus.tendermint.timeout_commit', '5s',
        '--consensus.tendermint.empty_block_interval', '0s',
        '--consensus.tendermint.max_tx_size', '32kb',
        '--consensus.backend', 'tendermint',
        '--scheduler.max_validators', '100',
        '--scheduler.max_per_entity', '1',
        '--scheduler.min_validators', '3',
    ]

    add_entities_from_directory(
        genesis_command, unpacked_entities_path, add_nodes=not oasis_test_only)

    if oasis_test_only:
        add_entities_from_directory(
            genesis_command, test_entities_path, add_nodes=oasis_test_only)

    # Run genesis command
    subprocess.check_call(genesis_command)


if __name__ == '__main__':
    main()
