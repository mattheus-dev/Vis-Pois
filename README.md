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
4. A primeira linha (cabeçalho) é utilizada para mapear as colunas
5. Cada linha subsequente é dividida em campos usando `strings.Split` (suporta separadores `,` ou `;`)
6. Os dados são convertidos para a estrutura `Record`
7. Todas as estruturas são processadas e retornadas como uma tabela formatada

O processo inclui tratamento de erros para:
- Arquivos vazios
- Linhas mal formatadas
- Erros ao abrir ou ler o arquivo
- Conversão de tipos de dados

## Estrutura da entidade Record

```go
type Record struct {
	ID                 string  `json:"id"`
	Name               string  `json:"name"`
	Price              float64 `json:"price"`
	Stock              int     `json:"stock"`
	Category           string  `json:"category"`
	Subcategory        string  `json:"subcategory"`
	Brand              string  `json:"brand"`
	Description        string  `json:"description"`
	ImageURL           string  `json:"image_url"`
	Weight             float64 `json:"weight"`
	Dimensions         string  `json:"dimensions"`
	Color              string  `json:"color"`
	Material           string  `json:"material"`
	CountryOfOrigin    string  `json:"country_of_origin"`
	Manufacturer       string  `json:"manufacturer"`
	SKU                string  `json:"sku"`
	Barcode            string  `json:"barcode"`
	TaxRate            float64 `json:"tax_rate"`
	DiscountPercentage float64 `json:"discount_percentage"`
	Rating             float64 `json:"rating"`
	ReviewCount        int     `json:"review_count"`
	IsActive           bool    `json:"is_active"`
	CreatedAt          string  `json:"created_at"`
	UpdatedAt          string  `json:"updated_at"`
}
```

## Exemplo de CSV

Aqui está um exemplo de arquivo CSV simples que pode ser usado para teste:

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
```

## Exemplo de CSV Complexo

A aplicação também suporta processamento de CSVs com estrutura mais complexa:

```
ID,Nome,Descricao,Marca,Categoria,Subcategoria,Preco,Estoque,Peso,Dimensoes,Cor,Material,PaisOrigem,Fabricante,SKU,CodigoBarras,TaxaImposto,PercentualDesconto,Avaliacao,NumeroAvaliacoes,Ativo,DataCriacao,DataAtualizacao,ImagemURL
TEC-MEC-001,Teclado Mecânico,Teclado mecânico com switches Cherry MX,HyperX,Periféricos,Teclados,299.90,45,0.85,44x14x4cm,Preto,Plástico e Alumínio,China,Kingston,TEC-MEC-001,7891234567890,12.5,5,4.7,128,true,2023-06-15,2023-08-20,https://example.com/images/teclado.jpg
```

## Como testar o endpoint

### Usando Curl

```bash
curl -X POST -F "file=@caminho/para/seu/arquivo.csv" http://localhost:8080/leitura/teste
```

### Usando Postman

1. Abra o Postman
2. Crie uma requisição POST para `http://localhost:8080/leitura/teste`
3. Na aba "Body", selecione "form-data"
4. Adicione uma chave chamada "file" e selecione o tipo "File"
5. Selecione o arquivo CSV para upload
6. Clique em "Send" para enviar a requisição

## Executando a Aplicação

Para executar a aplicação, use o seguinte comando na raiz do projeto:

```bash
go run main.go
```

A aplicação estará disponível em `http://localhost:8080`.
