[![Go build binary and upload release](https://github.com/berejant/kneu-github-webhooks/actions/workflows/release.yaml/badge.svg)](https://github.com/berejant/kneu-github-webhooks/actions/workflows/release.yaml)

# Install
```bash

# Download binary
cd {CI_SCRIPT_DIR}
wget -c https://github.com/berejant/kneu-github-webhooks/releases/latest/download/kneu-github-webhooks-linux-amd64.tar.gz -O - | tar -xz

# Set configuration
wget https://raw.githubusercontent.com/berejant/kneu-github-webhooks/main/.env.example -O .env
nano .env

# configure supervisor job
wget https://raw.githubusercontent.com/berejant/kneu-github-webhooks/main/supervisor-job-example.conf
nano supervisor-job-example.conf
mv supervisor-job-example.conf YOUR_JOB_NAME.conf
sudo ln -s $(pwd)/YOUR_JOB_NAME.conf /etc/supervisor/conf.d/
sudo service supervisor restart

```
