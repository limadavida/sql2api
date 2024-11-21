# Makefile
.PHONY: setup lint install golint staticcheck

# Diretório do GOPATH
GOPATH := $(shell go env GOPATH)

# Caminho do binário do GOPATH
BINPATH := $(GOPATH)/bin

PROJECT := "sql2api"

# Função para verificar se GOPATH está no PATH
check_path:
	@echo "Verificando se $(BINPATH) está no PATH..."
	@if ! echo $$PATH | grep -q $(BINPATH); then \
		echo "O diretório $(BINPATH) não está no PATH."; \
		echo "Adicione a linha 'export PATH=\$$PATH:$(BINPATH)' ao seu shell."; \
		echo "Depois, execute 'source ~/.bash_profile' ou 'source ~/.bashrc' ou 'source ~/.zshrc';"; \
		exit 1; \
	fi

# Instalar as ferramentas de lint
install-lint:
	@echo "Instalando golint e staticcheck..."
	go install golang.org/x/lint/golint@latest
	go install honnef.co/go/tools/cmd/staticcheck@latest

# Executar o golint
golint: check_path
	export PATH=$PATH:$(go env GOPATH)/bin
	@echo "Executando golint..."
	golint ./...

# Executar o staticcheck
staticcheck: check_path
	@echo "Executando staticcheck..."
	staticcheck ./...

# Regra padrão para executar lint
lint: staticcheck


test: check_path
	@echo "Executando testes..."
	go test ./... -v

graph_dummy: 
	@echo "Gerando dependencias em graph"
	go mod graph > graph.dot
	chmod +x ./tograph.sh
	./tograph.sh > fgraph.dot
	dot -Tpng -Gsize=8,8\! fgraph.dot -o fgraph.png

graph:
	go install github.com/loov/goda@latest
	goda graph -cluster -short "github.com/limadavida/$(PROJECT)/..."| dot -Tsvg -o cluster_$(PROJECT)_graph.svg
	goda graph "github.com/limadavida/$(PROJECT)/..." | dot -Tsvg -o pkg_$(PROJECT)_graph.svg

nocache:
	go clean -modcache
	go mod tidy
	go get -u

#	go list -m -u all

