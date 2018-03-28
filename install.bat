MKDIR %APPDATA%\Gitdo
MKDIR %APPDATA%\Gitdo\plugins
XCOPY /s/e/h/k/y %0\..\plugins %APPDATA%\Gitdo\plugins
MKDIR %APPDATA%\Gitdo\hooks
XCOPY /s/e/h/k/y %0\..\plugins %APPDATA%\Gitdo\plugins
COPY %0\..\secrets.json %APPDATA%\Gitdo