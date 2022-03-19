#!/bin/sh
# Tested with Bee version: 1.4.3-9546fedb

[[ -e "$(pwd)/.env" ]] && source $(pwd)/.env

if [[ -z "${BEE_BIN_PATH}" ]]; then
    bee_path="/usr/bin"
else
    bee_path=${BEE_BIN_PATH}
fi
bee_data_root=$(pwd)/bee_data_root
mkdir -p ${bee_data_root}

run_dfs() {
    dfs server --dataDir=${bee_data_root}/dfs --beePort=1633
}

run_bee_test() {
    echo "Using swap endpoint: ${BEE_SWAP_ENDPOINT}"
    [[ "$#" -ne 1 ]] && { echo "Usage : ./run_node -test NODE_ID"; exit 1; }
    node=$1
    port=$(($node+15))
    api_port=${port}33
    p2p_port=${port}34
    debug_port=${port}35
    data_dir="${bee_data_root}/bee_$node"

    echo "Starting node-$node"
    ${bee_path}/bee start \
        --api-addr=:${api_port} \
        --debug-api-enable=true \
        --debug-api-addr="localhost:${debug_port}" \
        --data-dir=${data_dir} \
        --verbosity=5 \
        --bootnode="" \
        --full-node=true \
        --network-id=10 \
        --mainnet=false \
        --warmup-time=0 \
        --resolver-options=${BEE_ENS_RESOLVER} \
        --swap-initial-deposit="1000000000000000" \
        --p2p-addr=:${p2p_port} \
        --swap-endpoint=${BEE_SWAP_ENDPOINT} \
        --cors-allowed-origins="*"
        # --cors-allowed-origins="http://localhost:8080,http://localhost:9090"
}

run_bee() {
    echo "Starting bee node using config at: ${BEE_CONFIG}"
    ${bee_path}/bee start --config=${BEE_CONFIG}
}

usage() {
    echo "usage: ${0} [option] NODE_ID"
    echo 'options:'
    echo '	-dfs    runs a dfs node (depends on a bee node)'
    echo '	-test   runs a bee node for local tests'
    echo '	-fly	runs a bee node with given config'
    echo
}

option="${1}" 
case ${option} in
   -dfs) run_dfs;;
   -test) run_bee_test ${@:2};;
   -fly) run_bee;;
   *) usage; exit 1;;
esac
