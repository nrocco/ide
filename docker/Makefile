.PHONY: stack
stack: containers
	docker stack deploy --prune --compose-file ide.yml ide


.PHONY: containers
containers: php5 php7 apache22 apache24 chefdk


.PHONY: apache22
apache22: php5
	docker build -t ide:apache22 apache22


.PHONY: apache24
apache24: php7
	docker build -t ide:apache24 apache24


.PHONY: php5
php5:
	docker build -t ide:php5 php5


.PHONY: php7
php7:
	docker build -t ide:php7 php7


.PHONY: chefdk
chefdk:
	docker build -t ide:chefdk chefdk


.PHONY: nodejs
nodejs:
	docker build -t ide:nodejs nodejs


.PHONY: ctags
ctags:
	docker build -t ide:ctags ctags
