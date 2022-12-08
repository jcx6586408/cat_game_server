@echo off

set localPath=.\bin\bin\

set host=118.195.244.48
set user=root
set remotePath=/home/sheep/bin/

for /r %localPath% %%i in (*) do (scp %%i %user%@%host%:%remotePath%)
