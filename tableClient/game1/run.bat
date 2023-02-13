@echo off
set excelLocalPath=wss://yinghuo-1.easygame2021.com:5101/

for /r %%f in (*.xlsx) do del %%f

main.exe %excelLocalPath% 0 1000 question1
main.exe %excelLocalPath% 1001 2000 question2
main.exe %excelLocalPath% 2001 3000 question3
main.exe %excelLocalPath% 3001 4000 question4
main.exe %excelLocalPath% 4001 5000 question5
main.exe %excelLocalPath% 5001 6000 question6
main.exe %excelLocalPath% 6001 7000 question7
main.exe %excelLocalPath% 7001 8000 question8
main.exe %excelLocalPath% 8001 9000 question9
main.exe %excelLocalPath% 9001 10000 question10
main.exe %excelLocalPath% 10001 11000 question11
