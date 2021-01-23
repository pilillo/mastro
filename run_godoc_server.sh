PORT=6060
MODNAME=$(go list -e)

echo "Check Project documentation at http://localhost:${PORT}/pkg/${MODNAME}"
godoc -http=:${PORT}
