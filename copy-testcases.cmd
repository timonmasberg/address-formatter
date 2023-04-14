echo off

goto(){
echo "Copying OpenCageData testcase files to templates folder"
cp -a address-formatting/testcases/countries/. testcases/
cp -a address-formatting/testcases/other/. testcases/
}

goto $@
exit

:(){
echo Copying OpenCageData testcase files to templates folder
robocopy address-formatting\testcases\countries testcases /COPY:D /E /IT
robocopy address-formatting\testcases\other testcases /COPY:D /E /IT
