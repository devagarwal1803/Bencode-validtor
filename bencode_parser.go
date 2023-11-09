package main

import (
	"errors"
	"fmt"
	"strconv"
)

func parseBencoding(data string) (interface{}, error) {
	if len(data) < 1 {
		return nil, errors.New("Empty input data")
	}
	return parseValue(data)
}

func parseValue(data string) (interface{}, error) {
	switch data[0] {
	case 'i':
		return parseInteger(data[1:])
	case 'l':
		return parseList(data[1:])
	case 'd':
		return parseDictionary(data[1:])
	default:
		return parseString(data)
	}
}

func parseInteger(data string) (int, error) {
	endIndex := 0
	for i := 0; i < len(data); i++ {
		if data[i] == 'e' {
			endIndex = i
			break
		}
	}
	if endIndex == 0 {
		return 0, errors.New("Invalid integer format")
	}
	intValue, err := strconv.Atoi(data[:endIndex])
	if err != nil {
		return 0, err
	}
	return intValue, nil
}

func parseList(data string) ([]interface{}, error) {
	list := make([]interface{}, 0)
	for data[0] != 'e' {
		value, err := parseValue(data)
		if err != nil {
			return nil, err
		}
		list = append(list, value)
		data = data[len(encodeValue(value)):] // Skip the length of the encoded value
	}
	return list, nil
}

func parseDictionary(data string) (map[string]interface{}, error) {
	dict := make(map[string]interface{})
	for len(data) > 0 && data[0] != 'e' {
		key, err := parseString(data)
		if err != nil {
			return nil, err
		}
		data = data[len(encodeString(key)):] // Skip the length of the encoded key
		value, err := parseValue(data)
		if err != nil {
			return nil, err
		}
		dict[key] = value
		data = data[len(encodeValue(value)):] // Skip the length of the encoded value
	}
	return dict, nil
}

func parseString(data string) (string, error) {
	colonIndex := 0
	for i := 0; i < len(data); i++ {
		if data[i] == ':' {
			colonIndex = i
			break
		}
	}
	if colonIndex == 0 {
		return "", errors.New("Invalid string format")
	}
	length, err := strconv.Atoi(data[:colonIndex])
	if err != nil {
		return "", err
	}
	startIndex := colonIndex + 1
	endIndex := startIndex + length
	if endIndex > len(data) {
		return "", errors.New("Invalid string length")
	}
	return data[startIndex:endIndex], nil
}

func encodeValue(value interface{}) string {
	switch v := value.(type) {
	case int:
		return encodeInt(v)
	case string:
		return encodeString(v)
	case []interface{}:
		return encodeList(v)
	case map[string]interface{}:
		return encodeDictionary(v)
	default:
		return ""
	}
}

func encodeInt(value int) string {
	return fmt.Sprintf("i%de", value)
}

func encodeString(value string) string {
	return fmt.Sprintf("%d:%s", len(value), value)
}

func encodeList(value []interface{}) string {
	encoded := "l"
	for _, item := range value {
		encoded += encodeValue(item)
	}
	encoded += "e"
	return encoded
}

func encodeDictionary(value map[string]interface{}) string {
	encoded := "d"
	for key, item := range value {
		encoded += encodeString(key) + encodeValue(item)
	}
	encoded += "e"
	return encoded
}

func main() {
	bencodedData := "d8:announce31:http://tracker.example.com:80806:length4:name11:example.txt5:i123ee"
	result, err := parseBencoding(bencodedData)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println(result)
	}
}
