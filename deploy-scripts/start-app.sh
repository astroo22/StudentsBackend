#!/bin/bash

LOGFILE="/var/log/start-app.log"

exec 2>>$LOGFILE

echo "Starting the backend service..." >> $LOGFILE
systemctl start backend.service

if [ $? -eq 0 ]; then
    echo "Successfully started the backend service." >> $LOGFILE
    exit 0
else
    echo "Failed to start the backend service." >> $LOGFILE
    exit 1
fi
