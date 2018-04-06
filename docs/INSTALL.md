Install
===========
## Windows
Right click the `windows_install.bat` file and run as admin.
This will copy the two folders to a new directory in your %AppData% path.
i.e. `C:\Users\YourName\AppData\Roaming\Gitdo`.

It needs to be ran as admin, as it moves the gitdo.exe to your `%WinDir%\System32` folder. I am currently looking for a better solution that still allows gitdo to be ran as a CMD prompt command.

## Mac / Linux
Open a terminal, and navigate to the downloaded and unnzipped directory.
Run the command with `./mac_linux_install.sh`. This will create a new hidden directory in your home folder, i.e. `/Users/YourName/.gitdo`.
The gitdo executable will be moved to the directory `/usr/local/bin` to be in your path. You can move it from here, as long as you move it to a location that is in your `$PATH`

## Notice
The gito directory created is where plugins will be found and ran, so new ones should be placed in their own folder (used as the name) inside the plugins directory. The hooks folder contains the necessary git hooks to run the gitdo commands at different times in the lifecycle. These can be ammended if you wish to add more hook options, just be aware of what gitdo does as to not cause a conflict.


