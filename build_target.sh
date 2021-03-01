TARGETS_DIR=targets

if [ "$#" -ne 1 ]; then
  echo "Usage $0 <target>"
  TARGETS=( $(ls -d ${TARGETS_DIR}/*) )
  echo "Available targets: ${TARGETS[@]} " 
else
  echo "Building target ${1}"
  docker build -t mastro-${1}:latest -f ${TARGETS_DIR}/${1}/Dockerfile .
fi
