package resp

import (
	"bufio"
	"bytes"
	"errors"
	"strconv"
)

// wraps a bufio.Reader to decode incoming RESP data
type Parser struct {
	reader *bufio.Reader
}

// this is the constructor of Parser
func NewParser(reader *bufio.Reader) *Parser {
	return &Parser{
		reader: reader,
	}
}

// readLine reads from the network until it hits a newline, then strips the \r\n
func (p *Parser) readLine() ([]byte, error) {
	line, err := p.reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}

	// since RESP uses \r\n to end every line, we strip it off
	cleanLine := bytes.TrimSuffix(line, []byte{'\r', '\n'})
	return cleanLine, nil
}

// Parse extracts a full RESP array into a []string
func (p *Parser) Parse() ([]string, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, err
	}
	if len(line) == 0 {
		return nil, errors.New("empty RESP payload")
	}

	// clients talking to a Redis server always wrap their commands in an Array (indicated by *) --> *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n
	if line[0] != '*' {
		return nil, errors.New("expected RESP array prefix")
	}
	// the rest of the line after * is the number of elements in the array
	numArgs, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return nil, err
	}

	// pre allocate slice to exact size for high performance
	args := make([]string, numArgs)

	for i := 0; i < numArgs; i++ {
		// read the bulk string prefix -> $
		strLine, err := p.readLine()
		if err != nil {
			return nil, err
		}
		if len(strLine) == 0 || strLine[0] != '$' {
			return nil, errors.New("expected RESP bulk string prefix")
		}
		// read the string data
		dataLine, err := p.readLine()
		if err != nil {
			return nil, err
		}
		//convert bytes to string and store
		args[i] = string(dataLine)
	}
	return args, nil
}
