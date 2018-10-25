# Set the variable value in *.tfvars file
# or using -var="do_token=..." CLI option
variable "do_token" {}

variable "region" {
    default = "nyc3"
}

variable "redis_pw" {}
variable "hostname" {
    default = "beta.stocktopus.io"
}
variable "ssh_key_path" {
    default = "~/.ssh/id_rsa.pub"
}
variable "ssh_private_key_path" {
    default = "~/.ssh/id_rsa"
}

# Configure the DigitalOcean Provider
provider "digitalocean" {
  token = "${var.do_token}"
}

resource "digitalocean_ssh_key" "default" {
   name = "default" 
   public_key = "${file("${var.ssh_key_path}")}"
}

# Create redis.conf file
data "template_file" "redisconf" {
    template = "${file("redis.conf")}"

    vars {
        redis_pw = "${var.redis_pw}"
    }
}

# Create separate redis droplet
resource "digitalocean_droplet" "redis" {
    image = "ubuntu-18-04-x64"
    name = "stocktopus-redis"
    region = "${var.region}"
    size = "s-1vcpu-1gb"
    monitoring = true
    ssh_keys = ["${digitalocean_ssh_key.default.fingerprint}"]
    tags = ["stocktopus", "terraform"]

    provisioner "file" {
        content = "${data.template_file.redisconf.rendered}"
        destination = "/etc/redis.conf"

        connection {
            type = "ssh"
            user = "root"
            private_key = "${file("${var.ssh_private_key_path}")}"
        }
    }

    provisioner "remote-exec" {
        script = "docker.sh"

        connection {
            type = "ssh"
            user = "root"
            private_key = "${file("${var.ssh_private_key_path}")}"
        }
    }

    provisioner "remote-exec" {
        inline = [
            "mkdir /data",
            "docker pull redis",
            "docker run --name redis -d -p 6379:6379 -v /data:/data redis redis-server /etc/redis.conf",
        ]

        connection {
            type = "ssh"
            user = "root"
            private_key = "${file("${var.ssh_private_key_path}")}"
        }
    }
}

output "stocktopus_redis_ip" {
    description = "redis ipv4 address"
    value = "${digitalocean_droplet.redis.ipv4_address}"
}

# Create the stocktopus main droplet
resource "digitalocean_droplet" "stocktopus" {
    image = "ubuntu-18-04-x64"
    name = "stocktopus"
    region = "${var.region}"
    size = "s-1vcpu-1gb"
    monitoring = true
    tags = ["stocktopus", "terraform"]
    ssh_keys = ["${digitalocean_ssh_key.default.fingerprint}"]

    provisioner "remote-exec" {
        script = "docker.sh"

        connection {
            type = "ssh"
            user = "root"
            private_key = "${file("${var.ssh_private_key_path}")}"
        }
    }

    provisioner "remote-exec" {
        inline = [
            "mkdir /cert",
            "docker pull quay.io/thorfour/stocktopus:v1.3.2",
            "docker run --name stocktopus -d -p 80:80 -p 443:443 -e REDISADDR=${digitalocean_droplet.redis.ipv4_address} -e REDISPW=${var.redis_pw} -v /cert:/cert quay.io/thorfour/stocktopus:v1.3.2 /server -host ${var.hostname} -c /cert",
        ]

        connection {
            type = "ssh"
            user = "root"
            private_key = "${file("${var.ssh_private_key_path}")}"
        }
    }
}

output "stocktopus_ip" {
    description = "stocktopus main ipv4"
    value = "${digitalocean_droplet.stocktopus.ipv4_address}"
}
