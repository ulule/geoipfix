source scripts/version.sh

kubectl create -f scripts/kubernetes/geoipfix-config.yml

DEPLOYMENT=$(envsubst < scripts/kubernetes/geoipfix-deployment.yml)
cat <<EOF | kubectl create -f -
${DEPLOYMENT}
EOF

kubectl create -f scripts/kubernetes/geoipfix-service.yml
