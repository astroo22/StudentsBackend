#!/bin/bash

LOGFILE="/var/log/stop-app.log"

exec 2>>$LOGFILE

if pgrep students > /dev/null; then
    echo "Stopping the students service..." >> $LOGFILE
    pkill students
    if [ $? -eq 0 ]; then
        echo "Successfully stopped the students service." >> $LOGFILE
    else
        echo "Failed to stop the students service." >> $LOGFILE
    fi
else
    echo "Students service is not running. Nothing to stop." >> $LOGFILE
fi

exit 0