source scripts/version.sh

make docker-build

docker build -t gcr.io/${GCLOUD_PROJECT_ID}/geoipfix:${GEOIPFIX_VERSION} .
gcloud docker -- push gcr.io/${GCLOUD_PROJECT_ID}/geoipfix:${GEOIPFIX_VERSION}
