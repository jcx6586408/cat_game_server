@echo off

set host=118.195.244.48
set centerhost=82.157.137.166
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
for /r %localPath% %%i in (*) do (scp -r %%i %user%@%centerhost%:%remotePath%)
goto end

:conf
for /r %conflocalPath% %%i in (*) do (scp -r %%i %user%@%centerhost%:%confPath%)
goto end

:excel
for /r %excelLocalPath% %%i in (*) do (scp -r %%i %user%@%centerhost%:%excelRemotePath%)
goto end

:iplocation
for /r %iplocationLocalPath% %%i in (*) do (scp -r %%i %user%@%centerhost%:%iplocationRemotePath%)
goto end

:ssl
for /r %sslLocalPath% %%i in (*) do (scp -r %%i %user%@%centerhost%:%sslRemotePath%)
goto end

:keystore
@REM 中心服配置文件
scp -r ./kill.sh %user%@%centerhost%:%remoteBase%
scp -r ./run.sh %user%@%centerhost%:%remoteBase%
scp -r ./copy.sh %user%@%centerhost%:%remoteBase%
scp -r ./dockercreate.sh %user%@%centerhost%:%remoteBase%
scp -r ./dockerrmi.sh %user%@%centerhost%:%remoteBase%
scp -r ./dockerdele.sh %user%@%centerhost%:%remoteBase%
scp -r ./chmod.sh %user%@%centerhost%:%remoteBase%

scp -r ./runLeaf.sh %user%@%centerhost%:%remoteBase%
scp -r ./runRank.sh %user%@%centerhost%:%remoteBase%
scp -r ./runCenter.sh %user%@%centerhost%:%remoteBase%
scp -r ./runSheep.sh %user%@%centerhost%:%remoteBase%

scp -r ./dockerRunLeaf.sh %user%@%centerhost%:%remoteBase%
scp -r ./dockerRunRank.sh %user%@%centerhost%:%remoteBase%
scp -r ./dockerRunCenter.sh %user%@%centerhost%:%remoteBase%
scp -r ./dockerfilebuild.sh %user%@%centerhost%:%remoteBase%
scp -r ./Dockerfile %user%@%centerhost%:%remoteBase%
scp -r ./docker-compose.yaml %user%@%centerhost%:%remoteBase%
scp -r ./dockercomposePs.sh %user%@%centerhost%:%remoteBase%
scp -r ./dockercomposeRun.sh %user%@%centerhost%:%remoteBase%
scp -r ./dockercomposeRunD.sh %user%@%centerhost%:%remoteBase%
scp -r ./dockercomposeStop.sh %user%@%centerhost%:%remoteBase%
scp -r ./dockercomposeLogs.sh %user%@%centerhost%:%remoteBase%


:end
exit