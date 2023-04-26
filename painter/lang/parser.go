package lang

import (
	"bufio"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	var res []painter.Operation
	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		commandLine := scanner.Text()
		op, err := parseCommand(commandLine)

		if err != nil {
			return res, err
		}

		if op != nil { // TODO: remove condition in future
			res = append(res, op)
		}
	}

	return res, nil
}

func parseCommand(commandLine string) (painter.Operation, error) {
	parsedCommand := strings.Fields(commandLine)
	commandName := parsedCommand[0]
	params := parsedCommand[1:]

	switch commandName {
	case "white":
		return painter.OperationFunc(painter.WhiteFill), nil
	case "green":
		return painter.OperationFunc(painter.GreenFill), nil
	case "update":
		return painter.UpdateOp, nil
	case "bgrect":
		parsedParams, err := parseParams(params, 4)
		if err != nil {
			return nil, err
		}
		return nil, nil // TODO: create struct
	case "figure":
		parsedParams, err := parseParams(params, 2)
		if err != nil {
			return nil, err
		}
		return nil, nil // TODO: create struct
	case "move":
		parsedParams, err := parseParams(params, 2)
		if err != nil {
			return nil, err
		}
		return nil, nil // TODO: create struct
	case "reset":
		return nil, nil // TODO: add action
	default:
		return nil, errors.New("unknown command")
	}
}

func parseParams(params []string, length int) ([]float32, error) {
	var res []float32

	if len(params) != length {
		return nil, errors.New("invalid params count")
	}

	for _, item := range params {
		floatNum, err := strconv.ParseFloat(item, 32)
		if err != nil {
			return nil, errors.New("invalid params")
		}

		if floatNum < 0 || floatNum > 1 {
			return nil, errors.New("invalid coordinates")
		}

		res = append(res, float32(floatNum))
	}

	return res, nil
}
