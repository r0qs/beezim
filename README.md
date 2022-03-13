# BeeZIM Mirror

This repository provides tools to publish copies of entire websites on Swarm.

Thankfully to [OpenZIM Project](https://wiki.openzim.org/wiki/OpenZIM), which is currently maintained by a not-for-profit entity named [Kiwix](https://www.kiwix.org/en/), websites can be highly compressed into a single file and easily shared by users or viewed in devices with very limited computational resources. 

The compressed files follow the [ZIM file format](https://wiki.openzim.org/wiki/ZIM_file_format), and according to Kiwix, the entire Wikipedia can be compressed in an 80Gb zim file, containing more than 6 million articles, with images!
The main purpose of Zim file format is to store websites content for offline usage.

Our goal is to use the zim files to mirror websites on decentralized storage like [Swarm](https://www.ethswarm.org/), providing censorship resistance to the websites that we love so much, like [Wikipedia](https://www.wikipedia.org/) and [Project Gutenberg](https://www.gutenberg.org/). We want to make a version of the [Internet Archive](https://archive.org/web/) that cannot be censored.


## How it works

// TODO: explain how it works

Files are downloaded from the zim, parsed and extracted as tar archives
The parser embeds metadata and a search engine to the archives, and uses swarm as a kv store for indexed data on each mirrored website
Each metadata file is also uploaded to swarm and queried from the js when searching for articles or words in the main page
The archives are then uploaded to swarm and its reference (i.e., the manifest address) and the metadata are also stored locally, but not the zims.
They are deleted after completion of the upload.

Make a ENS entry to the root manifest and point to: https://beezim.bzz.link or https://bzzim.bzz.link

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

// TODO: add description, explain sub commands and add example for each command
### Configure the Bee environment
Bee needs to be installed before using Beezim. See [.env-example](.env-example) for an example of the neccessary configuration parameters.
Create a file named **.env** with configuration parameters for your system.

### Download ZIM files

```
go run cli/main.go download --kiwix=wikipedia --zim=wikipedia_es_climate_change_mini_2022-02.zim 
```

### Parse ZIM files and embed a search engine

This converts the zim files to tar archives and embed information to them (JS, CSS, HTML) and a search engine using the Xapian index.

```
go run cli/main.go parse --zim=wikipedia_es_climate_change_mini_2022-02.zim --enable-search
```

### Parse ZIM files without embedded search

This converts the zim files to tar archives and embed information to them (JS, CSS, HTML).

```
go run cli/main.go parse --zim=wikipedia_es_climate_change_mini_2022-02.zim
```

### Upload the TAR to Swarm

```
go run cli/main.go upload --tar=wikipedia_es_climate_change_mini_2022-02.tar --gateway
```
```
go run cli/main.go upload --tar=wikipedia_es_climate_change_mini_2022-02.tar --batch-id=8e747b4aefe21a9c902337058f7aad71aa3170a9f399ece6f0bdb9f1ec432685
```

```
go run cli/main.go upload all --batch-id=8e747b4aefe21a9c902337058f7aad71aa3170a9f399ece6f0bdb9f1ec432685
```

```
go run cli/main.go upload --kiwix="gutenberg" all --batch-id=8e747b4aefe21a9c902337058f7aad71aa3170a9f399ece6f0bdb9f1ec432685
```

### Mirror

Performs the download `->` parser `->` upload for one or many zims. (not finished yet)

```
go run cli/main.go mirror --kiwix=gutenberg --zim=gutenberg_af_all_2022-03.zim --batch-id=8e747b4aefe21a9c902337058f7aad71aa3170a9f399ece6f0bdb9f1ec432685
```

```
go run cli/main.go mirror --kiwix=others --zim=alpinelinux_en_all_nopic_2021-03.zim --bee-api-url=http://localhost:1733 --bee-debug-api-url=http://localhost:1735 --batch-id=388b9a93fc084d350b2320bedacb3a88779867d956b20a2716512138bc88eac0
```
