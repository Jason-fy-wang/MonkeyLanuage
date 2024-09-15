package parser

const (
	_ int = iota
	LOWEST
	EQUALS      //=
	LESSGREATER // <  or >
	SUM         // + -
	PRODUCT     // 5*5 ,  10/2
	PREFIX      // -X or !X
	CALL        //add(5,5)
)
