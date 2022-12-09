# Finding OS architecture

$is64Bit = Test-Path 'Env:ProgramFiles(x86)'
$PLATFORM="Unidentified Operating System"
# Identifying the Operting system Architecture
if($is64Bit){
    $PLATFORM="windows_amd64"
}else {
    $PLATFORM="windows_386"
}

Add-Type -AssemblyName System.IO.Compression.FileSystem
function Unzip
{
    param([string]$zipfile, [string]$outpath)

    [System.IO.Compression.ZipFile]::ExtractToDirectory($zipfile, $outpath)
}

$OUT_PATH= $args[0] + "\cloudfix-linter\"
New-Item -Path $OUT_PATH -ItemType Directory
# Installing Terraform
$TERRAFORM_VERSION="1.2.6"
$FILE_PATH=$OUT_PATH+"terraform.zip"
Write-Output "Installing terraform........"
Invoke-WebRequest -URI https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_${PLATFORM}.zip -OutFile $FILE_PATH
Unzip $FILE_PATH $OUT_PATH
Remove-Item $FILE_PATH
$TEMP=$OUT_PATH+"terraform.exe"
Set-Alias -Name terraform -Value $TEMP -Scope Global
Write-Output "Terraform installed successfully"


# Installing yor_trace 
$YOR_VERSION="0.1.156"
$FILE_PATH=$OUT_PATH+"yor_trace.zip"
Write-Output "Installing yor_trace........"
Invoke-WebRequest -URI https://github.com/bridgecrewio/yor/releases/download/${YOR_VERSION}/yor_${YOR_VERSION}_${PLATFORM}.zip -OutFile $FILE_PATH
Unzip $FILE_PATH $OUT_PATH
Remove-Item $FILE_PATH
$TEMP=$OUT_PATH+"yor.exe"
Set-Alias -Name yor -Value $TEMP -Scope Global
Write-Output "Yor installed successfully"

#Installing tflint
# higher version have breaking changes to the plugin system and hence we can't install them without changing the plugin
$TFLINT_VERSION="v0.39.3"
$FILE_PATH=Get-Location
$FILE_PATH=$OUT_PATH+"tflint.zip"
Write-Output "Installing tflint........"
Invoke-WebRequest -URI https://github.com/terraform-linters/tflint/releases/download/${TFLINT_VERSION}/tflint_${PLATFORM}.zip -OutFile $FILE_PATH
Unzip $FILE_PATH $OUT_PATH
Remove-Item $FILE_PATH
$TEMP=$OUT_PATH+"tflint.exe"
Set-Alias -Name tflint -Value $TEMP -Scope Global
Write-Output "Tflint installed successfully"


# Install cloudfix-linter
Write-Output "Installing cloudfix-linter........"
$OUT_PATH_CFT=$OUT_PATH+"cloudfix-linter.exe"
Invoke-WebRequest -URI https://github.com/trilogy-group/cloudfix-linter/releases/latest/download/cloudfix-linter_${PLATFORM}.exe -OutFile $OUT_PATH_CFT
$TEMP=$OUT_PATH+"cloudfix-linter.exe"
Set-Alias -Name cloudfix-linter -Value $TEMP -Scope Global
Write-Output "Cloudfix-linter installed successfully"
