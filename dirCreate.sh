mkdir conf
#如果文件夹不存在，则创建文件夹
tempPath="conf/"
if [ ! -d "$tempPath" ]; then
mkdir conf
fi

tempPath="bin/"
if [ ! -d "$tempPath" ]; then
mkdir bin
fi


tempPath="IP2LOCATION-LITE-DB3.IPV6.BIN/"
if [ ! -d "$tempPath" ]; then
mkdir IP2LOCATION-LITE-DB3.IPV6.BIN
fi

tempPath="table/"
if [ ! -d "$tempPath" ]; then
mkdir table
fi

tempPath="ssl/"
if [ ! -d "$tempPath" ]; then
mkdir ssl
fi

tempPath="ssh/Nginx/"
if [ ! -d "$tempPath" ]; then
mkdir ssh/Nginx
fi
