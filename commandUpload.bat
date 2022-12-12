@echo off

set host=118.195.244.48
set user=root

set remoteBase=/home/Test/

set localPath=.\
set remotePath=%remoteBase%


for /r %localPath% %%i in (*.sh) do (scp -r %%i %user%@%host%:%remotePath%)
