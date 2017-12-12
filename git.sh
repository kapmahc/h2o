#!/bin/sh

root=ant-design-pro-master

# add plugins
for p in nut forum reading survey erp pos mall ops/vpn ops/mail
do
  git add ${root}/src/routes/${p}
done

# add source files
for f in common/nav.js layouts/Application.js layouts/Dashboard.js
do
  git add ${root}/src/${f}
done

git status
