MKDIR %APPDATA%\gitdo
MKDIR %APPDATA%\gitdo\plugins
XCOPY /s/e/h/k/y %0\..\plugins %APPDATA%\gitdo\plugins
MKDIR %APPDATA%\gitdo\hooks
XCOPY /s/e/h/k/y %0\..\hooks %APPDATA%\gitdo\hooks

@echo off
reg Query "HKLM\Hardware\Description\System\CentralProcessor\0" | find /i "x86" > NUL && set OS=32BIT || set OS=64BIT

echo %OS%
@echo on
if %OS%==32BIT COPY %0\..\gitdo_win_32.exe %windir%\System32\gitdo.exe
if %OS%==64BIT COPY %0\..\gitdo_win_64.exe %windir%\System32\gitdo.exe
PAUSE
