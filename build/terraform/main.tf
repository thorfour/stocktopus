# Set the variable value in *.tfvars file
# or using -var="do_token=..." CLI option
variable "do_token" {}

variable "droplet_name" {}
variable "region" {}
variable "ssh_keys" {}

# Configure the DigitalOcean Provider
provider "digitalocean" {
  token = "${var.do_token}"
}

# Create a web server
resource "digitalocean_droplet" "web" {
    image = "ubuntu-18-04-x64"
    name = "${var.droplet_name}"
    region = "${var.region}"
    size = "s-1vcpu-1gb"
    monitoring = true
    ssh_keys = "${var.ssh_keys}"
}
