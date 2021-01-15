package wsproto

import "encoding/json"

// RawCmdArg : not parsed args for command
type RawCmdArg []byte

// NextCommand : parse buff stream, return command type and next position
func NextCommand(buff []byte) (string, RawCmdArg, []byte, bool) {
	n := len(buff)
	isStringLiteral := false
	isEscaping := false
	cmdLen := 0
	for i := 0; i < n; i++ {
		if isStringLiteral && isEscaping {
			isEscaping = false
		} else if isStringLiteral {
			switch buff[i] {
			case '"':
				isStringLiteral = !isStringLiteral
			case '\\':
				isEscaping = true
			}
		} else if !isEscaping && buff[i] == ';' {
			cmdLen = i + 1
			break
		}
	}
	cmdType, rawCmdArg := parseCommand(buff[:cmdLen])
	return cmdType, rawCmdArg, buff[cmdLen:], cmdLen > 0
}

// BindCmdArgs : bind args for command
func (raw RawCmdArg) BindCmdArgs(args interface{}) error {
	return json.Unmarshal(raw, args)
}

func parseCommand(buff []byte) (string, RawCmdArg) {
	n := len(buff)
	i := 0
	for ; i < n; i++ {
		if buff[i] == ':' {
			break
		}
	}
	return string(buff[:i]), RawCmdArg(buff[i+1:])
}
