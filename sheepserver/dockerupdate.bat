@echo off

set host=118.195.244.48
set gamehost=82.156.167.77
set user=root

set remoteBase=/home/docker/sheepExam/

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
for /r %localPath% %%i in (*) do (scp -r %%i %user%@%gamehost%:%remotePath%)
goto end

:conf
for /r %conflocalPath% %%i in (*) do (scp -r %%i %user%@%gamehost%:%confPath%)
goto end

:excel
for /r %excelLocalPath% %%i in (*) do (scp -r %%i %user%@%gamehost%:%excelRemotePath%)
goto end

:iplocation
for /r %iplocationLocalPath% %%i in (*) do (scp -r %%i %user%@%gamehost%:%iplocationRemotePath%)
goto end

:ssl
for /r %sslLocalPath% %%i in (*) do (scp -r %%i %user%@%gamehost%:%sslRemotePath%)
goto end

:keystore
@REM 游戏服配置文件
scp -r ./kill.sh %user%@%gamehost%:%remoteBase%
scp -r ./run.sh %user%@%gamehost%:%remoteBase%
scp -r ./copy.sh %user%@%gamehost%:%remoteBase%
scp -r ./dockercreate.sh %user%@%gamehost%:%remoteBase%
scp -r ./dockerrmi.sh %user%@%gamehost%:%remoteBase%
scp -r ./dockerdele.sh %user%@%gamehost%:%remoteBase%

scp -r ./runLeaf.sh %user%@%gamehost%:%remoteBase%
scp -r ./runRank.sh %user%@%gamehost%:%remoteBase%
scp -r ./runCenter.sh %user%@%gamehost%:%remoteBase%

scp -r ./dockerRunLeaf.sh %user%@%gamehost%:%remoteBase%
scp -r ./dockerRunRank.sh %user%@%gamehost%:%remoteBase%
scp -r ./dockerRunCenter.sh %user%@%gamehost%:%remoteBase%
scp -r ./dockerfilebuild.sh %user%@%gamehost%:%remoteBase%
scp -r ./Dockerfile %user%@%gamehost%:%remoteBase%
scp -r ./docker-compose.yaml %user%@%centerhost%:%remoteBase%

:end
exit