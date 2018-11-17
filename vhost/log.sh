#/bin/bash
LOGPATH = /sakuraus/vhost/access.log
RESULT = /sakuraus/backup_log/$(date +%Y%m%d)/access.log
mkdir $RESULT
mv $LOGPATH $RESULT

