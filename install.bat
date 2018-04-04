MKDIR %APPDATA%\gitdo
MKDIR %APPDATA%\gitdo\plugins
XCOPY /s/e/h/k/y %0\..\plugins %APPDATA%\gitdo\plugins
MKDIR %APPDATA%\gitdo\hooks
XCOPY /s/e/h/k/y %0\..\hooks %APPDATA%\gitdo\hooks
COPY %0\..\gitdo_windows.exe %windir%\System32\gitdo.exe
