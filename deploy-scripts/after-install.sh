#!/bin/bash
LOGFILE="/var/log/after-install.log"

exec 2>>$LOGFILE

echo "Making students binary executable..." >> $LOGFILE
sudo chmod +x /var/www/backend/bin/students

echo "Granting read permissions to the config directory..." >> $LOGFILE
sudo chmod -R u+r /var/www/backend/config

echo "After-install steps completed successfully." >> $LOGFILE

exit 0