package client

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type readResponse struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

type NetworkResponse struct {
	Networks []Network `json:"networks"`
}

type Network struct {
	NetworkID   string `json:"network_id"`
	NetworkName string `json:"network_name"`
}

type StorageResponse struct {
	Storages []Storage `json:"storages"`
}

type Storage struct {
	HostID      string `json:"host_id"`
	HostName    string `json:"host_name"`
	Size        int    `json:"size"`
	Status      string `json:"status"`
	StorageID   string `json:"storage_id"`
	StorageName string `json:"storage_name"`
	Used        int    `json:"used"`
	VolumePath  string `json:"volume_path"`
}

type HostResponse struct {
	Hosts []Host `json:"hosts"`
}

type Host struct {
	FreeCpuCore  int    `json:"free_cpu_core"`
	FreeRamSize  int    `json:"free_ram_size"`
	HostID       string `json:"host_id"`
	HostName     string `json:"host_name"`
	Status       string `json:"status"`
	TotalCpuCore int    `json:"total_cpu_core"`
	TotalRamSize int    `json:"total_ram_size"`
}

type Guest struct {
	Autorun     int     `json:"autorun"`
	Description string  `json:"description"`
	GuestId     string  `json:"guest_id"`
	GuestName   string  `json:"guest_name"`
	Status      string  `json:"status"`
	StorageId   string  `json:"storage_id"`
	StorageName string  `json:"storage_name"`
	VcpuNum     int     `json:"vcpu_num"`
	Vdisks      []VDisk `json:"vdisks"`
	Vnics       []VNic  `json:"vnics"`
	VramSize    int     `json:"vram_size"`
}

type VNic struct {
	Mac         string `json:"mac"`
	Model       int    `json:"model"`
	NetworkID   string `json:"network_id"`
	NetworkName string `json:"network_name"`
	VnicID      string `json:"vnic_id"`
}

type VDisk struct {
	Controller int    `json:"controller"`
	Unmap      bool   `json:"unmap"`
	VdiskId    string `json:"vdisk_id"`
	VdiskSize  int    `json:"vdisk_size"`
}

type GuestInfo struct {
	auto_clean_task bool
	guest_id        string
	progress        int
	status          string
}

type CreateGuestVnicRequest struct {
}

type CreateGuestResponse struct {
	task_id string
}

func CreateGuest(apiInfo map[string]InfoData, host string, sid string, name string, storageId string, storageName string, vnics []interface{}, vdisks []interface{}) (CreateGuestResponse, error) {
	apiName := "SYNO.Virtualization.API.Guest"
	info := apiInfo[apiName]

	vnicList := createValidRequestMap(vnics, []string{"mac", "network_id", "network_name"})
	vdiskList := createValidRequestMap(vdisks, []string{"create_type", "vdisk_size", "image_id", "image_name"})

	vnicsString, _ := json.Marshal(vnicList)
	vdisksString, _ := json.Marshal(vdiskList)

	queryString := make(map[string]string)
	queryString["_sid"] = sid
	queryString["api"] = apiName
	queryString["method"] = "create"
	queryString["version"] = strconv.Itoa(info.MaxVersion)
	queryString["guest_name"] = name
	if storageId != "" {
		queryString["storage_id"] = storageId
	}
	if storageName != "" {
		queryString["storage_name"] = storageName
	}
	queryString["vnics"] = string(vnicsString)
	queryString["vdisks"] = string(vdisksString)

	wsUrl := host + "/webapi/entry.cgi"
	_, body, err := HttpCall(wsUrl, queryString)

	if err != nil {
		return CreateGuestResponse{}, err
	}

	var CreateGuestResponse CreateGuestResponse
	json.Unmarshal(body, &CreateGuestResponse)

	return CreateGuestResponse, nil
}

func ReadGuest(apiInfo map[string]InfoData, host string, sid string, name string) (Guest, error) {
	apiName := "SYNO.Virtualization.API.Guest"
	info := apiInfo[apiName]

	queryString := make(map[string]string)
	queryString["_sid"] = sid
	queryString["api"] = apiName
	queryString["method"] = "get"
	queryString["version"] = strconv.Itoa(info.MaxVersion)
	queryString["guest_name"] = name
	queryString["additional"] = "true"

	wsUrl := host + "/webapi/entry.cgi"

	_, body, err := HttpCall(wsUrl, queryString)
	if err != nil {
		return Guest{}, err
	}

	response := readResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return Guest{}, err
	}

	guestData, ok := response.Data.(map[string]interface{})
	if !ok {
		return Guest{}, fmt.Errorf("invalid guest data")
	}

	guestJson, err := json.Marshal(guestData)
	if err != nil {
		return Guest{}, err
	}

	guest := Guest{}
	err = json.Unmarshal(guestJson, &guest)
	if err != nil {
		return Guest{}, err
	}

	return guest, nil
}

func SetGuest(apiInfo map[string]InfoData, host string, sid string, oldName string, name string, autorun int, description string, vcpuNum int, vramSize int) ([]byte, error) {
	apiName := "SYNO.Virtualization.API.Guest"
	info := apiInfo[apiName]

	queryString := make(map[string]string)
	queryString["_sid"] = sid
	queryString["api"] = apiName
	queryString["method"] = "set"
	queryString["version"] = strconv.Itoa(info.MaxVersion)
	queryString["guest_name"] = oldName
	if oldName != name && name != "" {
		queryString["new_guest_name"] = name
	}
	queryString["autorun"] = strconv.Itoa(autorun)
	queryString["description"] = description
	if vcpuNum != 0 {
		queryString["vcpu_num"] = strconv.Itoa(vcpuNum)
	}
	if vramSize != 0 {
		queryString["vram_size"] = strconv.Itoa(vramSize)
	}

	wsUrl := host + "/webapi/entry.cgi"

	_, body, err := HttpCall(wsUrl, queryString)

	if err != nil {
		return nil, err
	}

	return body, nil
}

func UpdateGuest(apiInfo map[string]InfoData, host string, sid string, name string, newName string) (int, error) {
	apiName := "SYNO.Virtualization.API.Guest"
	info := apiInfo[apiName]

	queryString := make(map[string]string)
	queryString["_sid"] = sid
	queryString["api"] = apiName
	queryString["method"] = "set"
	queryString["version"] = strconv.Itoa(info.MaxVersion)
	queryString["guest_name"] = name
	queryString["new_guest_name"] = newName

	wsUrl := host + "/webapi/entry.cgi"

	statusCode, _, err := HttpCall(wsUrl, queryString)
	if err != nil {
		return statusCode, err
	}

	return statusCode, nil
}

func DeleteGuest(apiInfo map[string]InfoData, host string, sid string, name string) (int, error) {
	apiName := "SYNO.Virtualization.API.Guest"
	info := apiInfo[apiName]

	queryString := make(map[string]string)
	queryString["_sid"] = sid
	queryString["api"] = apiName
	queryString["method"] = "delete"
	queryString["version"] = strconv.Itoa(info.MaxVersion)
	queryString["guest_name"] = name

	wsUrl := host + "/webapi/entry.cgi"

	statusCode, _, err := HttpCall(wsUrl, queryString)
	if err != nil {
		return statusCode, err
	}

	return statusCode, nil
}

func PowerGuest(apiInfo map[string]InfoData, host string, sid string, name string, state bool) (int, error) {
	apiName := "SYNO.Virtualization.API.Guest.Action"
	info := apiInfo[apiName]

	queryString := make(map[string]string)
	queryString["_sid"] = sid
	queryString["api"] = apiName
	if state {
		queryString["method"] = "poweron"
	} else {
		queryString["method"] = "poweroff"
	}
	queryString["version"] = strconv.Itoa(info.MaxVersion)
	queryString["guest_name"] = name

	wsUrl := host + "/webapi/entry.cgi"

	statusCode, _, err := HttpCall(wsUrl, queryString)
	if err != nil {
		return statusCode, err
	}

	return statusCode, nil
}

func ListNetworks(apiInfo map[string]InfoData, host string, sid string) (NetworkResponse, error) {
	apiName := "SYNO.Virtualization.API.Network"
	info := apiInfo[apiName]

	queryString := make(map[string]string)
	queryString["_sid"] = sid
	queryString["api"] = apiName
	queryString["version"] = strconv.Itoa(info.MaxVersion)
	queryString["method"] = "list"

	wsUrl := host + "/webapi/entry.cgi"

	_, body, err := HttpCall(wsUrl, queryString)
	if err != nil {
		return NetworkResponse{}, err
	}

	response := readResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return NetworkResponse{}, err
	}

	storageData, ok := response.Data.(map[string]interface{})
	if !ok {
		return NetworkResponse{}, fmt.Errorf("invalid NetworkResponse data")
	}

	storageJson, err := json.Marshal(storageData)
	if err != nil {
		return NetworkResponse{}, err
	}

	networkResponse := NetworkResponse{}
	err = json.Unmarshal(storageJson, &networkResponse)
	if err != nil {
		return NetworkResponse{}, err
	}

	return networkResponse, nil
}

func ListStorages(apiInfo map[string]InfoData, host string, sid string) (StorageResponse, error) {
	apiName := "SYNO.Virtualization.API.Storage"
	info := apiInfo[apiName]

	queryString := make(map[string]string)
	queryString["_sid"] = sid
	queryString["api"] = apiName
	queryString["version"] = strconv.Itoa(info.MaxVersion)
	queryString["method"] = "list"

	wsUrl := host + "/webapi/entry.cgi"

	_, body, err := HttpCall(wsUrl, queryString)
	if err != nil {
		return StorageResponse{}, err
	}

	response := readResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return StorageResponse{}, err
	}

	storageData, ok := response.Data.(map[string]interface{})
	if !ok {
		return StorageResponse{}, fmt.Errorf("invalid StorageResponse data")
	}

	storageJson, err := json.Marshal(storageData)
	if err != nil {
		return StorageResponse{}, err
	}

	storageResponse := StorageResponse{}
	err = json.Unmarshal(storageJson, &storageResponse)
	if err != nil {
		return StorageResponse{}, err
	}

	return storageResponse, nil
}

func ListHosts(apiInfo map[string]InfoData, host string, sid string) (HostResponse, error) {
	apiName := "SYNO.Virtualization.API.Host"
	info := apiInfo[apiName]

	queryString := make(map[string]string)
	queryString["_sid"] = sid
	queryString["api"] = apiName
	queryString["version"] = strconv.Itoa(info.MaxVersion)
	queryString["method"] = "list"

	wsUrl := host + "/webapi/entry.cgi"

	_, body, err := HttpCall(wsUrl, queryString)
	if err != nil {
		return HostResponse{}, err
	}

	response := readResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return HostResponse{}, err
	}

	storageData, ok := response.Data.(map[string]interface{})
	if !ok {
		return HostResponse{}, fmt.Errorf("invalid HostResponse data")
	}

	storageJson, err := json.Marshal(storageData)
	if err != nil {
		return HostResponse{}, err
	}

	hostResponse := HostResponse{}
	err = json.Unmarshal(storageJson, &hostResponse)
	if err != nil {
		return HostResponse{}, err
	}

	return hostResponse, nil
}

func (g Guest) String() string {
	str := fmt.Sprintf("Guest:\n\tGuestName: %s\n\tGuestId: %s\n\tAutorun: %d\n\tDescription: %s\n\tStatus: %s\n\tStorageName: %s\n\tStorageId: %s\n\tVcpuNum: %d\n\tVramSize: %d\n\tVdisks: [\n", g.GuestName, g.GuestId, g.Autorun, g.Description, g.Status, g.StorageName, g.StorageId, g.VcpuNum, g.VramSize)
	for _, vdisk := range g.Vdisks {
		str += fmt.Sprintf("\t\t%s\n", vdisk.String())
	}
	str += "\t]\n\tVnics: [\n"
	for _, vnic := range g.Vnics {
		str += fmt.Sprintf("\t\t%s\n", vnic.String())
	}
	str += "\t]\n"
	return str
}

func (vnic VNic) String() string {
	return fmt.Sprintf("VNic:\n\tMac: %s\n\tModel: %d\n\tNetworkID: %s\n\tNetworkName: %s\n\tVnicID: %s", vnic.Mac, vnic.Model, vnic.NetworkID, vnic.NetworkName, vnic.VnicID)
}

func (vdisk VDisk) String() string {
	return fmt.Sprintf("VDisk:\n\tController: %d\n\tUnmap: %t\n\tVdiskId: %s\n\tVdiskSize: %d", vdisk.Controller, vdisk.Unmap, vdisk.VdiskId, vdisk.VdiskSize)
}

func (sr StorageResponse) String() string {
	str := "Storages:\n\t"
	for _, storage := range sr.Storages {
		str += fmt.Sprintf("\t\t%s\n", storage.String())
	}
	return str
}

func (s Storage) String() string {
	return fmt.Sprintf("HostID: %s\nHostName: %s\nSize: %d\nStatus: %s\nStorageID: %s\nStorageName: %s\nUsed: %d\nVolumePath: %s\n",
		s.HostID, s.HostName, s.Size, s.Status, s.StorageID, s.StorageName, s.Used, s.VolumePath)
}

func createValidRequestMap(input []interface{}, allowedKeys []string) []map[string]interface{} {
	var output []map[string]interface{}

	for _, elem := range input {
		if v, ok := elem.(map[string]interface{}); ok {
			filtered := make(map[string]interface{})
			for _, k := range allowedKeys {
				if val, found := v[k]; found && val != nil && val != "" {
					filtered[k] = val
				}
			}
			output = append(output, filtered)
		}
	}

	return output
}
