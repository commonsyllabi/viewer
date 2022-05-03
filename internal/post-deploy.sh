#!/usr/bin/env bash

systemctl --user restart cosyl.service
currentDate=`date`
echo Successful eployed at $currentDate >> deploy.log