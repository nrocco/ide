FROM chef/chefdk

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    wget --content-disposition https://releases.hashicorp.com/vagrant/1.9.5/vagrant_1.9.5_x86_64.deb -O /tmp/vagrant.deb && \
    dpkg -i /tmp/vagrant.deb && \
    apt-get install -qy build-essential && \
    vagrant plugin install vagrant-berkshelf vagrant-cloudstack vagrant-omnibus && \
    apt-get --purge remove -yq build-essential && \
    apt-get --purge autoremove -yq && \
    apt-get autoclean && \
    apt-get clean -yq && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
