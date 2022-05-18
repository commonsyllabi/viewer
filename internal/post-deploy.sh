#!/usr/bin/env bash

systemctl --user restart cosyl.service
currentDate=`date`
echo Successful deployed at $currentDate >> deploy.log