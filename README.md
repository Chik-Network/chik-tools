# Chik Tools

Collection of CLI tools for working with Chik Blockchain

## Installation

Download the correct executable file from the release page and run. If you are on debian/ubuntu, you can install using the apt repo, documented below.

### Homebrew Installation (MacOS)

1. `brew install chik-network/chik/chik-tools`

### Apt Repo Installation (Ubuntu/Debian)

#### Set up the repository

1. Update the `apt` package index and install packages to allow apt to use a repository over HTTPS:

```shell
sudo apt-get update

sudo apt-get install ca-certificates curl gnupg
```

2. Add Chik's official GPG Key:

```shell
curl -sL https://repo.chiknetwork.com/FD39E6D3.pubkey.asc | sudo gpg --dearmor -o /usr/share/keyrings/chik.gpg
```

3. Use the following command to set up the stable repository.

```shell 
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/chik.gpg] https://repo.chiknetwork.com/chik-tools/debian/ stable main" | sudo tee /etc/apt/sources.list.d/chik-tools.list > /dev/null
```

#### Install Chik Tools

1. Update the apt package index and install the latest version of Chik Tools

```shell
sudo apt-get update

sudo apt-get install chik-tools
```

### Go Install

`go install github.com/chik-network/chik-tools@latest`

### Docker

Prebuilt docker images are available.

Latest Release: `docker pull ghcr.io/chik-network/chik-tools:latest`

Latest Main Branch: `docker pull ghcr.io/chik-network/chik-tools:main`

Specific Tag: `docker pull ghcr.io/chik-network/chik-tools:0.1.0`

Note that you can choose a partial tag as well if you want the latest of a particular series:

Latest 0.1.z version: `docker pull ghcr.io/chik-network/chik-tools:0.1`

Latest 0.y.z version: `docker pull ghcr.io/chik-network/chik-tools:0`
