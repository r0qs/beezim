#!/bin/sh

# TODO: check if docker is present

ZXS_HOME=/home/emscripten/zxs
assets_dir=$(pwd)/indexer/assets/js/xapian
# TODO: check if exists before re-generate

# TODO: pass option to re-generate assets
mkdir -p ${assets_dir}

# Download zxs image with xapian compiled with emscripten compiler, and zxs installed.
docker pull ghcr.io/r0qs/zxs:latest

# Generate searcher and indexer web assembly library
docker run -it --rm --name zxs-searcher-0 \
	--user $(id -u):$(id -g) \
	--mount "type=bind,src=${assets_dir},dst=${ZXS_HOME}/dist" \
	zxs --entrypoint sh -c "npm run build && npm run preparelib"

# Build beezim image
docker build -t beezim -f Dockerfile .

# Run beezim connection to a localhost bee node
docker run -it --rm --name beezim-cli-0 \
	--user $(id -u):$(id -g) \
	--mount "type=bind,src=$(pwd)/.env,dst=/src/.env,readonly" \
	--mount "type=bind,src=$(pwd)/datadir,dst=/src/datadir" \
	--mount "type=bind,src=$(pwd)/indexer/assets/js/xapian,dst=/src/xapian" \
	-p 1733:1733 \
	-p 1735:1735 \
	--network="host" \
	beezim ./bin/beezim-cli mirror --zim=wikipedia_es_climate_change_mini_2022-02.zim --bee-api-url=http://localhost:1733 --bee-debug-api-url=http://localhost:1735 --batch-id=388b9a93fc084d350b2320bedacb3a88779867d956b20a2716512138bc88eac0 --enable-search