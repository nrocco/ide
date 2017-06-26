FROM ubuntu:xenial

ENV DEBIAN_FRONTEND noninteractive
ENV APACHE_LOCK_DIR=/var/lock/apache2
ENV APACHE_LOG_DIR=/var/log/apache2
ENV APACHE_PID_FILE=/var/run/apache2.pid
ENV APACHE_RUN_DIR=/var/run/apache2
ENV APACHE_RUN_GROUP=www-data
ENV APACHE_RUN_USER=www-data
ENV COMPOSER_ALLOW_SUPERUSER 1
ENV COMPOSER_HOME /composer
ENV IDE_WEBROOT=/var/www
ENV LANG=C
ENV PATH /composer/vendor/bin:$PATH

# Install apt dependencies
RUN apt-get update && \
    apt-get -yq install \
        apache2 \
        libapache2-mod-php7.0 \
        apache2-utils \
        git \
        snmp \
        php-xdebug \
        php7.0-bcmath \
        php7.0-cli \
        php7.0-curl \
        php7.0-intl \
        php7.0-mbstring \
        php7.0-mcrypt \
        php7.0-mysql \
        php7.0-snmp \
        php7.0-soap \
        php7.0-sqlite3 \
        php7.0-xmlrpc \
        php7.0-xsl \
        php7.0-xml \
        php7.0-zip && \
    rm /etc/php/7.0/cli/conf.d/20-xdebug.ini && \
    rm -rf /var/lib/apt/lists/*

# Install composer
RUN php -r "copy('https://getcomposer.org/installer', 'composer-setup.php');" && \
    php -r "if (hash_file('SHA384', 'composer-setup.php') === '669656bab3166a7aff8a7506b8cb2d1c292f042046c5a994c43155c0be6190fa0355160742ab2e1c88d40d5be660b410') { echo 'Installer verified'; } else { echo 'Installer corrupt'; unlink('composer-setup.php'); } echo PHP_EOL;" && \
    php composer-setup.php --install-dir=/usr/local/bin --filename=composer && \
    php -r "unlink('composer-setup.php');"

# HACK: get rid of this hard coded uid and gid
RUN useradd -u 1000 -g 100 -s /bin/bash user && \
    echo 'StrictHostKeyChecking no' >> /etc/ssh/ssh_config && echo 'UserKnownHostsFile /dev/null' >> /etc/ssh/ssh_config && \
    a2enmod vhost_alias rewrite

# logs should go to stdout / stderr
RUN set -ex \
    && mkdir -p $APACHE_RUN_DIR \
    && ln -sfT /dev/stderr "$APACHE_LOG_DIR/error.log" \
    && ln -sfT /dev/stdout "$APACHE_LOG_DIR/access.log" \
    && ln -sfT /dev/stdout "$APACHE_LOG_DIR/other_vhosts_access.log" \
    && rm -f /etc/php/7.0/apache2/conf.d/20-xdebug.ini

COPY apache.conf /etc/apache2/sites-enabled/000-default.conf

ENTRYPOINT ["php"]
