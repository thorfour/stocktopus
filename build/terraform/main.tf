# Set the variable value in *.tfvars file
# or using -var="do_token=..." CLI option
variable "do_token" {}

variable "droplet_name" {}
variable "region" {}
variable "ssh_keys" {}
variable "redis_pw" {}
variable "redis_addr" {}

# Configure the DigitalOcean Provider
provider "digitalocean" {
  token = "${var.do_token}"
}

# Create the stocktopus main droplet
resource "digitalocean_droplet" "stocktopus" {
    image = "ubuntu-18-04-x64"
    name = "stocktopus"
    region = "${var.region}"
    size = "s-1vcpu-1gb"
    monitoring = true
    ssh_keys = "${var.ssh_keys}"

    provisioner "docker_install" {
        command = "sudo apt install docker"
    }

    provisioner "stocktopus_install" {
        command = "docker pull quay.io/thorfour/stocktopus:v1.3.2"
    }

    provisioner "stocktopus_run" {
        command = "docker run -d -p 80:80 -p 443:443 -e REDISADDR=${var.redis_addr} -e REDISPW=${var.redis_pw} quay.io/thorfour/stocktopus:v1.3.2"
    }
}

resource "digitalocean_droplet" "redis" {

}
