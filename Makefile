include .envrc

.PHONY: run/api
run/api:
	@echo 'Running Product API...'
	@go run ./cmd/api -port=3000 -env=production -db-dsn=${PRODUCTS_DB_DSN}

.PHONY: db/psql
db/psql:
	psql ${PRODUCTS_DB_DSN}

.PHONY: db/migrations/new
db/migrations/new:
	@echo 'Creating migration files for ${name}...'
	migrate create -seq -ext=.sql -dir=./migrations ${name}

.PHONY: db/migrations/up
db/migrations/up:
	@echo 'Running up migrations...'
	migrate -path=./migrations -database ${PRODUCTS_DB_DSN} up

.PHONY: displayProduct
displayProduct:
	@echo 'Displaying Product'; \
	curl -i localhost:3000/displayProduct/${id} 

.PHONY: deleteProduct
deleteProduct:
	@echo 'Deleting Product'; \
	curl -X DELETE localhost:3000/deleteProduct/${id} 

.PHONY: displayAllProducts
displayAllProducts:
	@echo 'Deleting Product'; \
	curl -i localhost:3000/displayAllProducts

.PHONY: listProducts
listProducts:
	@echo 'Deleting Product'; \
	curl -i localhost:3000/displayAllProducts?${filter}

.PHONY: updateProduct
updateProduct:
	@echo 'Updating Product ${id}'; \
	curl -X PATCH localhost:3000/updateProduct/${id} -d '{"name":"Spoom", "Description":"Mouth!"}'


.PHONY: createReview
createReview:
	@echo 'Creating Review'; \
	BODY=${CREATEREVIEW}; \
	curl -X POST -d "$$BODY" localhost:3000/product/${id}/createReview ; \
	
.PHONY: getReview
getReview:
	@echo 'Displaying Review'; \
	curl -i localhost:3000/product/${id}/getReview/${rid}
	
.PHONY: updateReview
updateReview:
	@echo 'Updating Review'; \
	curl -X PATCH localhost:3000/product/${id}/updateReview/${rid} -d '{"rating":1, "comment":"Yes!"}'
	
.PHONY: deleteReview
deleteReview:
	@echo 'Deleting Product'; \
	curl -X DELETE localhost:3000/product/${id}/deleteReview/${rid}

.PHONY: displayAllReviews
displayAllReviews:
	@echo 'Displaying Reviews'; \
	curl -i localhost:3000/reviews?${filter}

.PHONY: createProduct
createProduct:
	@echo 'Creating Product'; \
    BODY=${CREATEPRODUCT}; \
	curl -i -d "$$BODY" localhost:3000/createProduct ; \
	echo 'create a product'
