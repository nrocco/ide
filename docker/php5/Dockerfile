FROM ubuntu:precise

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get -yq install \
        git \
        snmp \
        php5-cli \
        php5-curl \
        php5-mcrypt \
        php5-sqlite \
        php5-intl \
        php5-mysql \
        php5-snmp \
        php5-xdebug \
        php5-xsl \
        php5-xmlrpc && \
    rm -rf /var/lib/apt/lists/*

RUN php -r "copy('https://getcomposer.org/installer', 'composer-setup.php');" && \
    php -r "if (hash_file('SHA384', 'composer-setup.php') === '669656bab3166a7aff8a7506b8cb2d1c292f042046c5a994c43155c0be6190fa0355160742ab2e1c88d40d5be660b410') { echo 'Installer verified'; } else { echo 'Installer corrupt'; unlink('composer-setup.php'); } echo PHP_EOL;" && \
    php composer-setup.php --install-dir=/usr/local/bin --filename=composer && \
    php -r "unlink('composer-setup.php');"

RUN useradd -u 1000 -g 100 -s /bin/bash user

ENV COMPOSER_HOME /composer
ENV PATH /composer/vendor/bin:$PATH
ENV COMPOSER_ALLOW_SUPERUSER 1

RUN echo 'StrictHostKeyChecking no' >> /etc/ssh/ssh_config && echo 'UserKnownHostsFile /dev/null' >> /etc/ssh/ssh_config

ENTRYPOINT ["php"]
