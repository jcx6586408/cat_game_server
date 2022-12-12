@echo off

set host=118.195.244.48
set user=root

set remoteBase=/home/sheep/

set localPath=.\bin\bin\
set remotePath=%remoteBase%bin/

set conflocalPath=.\conf\remoteConf\
set confPath=%remoteBase%conf/

set excelLocalPath=.\table\
set excelRemotePath=%remoteBase%table/

set iplocationLocalPath=.\IP2LOCATION-LITE-DB3.IPV6.BIN\
set iplocationRemotePath=%remoteBase%IP2LOCATION-LITE-DB3.IPV6.BIN/

set sslLocalPath=.\ssl\Nginx\
set sslRemotePath=%remoteBase%ssl/Nginx/

if %1==1 (goto program) 
if %1==2 (goto conf) 
if %1==3 (goto excel) 
if %1==4 (goto iplocation) 
if %1==5 (goto ssl) 



:program
for /r %localPath% %%i in (*) do (scp -r %%i %user%@%host%:%remotePath%)
goto end

:conf
for /r %conflocalPath% %%i in (*) do (scp -r %%i %user%@%host%:%confPath%)
goto end

:excel
for /r %excelLocalPath% %%i in (*) do (scp -r %%i %user%@%host%:%excelRemotePath%)
goto end

:iplocation
for /r %iplocationLocalPath% %%i in (*) do (scp -r %%i %user%@%host%:%iplocationRemotePath%)
goto end

:ssl
for /r %sslLocalPath% %%i in (*) do (scp -r %%i %user%@%host%:%sslRemotePath%)

:end
exit