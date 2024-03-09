package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func main() {
	// Declaration
	var input map[string]interface{}

	// Reading file
	data, err := os.ReadFile("input.json")
	if err != nil {
		fmt.Println(err)
	}

	// JSON conversion
	json.Unmarshal(data, &input)

	// Performing Transition Operation
	output := constructMap(input)

	// Priting output in human readable format
	jsonString, err := json.MarshalIndent(output, "", "    ")
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Print the formatted JSON string
	fmt.Println(string(jsonString))
}

func constructMap(input map[string]interface{}) map[string]interface{} {
	// Reading Map data type

	output := make(map[string]interface{})

	for pk, pv := range input {
		if pk == "" {
			continue
		}
		for k, v := range pv.(map[string]interface{}) {
			if k == "" || v == "" {
				continue
			}

			switch k {
			case "N":
				status, rtnData := validateType("number", v)

				if status {
					output[pk] = rtnData
				}

			case "S":
				status, rtnData := validateType("string", v)

				if status {
					output[pk] = rtnData
				}
			case "M":
				data, ok := v.(map[string]interface{})

				if !ok {
					output[pk] = v
				} else {
					output[pk] = constructMap(data)
				}

			case "L":
				rtnData := constructList(v)

				if len(rtnData) > 0 {
					output[pk] = rtnData
				}
			case "BOOL":
				status, rtnData := validateType("bool", v)

				if status {
					output[pk] = rtnData
				}

			case "NULL":
				status, rtnData := validateType("null", v)

				if status {
					output[pk] = rtnData
				}

			default:
				output[pk] = v
			}
		}

	}

	return output

}

func constructList(data interface{}) []interface{} {
	// Reading list of array data type

	var out []interface{}

	switch data.(type) {
	case []interface{}:

		for _, pv := range data.([]interface{}) {

			if reflect.TypeOf(pv).Kind() == reflect.Map {
				for k, v := range pv.(map[string]interface{}) {
					if k == "" || v == "" {
						continue
					}

					switch k {
					case "S":
						status, rtnData := validateType("string", v)

						if status {
							out = append(out, rtnData)
						}
					case "N":
						status, rtnData := validateType("number", v)

						if status {
							out = append(out, rtnData)
						}
					case "BOOL":
						status, rtnData := validateType("bool", v)

						if status {
							out = append(out, rtnData)
						}
					case "NULL":

						status, rtnData := validateType("null", v)

						if status {
							out = append(out, rtnData)
						}
					}
				}
			}
		}

	}

	return out
}

func validateType(dataType string, val interface{}) (bool, interface{}) {

	// Utility function to validate and convert the type

	switch dataType {
	case "string":
		cnvData, ok := val.(string)
		if !ok {
			return false, val
		}

		return true, strings.Trim(cnvData, " ")

	case "number":
		strValue := val.(string)

		// Convert the string to a float64
		cnvData, err := strconv.ParseFloat(strValue, 64)

		if err != nil {
			return false, val
		}

		return true, cnvData
	case "bool":
		if strings.ToLower(val.(string)) == "true" || strings.ToLower(val.(string)) == "t" || strings.ToLower(val.(string)) == "1" {
			return true, true
		} else if strings.ToLower(val.(string)) == "false" || strings.ToLower(val.(string)) == "f" || strings.ToLower(val.(string)) == "0" {
			return true, false
		}

		return false, val

	case "null":
		return val.(string) == "", val.(string)
	default:
		return false, val
	}

}
