TARGETS_DIR=targets

TARGETS=( $(ls -d ${TARGETS_DIR}/*) ) 

if [ "$#" -ne 1 ]; then
  echo "Usage $0 <target>"
  echo "Available targets: ${TARGETS[@]} "
  exit  
fi

echo "Building target ${1}"
docker build -t mastro-${1}:latest -f ${TARGETS_DIR}/${1}/Dockerfile .