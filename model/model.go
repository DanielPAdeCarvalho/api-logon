package model

type User struct {
	Nome           string `json:"nome"`
	CPF            string `json:"cpf"`
	NomeCompleto   string `json:"nome-completo"`
	DataNascimento string `json:"data-nascimento"`
	Senha          string `json:"senha"`
	DataCriacao    string `json:"data-criacao"`
}

type Client struct {
	Nome        string `json:"nome"`
	Sobrenome   string `json:"sobrenome"`
	Email       string `json:"email"`
	Senha       string `json:"senha"`
	DataCriacao string `json:"data-criacao"`
}
