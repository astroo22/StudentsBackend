#!/bin/bash

LOGFILE="/var/log/start-app.log"

exec 2>>$LOGFILE

echo "Starting the students service..." >> $LOGFILE
systemctl start students.service

if [ $? -eq 0 ]; then
    echo "Successfully started the students service." >> $LOGFILE
    exit 0
else
    echo "Failed to start the students service." >> $LOGFILE
    exit 1
fi
