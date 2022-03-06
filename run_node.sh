#!/bin/sh
# Tested with Bee version: 1.4.3-9546fedb

[[ -e "$(pwd)/.env" ]] && source $(pwd)/.env

if [[ -z "${BEE_BIN_PATH}" ]]; then
    bee_path="/usr/bin"
else
    bee_path=${BEE_BIN_PATH}
fi
data_root=$(pwd)/data_root
mkdir -p ${data_root}

run_dfs() {
    dfs server --dataDir=${data_root}/dfs --beePort=1633
}

run_bee() {
    echo "Using swap endpoint: ${BEE_SWAP_ENDPOINT}"
    [[ "$#" -ne 1 ]] && { echo "Usage : ./run_node -r NODE_ID"; exit 1; }
    node=$1
    port=$(($node+15))
    api_port=${port}33
    p2p_port=${port}34
    debug_port=${port}35
    data_dir="${data_root}/bee_$node"

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

usage() {
    echo "usage: ${0} [option] NODE_ID"
    echo 'options:'
    echo '	-f	runs a dfs node (depends on a bee node)'
    echo '	-r	runs a bee node'
    echo
}

option="${1}" 
case ${option} in
   -f) run_dfs;;
   -r) run_bee ${@:2};;
   *) usage; exit 1;;
esac
