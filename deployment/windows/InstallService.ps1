$source = "$PSScriptRoot/../bin/vmware-controller.exe"
$destination = "C:\Program Files\VmwareController\vmware-controller.exe"

Copy-Item  $source $destination -Force -Recurse

sc.exe create VMWareControllerServer binPath= "$destination"
