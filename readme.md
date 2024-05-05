
1o Ao rodar o programa main.go, o banco de dados será criado automaticamente pelo docker compose. A tabela orders é verificada e criada ao executar o main.go, caso não exista.
```
CREATE TABLE orders (id varchar(255) NOT NULL, price float NOT NULL, tax float NOT NULL, final_price float NOT NULL, PRIMARY KEY (id))
```

2o para rodar o projeto, basta executar o seguinte comando:
```
go run main.go
```

Serviços que irá servir e suas respectivas portas:
Starting web server on port :8000
Starting gRPC server on port 50051
Starting GraphQL server on port 8080

3o Para testar o serviço REST, pode-se usar o curl, que é uma ferramenta de linha de comando para transferir dados com URL ou os seguintes arquivo .http.
```
/golangtarefa3/api/create_order.http   
/golangtarefa3/api/list_order.http
/golangtarefa3/api/list_orders.http
```
4o Para testar o serviço gRPC, pode-se usar o evans, que é uma ferramenta de linha de comando para interagir com serviços gRPC.
```
evans -r repl
```
5o Para testar o serviço graphql, pode-se usar o insomnia, que é uma ferramenta de linha de comando para interagir com serviços graphql.
criar uma ordem:
```
mutation{
  createOrder(input:{id:"123",Price:12.5,Tax:2.0}){id}
}
```
listar uma ordem:
```
getOrder(id:"123"){id,Price,Tax,FinalPrice}
```
listar todas as ordens
```
listOrders(){id,Price,Tax,FinalPrice}
```