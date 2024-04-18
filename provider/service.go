package provider

import (
	"log"
	"path/filepath"
	"strconv"

	"github.com/appkins-org/terraform-provider-synology/client"
)

type FileItemService struct {
	synologyClient client.SynologyClient
}

func (service FileItemService) Create(filename string, contents []byte) error {
	log.Println("Create " + string(contents))
	path, filename := getPathAndFilenameFromFullPath(filename)
	return service.synologyClient.Upload(path, true, true, filename, contents)
}

func (service FileItemService) Read(filename string) ([]byte, error) {
	return service.synologyClient.Download(filename)
}

func (service FileItemService) Update(filename string, contents []byte) ([]byte, error) {
	log.Println("Update " + string(contents))

	path, filename := getPathAndFilenameFromFullPath(filename)
	err := service.synologyClient.Upload(path, true, true, filename, contents)
	return contents, err
}

func (service FileItemService) Delete(filename string) error {
	return service.synologyClient.Delete(filename, false)

}

func getPathAndFilenameFromFullPath(fullPath string) (string, string) {
	return filepath.Dir(fullPath), filepath.Base(fullPath)
}

type FolderItemService struct {
	synologyClient client.SynologyClient
}

func (service FolderItemService) Create(path string) error {
	log.Println("Create Folder" + string(path))
	basePath, name := getPathAndFilenameFromFullPath(path)
	_, error := service.synologyClient.CreateFolder(basePath, name, true, "")
	return error
}

func (service FolderItemService) Delete(path string) error {
	log.Println("Delete Folder" + string(path))
	return service.synologyClient.Delete(path, true)

}

type GuestService struct {
	synologyClient client.SynologyClient
}

func (service GuestService) Create(name string, storageID string, storageName string, vnics []interface{}, vdisks []interface{}) error {
	log.Println("Create VMM Guest " + string(name))
	_, err := service.synologyClient.CreateGuest(name, storageID, storageName, vnics, vdisks)
	return err
}

func (service GuestService) Set(oldName string, name string, autorun int, description string, vcpuNum int, vramSize int) error {
	log.Println("Setting values on VMM Guest " + string(name))
	err := service.synologyClient.SetGuest(oldName, name, autorun, description, vcpuNum, vramSize)
	return err
}

func (service GuestService) Read(name string) (client.Guest, error) {
	log.Println("Read VMM Guest " + string(name))
	content, err := service.synologyClient.ReadGuest(name)
	return content, err
}

func (service GuestService) Update(name string, newName string) error {
	log.Println("Update VMM Guest from " + string(name) + " to " + string(newName))
	err := service.synologyClient.UpdateGuest(name, newName)
	return err
}

func (service GuestService) Delete(name string) error {
	log.Println("Delete VMM Guest " + string(name))
	err := service.synologyClient.DeleteGuest(name)
	return err
}

func (service GuestService) Power(name string, state bool) error {
	log.Println("Setting VMM Guest " + string(name) + " poweron state to " + strconv.FormatBool(state))
	err := service.synologyClient.PowerGuest(name, state)
	return err
}

type StorageGuestService struct {
	synologyClient client.SynologyClient
}

func (service StorageGuestService) Read() (client.StorageResponse, error) {
	content, err := service.synologyClient.ReadStorageGuest()
	return content, err
}

type NetworkGuestService struct {
	synologyClient client.SynologyClient
}

func (service NetworkGuestService) Read() (client.NetworkResponse, error) {
	content, err := service.synologyClient.ReadNetworkGuest()
	return content, err
}

type HostGuestService struct {
	synologyClient client.SynologyClient
}

func (service HostGuestService) Read() (client.HostResponse, error) {
	content, err := service.synologyClient.ReadHostGuest()
	return content, err
}
