@echo off

set host=118.195.244.48
set centerhost=82.157.137.166
set gamehost=82.156.167.77
set user=root

set remoteBase=/home/sheep/

set localPath=..\bin\bin\
set remotePath=%remoteBase%bin/

set conflocalPath=..\conf\normalConf\
set confPath=%remoteBase%conf/

set excelLocalPath=..\table\
set excelRemotePath=%remoteBase%table/

set iplocationLocalPath=..\IP2LOCATION-LITE-DB3.IPV6.BIN\
set iplocationRemotePath=%remoteBase%IP2LOCATION-LITE-DB3.IPV6.BIN/

set sslLocalPath=..\easygame2021.com_nginx\
set sslRemotePath=%remoteBase%easygame2021.com_nginx/

set keyLocalPath=C:\Users\Dell\.ssh\id_rsa.pub
set keyRemotePath=/root/.ssh/

if %1==1 (goto program) 
if %1==2 (goto conf) 
if %1==3 (goto excel) 
if %1==4 (goto iplocation) 
if %1==5 (goto ssl) 
if %1==6 (goto keystore)



:program
for /r %localPath% %%i in (*) do (scp -r %%i %user%@%centerhost%:%remotePath%)
for /r %localPath% %%i in (*) do (scp -r %%i %user%@%gamehost%:%remotePath%)
goto end

:conf
for /r %conflocalPath% %%i in (*) do (scp -r %%i %user%@%centerhost%:%confPath%)
for /r %conflocalPath% %%i in (*) do (scp -r %%i %user%@%gamehost%:%confPath%)
goto end

:excel
for /r %excelLocalPath% %%i in (*) do (scp -r %%i %user%@%centerhost%:%excelRemotePath%)
for /r %excelLocalPath% %%i in (*) do (scp -r %%i %user%@%gamehost%:%excelRemotePath%)
goto end

:iplocation
for /r %iplocationLocalPath% %%i in (*) do (scp -r %%i %user%@%centerhost%:%iplocationRemotePath%)
for /r %iplocationLocalPath% %%i in (*) do (scp -r %%i %user%@%gamehost%:%iplocationRemotePath%)
goto end

:ssl
for /r %sslLocalPath% %%i in (*) do (scp -r %%i %user%@%centerhost%:%sslRemotePath%)
for /r %sslLocalPath% %%i in (*) do (scp -r %%i %user%@%gamehost%:%sslRemotePath%)

:keystore
scp -r  %keyLocalPath%  %user%@%host%:%keyRemotePath%
scp -r ./kill.sh %user%@%centerhost%:%remoteBase%
scp -r ./run.sh %user%@%centerhost%:%remoteBase%

:end
exit