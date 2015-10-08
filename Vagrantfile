# -*- mode: ruby -*-
# vi: set ft=ruby :

VAGRANTFILE_API_VERSION = "2"

INSTALL_DEPS=<<EOF
apt-get install memcached

# Install Go
wget https://storage.googleapis.com/golang/go1.5.1.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.5.1.linux-amd64.tar.gz

mkdir -p /opt/go
chmod 777 -R /opt/go

echo 'export GOROOT"=/usr/local/go"' > /etc/profile.d/go.sh
echo 'export GOPATH"=/opt/go"' >> /etc/profile.d/go.sh
echo 'export PATH="$PATH:$GOROOT/bin"' >> /etc/profile.d/go.sh

grep 'cd /vagrant' /home/vagrant/.bashrc ||
    echo 'cd /vagrant' >> /home/vagrant/.bashrc
EOF

Vagrant.configure(VAGRANTFILE_API_VERSION) do |config|

    config.vm.box = "ubuntu/vivid64"

    config.vm.define "dev" do |dev|
        dev.vm.provision "shell", inline: "#{INSTALL_DEPS}"
    end

    if Vagrant.has_plugin?("vagrant-cachier")
        config.cache.scope = :box
    end

end
