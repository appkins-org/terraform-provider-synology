package client

import (
	"encoding/json"
	"strconv"
)

type PackageInfo struct {
	Is_manager               bool
	Support_sharing          bool
	Support_virtual_protocol string
	Hostname                 string
}

type PackageFile struct {
	Path       string
	Name       string
	Isdir      bool
	Additional interface{}
}

type PackageInstallResponse struct {
	Progress string `json:"progress"`
	TaskId   string `json:"taskid"`
}

func GetPackageInfo(apiInfo map[string]InfoData, host string, sid string) (PackageInfo, error) {

	queryString := make(map[string]string)
	queryString["method"] = "get"
	queryString["_sid"] = sid

	_, apiResponse, err := CallAPI(host, "SYNO.Core.Package", apiInfo, queryString)

	if err != nil {
		return PackageInfo{}, err
	}
	var packageInfo PackageInfo

	json.Unmarshal(apiResponse.Data, &packageInfo)

	return packageInfo, nil
}

func CheckPackageInstall(apiInfo map[string]InfoData, host string, sid string) (bool, error) {
	apiName := "SYNO.Core.Package"
	info := apiInfo[apiName]

	queryString := make(map[string]string)
	queryString["method"] = "install_check"
	queryString["_sid"] = sid
	queryString["api"] = apiName
	queryString["version"] = strconv.Itoa(info.MaxVersion)
	wsUrl := host + "/webapi/" + info.Path

	_, apiResponse, err := CallAPI(wsUrl, apiName, apiInfo, queryString)

	return apiResponse.Success, err
}

func InstallPackage(apiInfo map[string]InfoData, host string, sid string, name string, url string) (bool, error) {
	apiName := "SYNO.Core.Package.Installation"
	info := apiInfo[apiName]

	queryString := make(map[string]string)
	queryString["_sid"] = sid
	queryString["api"] = apiName
	queryString["method"] = "uninstall"
	queryString["version"] = strconv.Itoa(info.MaxVersion)
	queryString["name"] = name
	queryString["url"] = url
	queryString["checksum"] = ""
	queryString["filesize"] = "0"
	queryString["type"] = "0"
	queryString["blqinst"] = "false"
	queryString["operation"] = "install"

	_, apiResponse, err := CallAPI(host, apiName, apiInfo, queryString)
	if err != nil {
		return false, err
	}
	return apiResponse.Success, nil
}

func UninstallPackage(apiInfo map[string]InfoData, host string, sid string, packageId string) (bool, error) {
	apiName := "SYNO.Core.Package.Uninstallation"

	queryString := make(map[string]string)
	queryString["api"] = "SYNO.Core.Package.Uninstallation"
	queryString["version"] = "1"
	queryString["method"] = "uninstall"
	queryString["id"] = packageId
	queryString["extra_values"] = "{\"wizard_keep_data\":false,\"wizard_delete_data\":true}"
	queryString["_sid"] = sid
	queryString["api"] = apiName
	wsUrl := host + "/webapi/entry.cgi"

	_, _, err := HttpCall(wsUrl, queryString)

	return err == nil, err
}
