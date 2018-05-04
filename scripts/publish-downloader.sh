docker build -f Dockerfile.downloader -t gcr.io/${GCLOUD_PROJECT_ID}/geoipfix-downloader:0.1.2 .
gcloud docker -- push gcr.io/${GCLOUD_PROJECT_ID}/geoipfix-downloader:0.1.2
