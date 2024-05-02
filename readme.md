
1o precisa criar a tabela, pode ser feito com o seguinte comando:
```
CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))
```

ou então se você executar algum dos testes de repositório golangtarefa3/internal/infra/database/order_repository_test.go

qualquer um dos testes irá executar o SetupSuite que irá criar a tabela orders.

2o para rodar o projeto, basta executar o seguinte comando:
```
go run main.go
```

Serviços:
Starting web server on port :8000
Starting gRPC server on port 50051
Starting GraphQL server on port 8080


teste grpc
evans -r repl