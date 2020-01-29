#!/bin/sh
# This is a comment!

if [ "$1" = "staging" ]
then
  echo Building wallet with staging Enviorment

  NODE_ENV=staging VUE_APP_CHAIN=rns-test-01 VUE_APP_CLAIM_URL=https://proxy.testnet.color-platform.rnssol.com:8000/claim yarn build

elif [ "$1" = "production" ]
then
  echo Building wallet with production Enviorment

  NODE_ENV=production VUE_APP_CHAIN=color-mainbeta-2 VUE_APP_CLAIM_URL=https://getcoins.mainbeta-2.color-platform.org/claim yarn build
else 
  echo Environment not provided, e.g staging
fi

#npm install
#STARGATE=https://proxy.testnet.color-platform.rnssol.com:9071 RPC=https://rpc.testnet.color-platform.rnssol.com yarn build
#yarn serve:dist


