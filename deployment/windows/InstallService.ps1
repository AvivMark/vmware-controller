if (-Not ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] 'Administrator')) {
    if ([int](Get-CimInstance -Class Win32_OperatingSystem | Select-Object -ExpandProperty BuildNumber) -ge 6000) {
     $CommandLine = "-File `"" + $MyInvocation.MyCommand.Path + "`" " + $MyInvocation.UnboundArguments
     Start-Process -FilePath PowerShell.exe -Verb Runas -ArgumentList $CommandLine
     Exit
    }
   }
   

$binName = "vmware-controller.exe"
$source = "$PSScriptRoot/../../bin/$binName"
$programFolder =  "C:\Program Files\VmwareController"

$destination = "$programFolder\$binName"
New-Item -Path "$destination" -Type Directory -Force
Copy-Item  $source $destination -Force -Recurse


$serviceName = "vmware-controller"
$service = Get-Service -Name $serviceName -ErrorAction SilentlyContinue
if($service -ne $null)
{
    sc.exe delete $serviceName
} 

sc.exe create $serviceName binPath= "$destination"


sc.exe config $serviceName start= auto
sc.exe start $serviceName 

[System.Windows.Forms.MessageBox]::Show("Deployment completed", "vmware-controller-deployment", [System.Windows.Forms.MessageBoxButtons]::OKCancel)

Read-Host
