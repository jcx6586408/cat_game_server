@echo off
set excelLocalPath=ws://localhost:3653/

for /r %%f in (*.xlsx) do del %%f

main.exe %excelLocalPath% 0 1000 question1
main.exe %excelLocalPath% 1001 2000 question2
main.exe %excelLocalPath% 2001 3000 question3
main.exe %excelLocalPath% 3001 4000 question4
main.exe %excelLocalPath% 4001 5000 question5
main.exe %excelLocalPath% 5001 6000 question6
main.exe %excelLocalPath% 6001 7000 question7
main.exe %excelLocalPath% 7001 8000 question8

copy *.xls question.xlsx