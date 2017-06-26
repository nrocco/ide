.PHONY: stack
stack: containers
	docker stack deploy --prune --compose-file ide.yml ide


.PHONY: containers
containers: php5 php7 chefdk ctags nodejs


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
