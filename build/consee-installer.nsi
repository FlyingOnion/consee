; Consee Windows Installer Script
; Created for Consee application installation

; --------------------------------
; Basic Configuration
; --------------------------------
!define APPNAME "Consee"
!define VERSION "0.25.11"
!define COMPANY "Consee Team"
!define DESCRIPTION "Consee Installer"
!define URL "https://github.com/FlyingOnion/consee"

; Installer basic information
Name "${APPNAME}"
Caption "${APPNAME} Installer"
Icon "${NSISDIR}\Contrib\Graphics\Icons\modern-install.ico"
OutFile "consee-installer-v${VERSION}.exe"
InstallDir "D:\consee"
RequestExecutionLevel admin

; --------------------------------
; Interface Settings
; --------------------------------
!include "MUI2.nsh"

; Interface configuration
!define MUI_ABORTWARNING
!define MUI_ICON "${NSISDIR}\Contrib\Graphics\Icons\modern-install.ico"
!define MUI_UNICON "${NSISDIR}\Contrib\Graphics\Icons\modern-uninstall.ico"

; --------------------------------
; Custom Page Declarations
; --------------------------------
!macro MUI_CUSTOMPAGE_INIT COMMAND
  !insertmacro MUI_PAGE_INIT COMMAND
!macroend

; --------------------------------
; Pages Configuration
; --------------------------------

; 1. Welcome Page
!define MUI_WELCOMEPAGE_TITLE $(MUI_TEXT_WELCOME_INFO_TITLE)
!define MUI_WELCOMEPAGE_TEXT $(MUI_TEXT_WELCOME_INFO_TEXT)
!insertmacro MUI_PAGE_WELCOME

; 2. License Page
!define MUI_LICENSEPAGE_TEXT_TOP $(MUI_TEXT_LICENSE_TEXT_TOP)
!define MUI_LICENSEPAGE_FORCE_SELECTION
!insertmacro MUI_PAGE_LICENSE $(LICENSE_FILE)

; 3. Directory Page
!insertmacro MUI_PAGE_DIRECTORY

; 4. Installation Progress Page
!insertmacro MUI_PAGE_INSTFILES

; 5. Custom Page for User Information
Page custom ConsulInformationPage ConsulInformationPageLeave

; 6. Finish Page
!define MUI_FINISHPAGE_TEXT "$(MUI_FINISHPAGE_TEXT_CUSTOM)"
!define MUI_FINISHPAGE_RUN "$INSTDIR\bin\start_consee.ps1"
!define MUI_FINISHPAGE_RUN_TEXT "$(MUI_FINISHPAGE_RUN_TEXT_CUSTOM)"
!insertmacro MUI_PAGE_FINISH

; Uninstaller pages
!insertmacro MUI_UNPAGE_CONFIRM
!insertmacro MUI_UNPAGE_INSTFILES

; --------------------------------
; Languages
; --------------------------------
!insertmacro MUI_LANGUAGE "English"
!insertmacro MUI_LANGUAGE "SimpChinese"
LicenseLangString LICENSE_FILE ${LANG_ENGLISH} "license.txt"
LicenseLangString LICENSE_FILE ${LANG_SIMPCHINESE} "license-cn.txt"

; --------------------------------
; Language Strings (MUI Built-in Variables)
; --------------------------------

; English Strings

LangString MUI_TEXT_WELCOME_INFO_TITLE ${LANG_ENGLISH} "Welcome to the ${APPNAME} Installer"
LangString MUI_TEXT_WELCOME_INFO_TEXT ${LANG_ENGLISH} "This program will guide you through the installation of ${APPNAME}.$\r$\n$\r$\nProgram Version: v${VERSION}$\r$\n$\r$\nProgram URL: ${URL}$\r$\n$\r$\nClick Next to continue."

LangString MUI_TEXT_LICENSE_TEXT_TOP ${LANG_ENGLISH} "This software is licensed under the Mulan PSL v2. Please read the license carefully before installation."

LangString MUI_TEXT_DO_NOT_USE_DEFAULT_ADMIN_TOKEN ${LANG_ENGLISH} "Please make sure that admin token is correctly set. Do NOT use default admin token. This will cause potential access problems and secure risks." 

LangString MUI_FINISHPAGE_TEXT_CUSTOM ${LANG_ENGLISH} "${APPNAME} has been successfully installed on your computer."
LangString MUI_FINISHPAGE_RUN_TEXT_CUSTOM ${LANG_ENGLISH} "Run ${APPNAME}"

; Custom Page Strings
LangString CONSUL_INFO_TITLE ${LANG_ENGLISH} "Consee Initialization"
LangString CONSUL_INFO_SUBTITLE ${LANG_ENGLISH} "Consee depends on Consul service to run, please configure Consul basic information. You can modify it now, or modify it later in '$INSTDIR\config\config.yaml'."
LangString CONSUL_ADDRESS ${LANG_ENGLISH} "Consul Address:"
LangString CONSUL_DATACENTER ${LANG_ENGLISH} "Consul Datacenter:"
LangString CONSUL_ADMIN_TOKEN ${LANG_ENGLISH} "Consul Admin Token:"
LangString CONSUL_DEFAULT_ADMIN_TOKEN_CONFIRM ${LANG_ENGLISH} "Do you want to use the default token? This will bring serious security risks. We don't recommend it."

; Chinese Strings
LangString MUI_TEXT_WELCOME_INFO_TITLE ${LANG_SIMPCHINESE} "欢迎使用 ${APPNAME} 安装程序"
LangString MUI_TEXT_WELCOME_INFO_TEXT ${LANG_SIMPCHINESE} "${APPNAME} 安装程序将引导您完成 ${APPNAME} 的安装过程。$\r$\n$\r$\n程序版本: ${VERSION}$\r$\n$\r$\n程序网址: ${URL}$\r$\n$\r$\n点击“下一步”继续。"

LangString MUI_TEXT_LICENSE_TEXT_TOP ${LANG_SIMPCHINESE} "本软件为开源软件，使用木兰宽松许可证，第二版 Mulan PSL v2 许可证授权，请仔细阅读许可证内容。"

LangString MUI_FINISHPAGE_TEXT_CUSTOM ${LANG_SIMPCHINESE} "${APPNAME} 已成功安装在您的计算机上。"
LangString MUI_FINISHPAGE_RUN_TEXT_CUSTOM ${LANG_SIMPCHINESE} "运行 ${APPNAME}"

; Custom Page Strings
LangString CONSUL_INFO_TITLE ${LANG_SIMPCHINESE} "Consee 初始化"
LangString CONSUL_INFO_SUBTITLE ${LANG_SIMPCHINESE} "Consee 依赖 Consul 服务运行，请配置 Consul 基本信息。您可以现在修改，也可以稍后在 '$INSTDIR\config\config.yaml' 中修改。"
LangString CONSUL_ADDRESS ${LANG_SIMPCHINESE} "Consul 地址:"
LangString CONSUL_DATACENTER ${LANG_SIMPCHINESE} "Consul 数据中心:"
LangString CONSUL_ADMIN_TOKEN ${LANG_SIMPCHINESE} "Consul 管理员令牌:"
LangString CONSUL_DEFAULT_ADMIN_TOKEN_CONFIRM ${LANG_SIMPCHINESE} "您是否要使用默认令牌？这将带来严重的安全风险。我们不建议这样做。"

; --------------------------------
; Variables
; --------------------------------
Var ConseeAddress
Var ConsulAddress
Var ConsulDC
Var ConsulAdminToken

; --------------------------------
; Installer Sections
; --------------------------------

Section "Installing ${APPNAME}" SecMain
    SectionIn RO
    
    SetOutPath "$INSTDIR"
    
    ; Create installation directory
    CreateDirectory "$INSTDIR"
    CreateDirectory "$INSTDIR\bin"
    CreateDirectory "$INSTDIR\config"
    CreateDirectory "$INSTDIR\logs"
    CreateDirectory "$SMPROGRAMS\${APPNAME}"
    
    ; TODO: Add actual binary files here
    File /a /oname=$INSTDIR\bin\consee.exe "bin\consee.exe"
    ; File "bin\*.dll"
    
    ; Create placeholder files (remove when actual binaries are available)
    ; FileOpen $1 "$INSTDIR\bin\consee.exe" w
    ; FileClose $1

    ; Create scripts
    FileOpen $1 "$INSTDIR\bin\start_consee.ps1" w
    FileWrite $1 "Start-Process -WindowStyle hidden -FilePath $INSTDIR\bin\consee.exe$\r$\n"
    FileWrite $1 "echo Consee started successfully. Address: http://$ConseeAddress$\r$\n"
    FileWrite $1 "timeout /t 3"
    FileClose $1

    FileOpen $1 "$INSTDIR\bin\stop_consee.ps1" w
    FileWrite $1 "Stop-Process -Name consee.exe$\r$\n"
    FileWrite $1 "echo Consee stopped successfully.$\r$\n"
    FileWrite $1 "timeout /t 3"
    FileClose $1
    
    ; Create uninstaller
    WriteUninstaller "$INSTDIR\uninstall.exe"
    
    ; Create desktop shortcut
    CreateShortCut "$DESKTOP\Start${APPNAME}.lnk" "$INSTDIR\bin\start_consee.ps1"
    CreateShortCut "$DESKTOP\Stop${APPNAME}.lnk" "$INSTDIR\bin\stop_consee.ps1"
    CreateShortCut "$SMPROGRAMS\${APPNAME}\Start${APPNAME}.lnk" "$INSTDIR\bin\start_consee.ps1"
    CreateShortCut "$SMPROGRAMS\${APPNAME}\Stop${APPNAME}.lnk" "$INSTDIR\bin\stop_consee.ps1"

    
    ; Write registry keys for Windows
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "DisplayName" "${APPNAME}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "UninstallString" "$INSTDIR\uninstall.exe"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "DisplayIcon" "$INSTDIR\bin\consee.exe"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "Publisher" "${COMPANY}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "DisplayVersion" "${VERSION}"
    WriteRegStr HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}" "URLInfoAbout" "${URL}"
    
SectionEnd

; --------------------------------
; Functions
; --------------------------------

Function .onInit
    ; Initialize user information variables
    StrCpy $ConsulAddress "localhost:8500"
    StrCpy $ConsulDC "dc1"
    StrCpy $ConsulAdminToken "f8f0c0c0-f8f0-f8f0-f8f0-f8f0c0c0f8f0"
    StrCpy $ConseeAddress "localhost:3668"

    !insertmacro MUI_LANGDLL_DISPLAY
FunctionEnd

; Custom User Information Page Function
Function ConsulInformationPage
    !insertmacro MUI_HEADER_TEXT "$(CONSUL_INFO_TITLE)" "$(CONSUL_INFO_SUBTITLE)"
    
    nsDialogs::Create /NOUNLOAD 1018
    Pop $0
    
    ${If} $0 == error
        Abort
    ${EndIf}
    
    ; Create labels
    ${NSD_CreateLabel} 0 0 100% 12u "$(CONSUL_ADDRESS)"
    Pop $1
    
    ${NSD_CreateLabel} 0 40u 100% 12u "$(CONSUL_DATACENTER)"
    Pop $2
    
    ${NSD_CreateLabel} 0 80u 100% 24u "$(CONSUL_ADMIN_TOKEN)$\r$\n$(CONSUL_DEFAULT_ADMIN_TOKEN_CONFIRM)"
    Pop $3
    
    ; Create text boxes
    ${NSD_CreateText} 0 15u 100% 12u "$ConsulAddress"
    Pop $4
    
    ${NSD_CreateText} 0 55u 100% 12u "$ConsulDC"
    Pop $5
    
    ${NSD_CreateText} 0 107u 100% 12u "$ConsulAdminToken"
    Pop $6
    
    nsDialogs::Show
    
    ; Get values from text boxes
    ${NSD_GetText} $4 $ConsulAddress
    ${NSD_GetText} $5 $ConsulDC
    ${NSD_GetText} $6 $ConsulAdminToken
FunctionEnd

Function ConsulInformationPageLeave
    ${If} $ConsulAddress == ""
      StrCpy $ConsulAddress "localhost:8500"
    ${EndIf}
    ${If} $ConsulDC == ""
      StrCpy $ConsulDC "dc1"
    ${EndIf}
    FileOpen $1 "$INSTDIR\config\config.yaml" w
    FileWrite $1 "consul:$\r$\n"
    FileWrite $1 "  address: $ConsulAddress$\r$\n"
    FileWrite $1 "  datacenter: $ConsulDC$\r$\n"
    FileWrite $1 "  admin_token: $ConsulAdminToken$\r$\n$\r$\n"
    FileWrite $1 "log_level: info$\r$\n"
    FileWrite $1 "log_file: $INSTDIR\logs\consee.log$\r$\n"
    FileClose $1
FunctionEnd

; --------------------------------
; Uninstaller Section
; --------------------------------

Section "Uninstall"
    
    ; Remove files and directories
    RMDir /r "$INSTDIR\bin"
    RMDir /r "$INSTDIR\config"
    RMDir /r "$INSTDIR\logs"
    Delete "$INSTDIR\uninstall.exe"
    RMDir "$INSTDIR"
    
    ; Remove desktop shortcut
    Delete "$DESKTOP\Start${APPNAME}.lnk"
    Delete "$DESKTOP\Stop${APPNAME}.lnk"
    Delete "$SMPROGRAMS\${APPNAME}\Start${APPNAME}.lnk"
    Delete "$SMPROGRAMS\${APPNAME}\Stop${APPNAME}.lnk"
    RMDir /r "$SMPROGRAMS\${APPNAME}"
    
    ; Remove registry keys
    DeleteRegKey HKLM "Software\Microsoft\Windows\CurrentVersion\Uninstall\${APPNAME}"
    
SectionEnd