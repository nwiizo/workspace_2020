#!/bin/sh

HTTPD_COUNT=`ps auxf | grep [h]ttpd | wc -l`
if [ $HTTPD_COUNT -gt 550 ];then
        date
        echo $HTTPD_COUNT process restart
        systemctl status httpd
else
        date
        echo $HTTPD_COUNT process not restart
fi