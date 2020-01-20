#!/usr/bin/env bash

# Exit with error if any command returns non zero code.
set -e
# Exit with error if any undefined variable is referenced.
set -u

REGISTRY=${REGISTRY:-registry.ori.co}

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )/.." && pwd )
PULL=

USAGE="Usage: $0 [options] <version>

Options:
    -h          Display this help message.
    -p          Pull base image before building.

Positional arguments:
    <version>   Version of software (probably branch/tag name or commit sha).
                This is used as the image tag.
"

log() { echo "$@" >&2; }

while getopts "hp" opt; do
    case "$opt" in
        h)
            echo "$USAGE"
            exit
            ;;
        p)
            PULL=1
            ;;
        \?)
            log "$USAGE"
            log
            log "ERROR: Invalid option: -$OPTARG"
            exit 1
    esac
done
shift $((OPTIND-1))

if [ "$#" -eq 0 ]; then
    log "$USAGE"
    log
    log "ERROR: No version/tag specified."
    exit 1
fi
if [ "$#" -gt 1 ]; then
    log "$USAGE"
    log
    log "ERROR: Multiple version/tag specified."
    exit 1
fi
TAG=$1

if [ -z "$PULL" ]; then
    BUILD_ARGS=""
else
    BUILD_ARGS="--pull"
fi

CLI_IMG="$REGISTRY/oriindustries/ori-app/client:$TAG"
SRV_IMG="$REGISTRY/oriindustries/ori-app/client:$TAG"

log "Will build images"
log "cli:               $CLI_IMG"
log "srv:               $SRV_IMG"
log
log

log "Building client image"
log
set -x
docker build -t $CLI_IMG \
    --build-arg=ORIAPP_BUILD_VERSION=$TAG \
    --build-arg=ORIAPP_BUILD_SHA=${CI_COMMIT_SHA:-$(git rev-parse HEAD)} \
    --build-arg=ORIAPP_BUILD_DATE="$(date -u '+%Y-%M-%d %H:%m:%S')" \
    -f $DIR/docker/Dockerfile.client $DIR
set +x
log
log

log "Building server image"
log
set -x
docker build -t $SRV_IMG \
    --build-arg=ORIAPP_BUILD_VERSION=$TAG \
    --build-arg=ORIAPP_BUILD_SHA=${CI_COMMIT_SHA:-$(git rev-parse HEAD)} \
    --build-arg=ORIAPP_BUILD_DATE="$(date -u '+%Y-%M-%d %H:%m:%S')" \
    -f $DIR/docker/Dockerfile.srv $DIR
set +x
log
log



log "Built images"
log
log "cli:               $CLI_IMG"
log "srv:               $SRV_IMG"
