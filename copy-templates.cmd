echo off

goto(){
echo "Copying OpenCageData config files to templates folder"
cp -a OpenCageData-address-formatting/conf/. templates/
}

goto $@
exit

:(){
echo Copying OpenCageData config files to templates folder
robocopy OpenCageData-address-formatting\conf templates /COPY:D /E /IT

