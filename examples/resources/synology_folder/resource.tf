terraform {
  required_providers {
    synology = {
      version = "0.2.0"
      source  = "github.com/appkins-org/synology"
    }
  }
}

provider "synology" {
  url      = "<SYNOLOGY_ADDRESS>"
  username = "<SYNOLOGY_USERNAME>"
  password = "<SYNOLOGY_PASSWORD>"
  # these variables can be set as env vars in SYNOLOGY_ADDRESS SYNOLOGY_USERNAME and SYNOLOGY_PASSWORD
}

resource "synology_folder" "my-folder" {
  path = "/home/downloaded/sample-folder"
}