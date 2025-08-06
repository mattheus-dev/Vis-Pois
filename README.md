# CSV Processing Application

Esta aplicação é uma API REST que permite fazer upload e processamento de arquivos CSV seguindo uma arquitetura hexagonal. A aplicação foi projetada para processar arquivos CSV que serão usados na atualização do Elasticsearch.

## Estrutura da Aplicação

A aplicação segue a arquitetura hexagonal (ports and adapters) para separar as responsabilidades:

- **Domain**: Contém as entidades de negócio e portas (interfaces)
- **Application**: Contém a lógica de aplicação e serviços
- **Adapters**: Contém os adaptadores para interfaces externas (HTTP)
- **Infrastructure**: Contém implementações concretas para recursos externos

## Como funciona a leitura de CSV

O processo de leitura de arquivos CSV funciona da seguinte forma:

1. O arquivo é recebido via upload através do endpoint `/leitura/teste`
2. O arquivo é salvo temporariamente no sistema de arquivos local
3. O arquivo é aberto e lido linha por linha usando `bufio.NewScanner`
4. A primeira linha (cabeçalho) é ignorada
5. Cada linha subsequente é dividida em campos usando `strings.Split` (suporta separadores `,` ou `;`)
6. Os dados são convertidos para a estrutura `Record`
7. Todas as estruturas são logadas e retornadas como resposta JSON

O processo inclui tratamento de erros para:
- Arquivos vazios
- Linhas mal formatadas
- Erros ao abrir ou ler o arquivo

## Estrutura da entidade Record

```go
type Record struct {
    ID    string  `json:"id"`
    Name  string  `json:"name"`
    Price float64 `json:"price"`
    Stock int     `json:"stock"`
}
```

## Exemplo de CSV

Aqui está um exemplo de arquivo CSV que pode ser usado para teste:

```
ID;Nome;Preco;Estoque
001;Teclado Mecânico;299.90;45
002;Mouse Gamer;159.90;78
003;Headset;189.90;23
004;Mousepad;59.90;120
```

Ou usando vírgulas como separador:

```
ID,Nome,Preco,Estoque
001,Teclado Mecânico,299.90,45
002,Mouse Gamer,159.90,78
003,Headset,189.90,23
004,Mousepad,59.90,120
```

## Como testar o endpoint

### Usando cURL

```bash
curl -X POST \
  http://localhost:8080/leitura/teste \
  -H 'Content-Type: multipart/form-data' \
  -F 'file=@/caminho/para/seu/arquivo.csv'
```

### Usando Postman

1. Abra o Postman
2. Crie uma nova requisição POST para `http://localhost:8080/leitura/teste`
3. Na aba "Body", selecione "form-data"
4. Adicione um campo chamado "file" e selecione o tipo "File"
5. Clique no botão "Select Files" e escolha seu arquivo CSV
6. Clique em "Send"

## Executando a Aplicação

Para executar a aplicação, use o seguinte comando na raiz do projeto:

```bash
go run main.go
```

A aplicação estará disponível em `http://localhost:8080`.
