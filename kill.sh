kill -9 $( ps -e|grep leaf |awk '{print $1}')
kill -9 $( ps -e|grep rank |awk '{print $1}')