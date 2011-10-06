// Functions for reading and writing JSON objects
package percolation

import (
	"io/ioutil"
	"json"
	"os"
	"reflect"
	"strings"
)

// Copy values from the JSON object at filePath into object.
func CopyFromFile(filePath string, object interface{}) os.Error {
	jsonObject, err := ReadJSONFile(filePath)
	if err != nil {
		return err
	}
	CopyValues(jsonObject, object)
	return nil
}

// Copy values from the JSON string given into object.
func CopyFromString(jsonData string, object interface{}) os.Error {
	jsonObject, err := ReadJSONString(jsonData)
	if err != nil {
		return err
	}
	CopyValues(jsonObject, object)
	return nil
}

// Get a JSON object from the byte slice given.
func ReadJSONBytes(jsonData []byte) (*map[string]interface{}, os.Error) {
	jsonObject := make(map[string]interface{})
	err := json.Unmarshal(jsonData, &jsonObject)
	if err != nil {
		return nil, err
	}
	return &jsonObject, nil
}

// Get a JSON object from the file given.
func ReadJSONFile(filePath string) (*map[string]interface{}, os.Error) {
	fileContents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return ReadJSONBytes(fileContents)
}

// Get a JSON object from the string given.
func ReadJSONString(jsonData string) (*map[string]interface{}, os.Error) {
	jsonBytes, err := StringToBytes(jsonData)
	if err != nil {
		return nil, err
	}
	return ReadJSONBytes(jsonBytes)
}

// Write Marshal-able object to a new file at filePath
func WriteJSONFile(object interface{}, filePath string) os.Error {
	marshalled, err := json.Marshal(object)
	if err != nil {
		return err
	}
	jsonFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	if _, err := jsonFile.Write(marshalled); err != nil {
		return err
	}
	if err := jsonFile.Close(); err != nil {
		return err
	}
	return nil
}

// Look at each key in jsonObject and copy thats key's value into the
// corresponding field in object.
func CopyValues(jsonObject *map[string]interface{}, object interface{}) {
	// dereference the object pointer
	objectValue := reflect.Indirect(reflect.ValueOf(object))
	// iterate over all fields in the JSON object
	for key, value := range *jsonObject {
		// get a reference to the field in object
		field := objectValue.FieldByName(key)
		if !field.CanSet() {
			// Can't set, probably because this field doesn't
			// exist in object.  Skip it silently.
			continue
		}
		// recognize some numeric types which aren't available in JSON
		// (can extend this list)
		fieldType := field.Type().Name()
		if fieldType == "int" {
			value = int(value.(float64))
		} else if fieldType == "uint" {
			value = uint(value.(float64))
		}
		// set the field in object
		field.Set(reflect.ValueOf(value))
	}
}

// Convert string to byte slice
func StringToBytes(str string) ([]byte, os.Error) {
	reader := strings.NewReader(str)
	bytes := make([]byte, len(str))
	for seen := 0; seen < len(str); {
		n, err := reader.Read(bytes)
		if err != nil {
			return nil, err
		}
		seen += n
	}
	return bytes, nil
}
