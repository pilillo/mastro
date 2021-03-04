# Builder
PROJECT=mastro
TARGETS_DIR=targets
TARGETS=( $(ls -d ${TARGETS_DIR}/*) ) 

while [[ $# -gt 0 ]]; do
  key="${1}"
  case $key in
    -t|--target)
    TARGET=${2}
    shift
    ;;
    -o|--organization)
    ORGANIZATION=${2}
    shift
    ;;
    -p|--push)
    PUSH=true
    ;;
    -s|--static)
    STATIC=true
    ;;
  esac
  # shift just processed arg
  shift
done

if [ -z ${ORGANIZATION} ] || [ -z ${TARGET} ]; then
  echo "Usage $0 -t|--target <target> -o|--organization <organization> {-p|--push}"
  echo "Available targets: ${TARGETS[@]} "
  exit
fi

# in static builds - include kerberos support (for crawlers only)
# can be overridden
#GO_BUILD_TAGS="${GO_BUILD_TAGS:--tags=kerberos}"
GO_BUILD_TAGS="${GO_BUILD_TAGS:-}"

go_static_build() {
  # go module always on since we use it
  export GO111MODULE=on
  # static build by disabling CGO
  export CGO_ENABLED=0
  # info for cross compilation
  export GOOS=linux
  export GOARCH=amd64
  echo "Running - go build ${GO_BUILD_TAGS} -o ${ARTIFACT} ${LOCATION}"
  go build ${GO_BUILD_TAGS} -o ${ARTIFACT} ${LOCATION}
  echo "Compiled ${ARTIFACT} from ${LOCATION}"
}

# decide whether to push a latest and latest-static tag (can be disabled from outside)
PUSH_LATEST="${PUSH_LATEST:-true}"

dhub_push() {
  docker push ${IMAGE}:${BUILD_TAG}
  echo "pushed ${IMAGE}:${BUILD_TAG}"
  if [ ${PUSH_LATEST} ]; then
    LATEST="latest"$([[ $test == *"static" ]] && echo "-static")
    docker image tag ${IMAGE}:${BUILD_TAG} ${IMAGE}:${LATEST}
    docker push ${IMAGE}:${LATEST}
  fi
}

# default build tag
BUILD_TAG=$(date +%Y%m%d)

if [ "${TARGET}" == "all" ]; then
  echo "Building all-in-one image"
  ARTIFACT=${PROJECT}
  LOCATION="."
  IMAGE=${ORGANIZATION}/${PROJECT}
  if [ ${STATIC} ]; then
    go_static_build
    # move artifact to a fresh docker image
    BUILD_TAG=${BUILD_TAG}-static
    docker build --build-arg ARTIFACT -t ${IMAGE}:${BUILD_TAG} -f Dockerfile.static .
  else
    docker build -t ${IMAGE}:${BUILD_TAG} -f Dockerfile .
  fi
  
  if [ ${PUSH} ]; then
    echo "Pushing to dockerhub"
    dhub_push
  fi
elif [ -d "${TARGETS_DIR}/${TARGET}" ]; then
  echo "Building target ${TARGET}"
  ARTIFACT=${PROJECT}
  LOCATION="./${TARGETS_DIR}/${TARGET}"
  IMAGE=${ORGANIZATION}/${PROJECT}-${TARGET}

  if [ ${STATIC} ]; then
    go_static_build
    # move artifact to a fresh docker image
    BUILD_TAG=${BUILD_TAG}-static
    docker build --build-arg ARTIFACT -t ${IMAGE}:${BUILD_TAG} -f Dockerfile.static .
  else
    docker build -t ${IMAGE}:${BUILD_TAG} -f ${TARGETS_DIR}/${TARGET}/Dockerfile .
  fi

  if [ ${PUSH} ]; then
    echo "Pushing to dockerhub"
    dhub_push
  fi
else
   echo "Selected target path ${TARGET} does not exist"
   echo "Available targets: ${TARGETS[@]} "
fi