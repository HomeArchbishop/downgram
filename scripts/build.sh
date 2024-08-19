DEV_MODE=false
for arg in "$@"
do
  if [ "$arg" == "--dev" ]; then
    DEV_MODE=true
    echo "[BUILD] Development mode enabled."
    break
  fi
done

if [ "$(uname)" == "Darwin" ]; then
  SUFFIX=""
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
  SUFFIX=""
elif [ "$(expr substr $(uname -s) 1 10)" == "MINGW32_NT" ] || [ "$(expr substr $(uname -s) 1 10)" == "MINGW64_NT" ]; then
  SUFFIX=".exe"
fi

if [ -d "dist" ]; then
  rm -rf ./dist
fi

mkdir -p dist
chmod 755 ./dist

version=$(cat ./VERSION)
echo "[BUILD] Kaption version: $version"

echo "[BUILD] GO building..."
if [ "$DEV_MODE" = true ]; then
  go build -o ./dist/ -ldflags "-X main.Version=$version" ./cmd/downgram
else
  go build -o ./dist/ -ldflags "-X main.Version=$version -H windowsgui" ./cmd/downgram
  # go build -o ./dist/ -ldflags "-X main.Version=$version" ./cmd/downgram
fi

cp ./VERSION ./dist/
cp ./LICENSE ./dist/

if [ "$DEV_MODE" = true ]; then
  echo "[BUILD] Development mode enabled."
  cat CACHE > ./dist/CACHE
  cat session.json > ./dist/session.json
fi
