# BeeZIM Mirror Tool

This repository provides tools to publish copies of entire websites on Swarm.

Thankfully to [OpenZIM Project](https://wiki.openzim.org/wiki/OpenZIM), which is currently maintained by a not-for-profit entity named [Kiwix](https://www.kiwix.org/en/), websites can be highly compressed into a single file and easily shared by users or viewed in devices with very limited computational resources. 

The compressed files follow the [ZIM file format](https://wiki.openzim.org/wiki/ZIM_file_format), and according to Kiwix the entire Wikipedia can be compressed in an 80Gb zim file, containing more than 6 million articles, with images!
The main purpose of Zim file format is to store websites content for offline usage.

Our goal is to use the zim files to mirror websites on decentralized storages like [Swarm](https://www.ethswarm.org/), providing censorship resistance to the websites that we love so much, like [Wikipedia](https://www.wikipedia.org/) and [Project Gutenberg](https://www.gutenberg.org/). We want to make a version of the [Internet Archive](https://archive.org/web/) that cannot be censored.

## How it works

Beezim is a command-line tool that uses the [Bee](https://github.com/ethersphere/bee) API to upload content on Swarm.
Files in the ZIM format can be downloaded from a web2 mirror website or provided by the user if it is already stored locally.
The ZIMs are parsed, and a tar archive is generated from them.

The parser can optionally embed metadata, a text search engine and a search DApp to the archives.
Each tar archive is then uploaded to swarm, and its reference (i.e., the [manifest](https://docs.ethswarm.org/docs/access-the-swarm/upload-a-directory#upload-the-directory-containing-your-website) address) is returned as output. Please keep this stored so you can access your page later on. We plan to provide a key-value store to manage the metadata and references in the future.

The search engine can be enabled during the parsing by using the option `--enable-search`.
And allow users to query for uploaded articles containing data from the text or title of the ZIMs.
Beezim also embeds a navigation bar and pages to display information about the uploaded files and the project when the search is enabled.

The ZIM and/or tar files can be automatically deleted from the host machine after upload, using the option `--clean`.

The default behavior of Beezim is to `mirror` ZIMs to Swarm as they are, **without** append metadata or the search tool to it.
However, if you would like to be able to search on the uploaded content in a similar fashion provided by [Kiwix](https://library.kiwix.org/), but without relying on server-side services or database, you can try out our search tool!

Our search tool is called Zim Xapian Searchindex, or [zxs](https://github.com/r0qs/zxs) for short.
It is a WebAssembly library and javascript search tool that can read the indexes in the [Xapian](https://xapian.org/) format extracted from ZIM files under `X/fulltext/xapian` and `X/title/xapian`.

ZXS enables the search of indexed data in your browser using the Xapian database that is [already embedded in the ZIM files]((https://wiki.openzim.org/wiki/Search_indexes)) without interacting with a server.
It is based on Xapian search engine library and compiled for WebAssembly using the [Emscripten](https://emscripten.org/) compiler.
By using *Beezim* search tool, no server can monitor what you are searching for! Everything happens on your browser.
ZXS could also allow users to search contents without an internet connection,
embedding the javascript search tool and the WebAssembly engine directly in the ZIM files.
Although, this is not done yet!

## How to run

The current command-line tool has the following available commands:
```
Swarm zim mirror command-line tool

Usage:
  beezim [command]

Available Commands:
  download    download zim file
  help        Help about any command
  list        Shows a list of compressed websites currently maintained by Kiwix
  mirror      mirror kiwix zim repositories to swarm
  parse       parse zim file
  upload      upload zim file to swarm

Flags:
      --batch-amount int           bee postage batch amount (default 100000000)
      --batch-depth uint           bee postage batch depth (default 30)
      --batch-id string            bee postage batch ID
      --bee-api-url string         bee api url (default "http://localhost:1633")
      --bee-debug-api-url string   bee debug api url (default "http://localhost:1635")
      --clean                      delete all downloaded zim and generated tar files
      --datadir string             path to datadir directory (default "./datadir")
      --enable-search              enable search index
      --gas-price string           gas price for postage stamps purchase
      --gateway                    connect to the swarm public gateway (default "https://gateway.ethswarm.org")
  -h, --help                       help for beezim
      --kiwix string               name of the compressed website hosted by Kiwix. Run "list" to see all available options (default "wikipedia")

Use "beezim [command] --help" for more information about a command.
```


![main page](/images/main-page.png)

| Articles | List uploaded files | Files information | Search bar | Search Results |
| -------- | ------------------- | ----------------- | ---------- | -------------- |
| ![Articles](/images/article.png) | ![List Files](/images/files.png) | ![File Info](/images/files-open.png) | ![Search Bar](/images/search-bar.png) | ![Search Results](/images/search-results.png) |


### Configure the Bee environment

Beezim uploads files to Swarm by connecting to a bee node.
But you can use the `--gateway` option to upload directly to the [public swarm gateway](https://gateway.ethswarm.org/).
However the public gateway has a maximum upload limit of 10 MB per file.

Example using the gateway:
```
beezim mirror --zim=wikipedia_cr_all_maxi_2022-02.zim --gateway
```

For best experience and convenience it is recommended that you run your own bee node before try Beezim.
See [.env-example](.env-example) for an example of the necessary configuration parameters.
Create a file named **.env** with configuration parameters for your system.

### Download ZIM files

```
beezim download \
  --kiwix=wikipedia \
  --zim=wikipedia_es_climate_change_mini_2022-02.zim 
```

### Parse ZIM files and embed a search engine

This converts the zim files to tar archives and embed information to them (JS, CSS, HTML) and a search engine using the Xapian index.

```
beezim parse \
  --zim=wikipedia_es_climate_change_mini_2022-02.zim \
  --enable-search
```

### Parse ZIM files without embedded search

This converts the zim files to tar archives and embed information to them (JS, CSS, HTML).

```
beezim parse --zim=wikipedia_es_climate_change_mini_2022-02.zim
```

### Upload the TAR to Swarm

```
beezim upload --tar=wikipedia_es_climate_change_mini_2022-02.tar --gateway
```

```
beezim upload \
  --tar=wikipedia_es_climate_change_mini_2022-02.tar \
  --batch-id=8e747b4aefe21a9c902337058f7aad71aa3170a9f399ece6f0bdb9f1ec432685
```

```
beezim upload all \
  --batch-id=8e747b4aefe21a9c902337058f7aad71aa3170a9f399ece6f0bdb9f1ec432685
```

```
beezim upload --kiwix="gutenberg" all \
  --batch-id=8e747b4aefe21a9c902337058f7aad71aa3170a9f399ece6f0bdb9f1ec432685
```

### Mirror

Performs the download `->` parser `->` upload for one or many zims. (not finished yet)

```
beezim mirror --kiwix=gutenberg \
  --zim=gutenberg_af_all_2022-03.zim \
  --batch-id=8e747b4aefe21a9c902337058f7aad71aa3170a9f399ece6f0bdb9f1ec432685
```

```
beezim mirror --kiwix=others \
  --zim=alpinelinux_en_all_nopic_2021-03.zim \
  --bee-api-url=http://localhost:1733 \
  --bee-debug-api-url=http://localhost:1735 \
  --batch-id=388b9a93fc084d350b2320bedacb3a88779867d956b20a2716512138bc88eac0
```

## Using Docker to Build BeeZIM

### Without search engine

If you don't plan to use the search engine and would like to mirror ZIMs as they are.
You can just install BeeZIM in your machine and use it without the `--enable-search` option,
or build the BeeZIM docker image (not the docker compose).

```
docker build -t beezim -f Dockerfile .
```

### With the Search Engine and Search DApp Tool

Our search DApp depends on [Zim Xapian Searchindex](https://github.com/r0qs/zxs), a WebAssembly library
and javascript search tool that can read the search indexes extracted from ZIM files.

We also provide a `docker-compose.yml` to download the ZXS image and build Beezim in your local machine.
You can use it running the command below:

```
docker-compose -f docker-compose.yml up --build --remove-orphans && docker-compose rm -fsv
```

```
docker-compose run --rm \
  --user $(id -u):$(id -g) \
  beezim ./bin/beezim-cli mirror \
  --zim=wikipedia_es_climate_change_mini_2022-02.zim \
  --bee-api-url=http://localhost:1633 \
  --bee-debug-api-url=http://localhost:1635 \
  --batch-id=388b9a93fc084d350b2320bedacb3a88779867d956b20a2716512138bc88eac0 \
  --enable-search
```

There is also a script to simplify a bit the above command when running BeeZIM with docker:
```
./dc-beezim-cli mirror \
  --zim=wikipedia_es_climate_change_mini_2022-02.zim \
  --batch-id=388b9a93fc084d350b2320bedacb3a88779867d956b20a2716512138bc88eac0 \
  --enable-search
```