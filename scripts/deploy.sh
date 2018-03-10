source scripts/version.sh

kubectl apply -f scripts/kubernetes/geoipfix-config.yml

DEPLOYMENT=$(envsubst < scripts/kubernetes/geoipfix-deployment.yml)
cat <<EOF | kubectl apply -f -
${DEPLOYMENT}
EOF

kubectl apply -f scripts/kubernetes/geoipfix-service.yml
