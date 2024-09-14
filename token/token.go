package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT = "IDENT"
	INT   = "INT"

	ASSIGN = "="
	PLUS   = "+"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	FUNCTION = "FUNCTION"
	LET      = "LET"
)

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

var keyworkds = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

func LoopupIdentifier(ident string) TokenType {

	if tkType, ok := keyworkds[ident]; ok {
		return tkType
	}

	return IDENT
}
