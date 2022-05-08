#!/bin/bash -e

set -e

# won't test:
# docs/en/start/configuration/start-apollo.md
# docs/zh/start/configuration/start-apollo.md
# because the github workflow can not connect to the apollo server due to the great firewall
quickstarts="docs/en/start/configuration/start.md
  docs/zh/start/configuration/start.md
  docs/en/start/state/start.md
  docs/zh/start/state/start.md
  docs/en/start/pubsub/start.md
  docs/zh/start/pubsub/start.md
  docs/en/start/lock/start.md
  docs/zh/start/lock/start.md
  docs/en/start/sequencer/start.md
  docs/zh/start/sequencer/start.md
  docs/en/start/rpc/helloworld.md
  docs/zh/start/rpc/helloworld.md
  docs/zh/start/file/minio.md
  docs/en/start/api_plugin/helloworld.md
  docs/zh/start/api_plugin/helloworld.md
  docs/zh/start/actuator/start.md
  docs/en/start/actuator/start.md
  docs/zh/start/trace/trace.md
  docs/en/start/trace/trace.md
  docs/zh/start/trace/skywalking.md
  docs/en/start/wasm/start.md
  docs/zh/start/wasm/start.md
"
quickstarts="docs/en/start/wasm/start.md
"

# download mdx
if ! test -e $(pwd)/etc/script/mdx
then
  curl -o $(pwd)/etc/script/mdx https://raw.githubusercontent.com/seeflood/mdx/main/mdx
fi
chmod +x $(pwd)/etc/script/mdx


# release all resources
release_resource() {
  if killall layotto; then
    echo "layotto released"
  fi
  if killall etcd; then
    echo "etcd released"
  fi
  if killall go; then
    echo "golang processes released"
  fi

  # remove all the docker containers
  if [ $(docker ps -a -q | wc -l) -gt 0 ]; then
    docker rm -f $(docker ps -a -q)
  fi
}

release_resource

# download etcd
sh etc/script/download_etcd.sh

# test quickstarts
for doc in ${quickstarts}; do
  echo "Start testing $doc......"

  #./mdx docs/en/start/state/start.md
  $(pwd)/etc/script/mdx $doc

  echo "End testing $doc......"
  release_resource
done
