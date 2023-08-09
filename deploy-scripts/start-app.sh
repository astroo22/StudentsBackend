#!/bin/bash
LOGFILE="/var/log/start-app.log"

exec 2>>$LOGFILE

export APP_ENV=prod

# If the 'students' process is already running, kill it
if pgrep students > /dev/null; then
    echo "Students service is currently running. Stopping..." >> $LOGFILE
    pkill students
    sleep 2 # give second to allow it to shutdown
fi

echo "Starting the students service..." >> $LOGFILE
/var/www/backend/bin/students &
if [ $? -eq 0 ]; then
    echo "Successfully started the students service." >> $LOGFILE
    exit 0
else
    echo "Failed to start the students service." >> $LOGFILE
    exit 1
fi