#!/bin/bash

rm -r ~/.laconic2d || true
LACONIC2D_BIN=$(which laconic2d)
# configure laconic2d
$LACONIC2D_BIN config set config log_level "*:error,p2p:info,state:info,auction:info,bond:info,registry:info" --skip-validate
$LACONIC2D_BIN config set client chain-id demo
$LACONIC2D_BIN config set client keyring-backend test
$LACONIC2D_BIN keys add alice
$LACONIC2D_BIN keys add bob
$LACONIC2D_BIN init test --chain-id demo --default-denom photon
# update genesis
$LACONIC2D_BIN genesis add-genesis-account alice 10000000photon --keyring-backend test
$LACONIC2D_BIN genesis add-genesis-account bob 1000photon --keyring-backend test
# create default validator
$LACONIC2D_BIN genesis gentx alice 1000000photon --chain-id demo
$LACONIC2D_BIN genesis collect-gentxs
