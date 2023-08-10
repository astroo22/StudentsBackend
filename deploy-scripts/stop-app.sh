#!/bin/bash

LOGFILE="/var/log/stop-app.log"

exec 2>>$LOGFILE

echo "Stopping the students service..." >> $LOGFILE
systemctl stop backend.service

if [ $? -eq 0 ]; then
    echo "Successfully stopped the students service." >> $LOGFILE
    exit 0
else
    echo "Failed to stop the students service." >> $LOGFILE
    exit 1
fi