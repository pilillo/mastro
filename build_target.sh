PROJECT=mastro
TARGETS_DIR=targets
TARGETS=( $(ls -d ${TARGETS_DIR}/*) ) 

if [ "$#" -ne 3 ]; then
  echo "Usage $0 <target> <organization> <push>"
  echo "Available targets: ${TARGETS[@]} "
elif [ "${1}" == "all" ]; then
  echo "Building all-in-one image"
  BUILD_TAG=$(date +%Y%m%d)
  docker build -t ${2}/${PROJECT}:${BUILD_TAG} -f Dockerfile .
  [ "${3}" == "push" ] && docker push ${2}/${PROJECT}:${BUILD_TAG}
  docker image tag ${2}/${PROJECT}:${BUILD_TAG} ${2}/${PROJECT}:latest
  [ "${3}" == "push" ] && docker push ${2}/${PROJECT}:${BUILD_TAG}
elif [ -d "${TARGETS_DIR}/${1}" ]; then
  echo "Building target ${1}"
  BUILD_TAG=$(date +%Y%m%d)
  docker build -t ${2}/${PROJECT}-${1}:${BUILD_TAG} -f ${TARGETS_DIR}/${1}/Dockerfile .
  [ "${3}" == "push" ] && docker push ${2}/${PROJECT}:${BUILD_TAG}
  docker image tag ${2}/${PROJECT}-${1}:${BUILD_TAG} ${2}/${PROJECT}-${1}:latest
  [ "${3}" == "push" ] && docker push ${2}/${PROJECT}:latest
else
  echo "Selected target path ${1} does not exist"
  echo "Available targets: ${TARGETS[@]} "
fi