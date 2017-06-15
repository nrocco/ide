FROM ide:php7

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update && \
    apt-get -yq install \
        apache2 \
        libapache2-mod-php7.0 \
        apache2-utils && \
    rm -rf /var/lib/apt/lists/*

RUN a2enmod vhost_alias rewrite

ENV APACHE_RUN_USER=www-data \
    APACHE_RUN_GROUP=www-data \
    APACHE_PID_FILE=/var/run/apache2.pid \
    APACHE_RUN_DIR=/var/run/apache2 \
    APACHE_LOCK_DIR=/var/lock/apache2 \
    APACHE_LOG_DIR=/var/log/apache2 \
    LANG=C

# logs should go to stdout / stderr
RUN set -ex \
    && mkdir -p $APACHE_RUN_DIR \
    && ln -sfT /dev/stderr "$APACHE_LOG_DIR/error.log" \
    && ln -sfT /dev/stdout "$APACHE_LOG_DIR/access.log" \
    && ln -sfT /dev/stdout "$APACHE_LOG_DIR/other_vhosts_access.log"

EXPOSE 80

COPY apache.conf /etc/apache2/sites-enabled/000-default.conf

WORKDIR /var/www
VOLUME ["/var/www"]

ENTRYPOINT ["apache2", "-DFOREGROUND"]
