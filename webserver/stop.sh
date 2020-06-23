sudo kill -9 $(ps -ef|grep wentmin |gawk '$0 !~/grep/ {print $2}' |tr -s '\n' ' ')
