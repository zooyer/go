package fastdfs

import (
	"io"
	"strconv"
	"strings"
	"reflect"
	"fmt"
	"io/ioutil"
	"bytes"
	"bufio"
	"errors"
	"golang.org/x/text/encoding/simplifiedchinese"
)

type NameValuePair struct {
	name,value string
}

func NewNameValuePair(name,value string) *NameValuePair {
	return &NameValuePair{
		name:name,
		value:value,
	}
}

func (p *NameValuePair) GetName() string {
	return p.name
}

func (p *NameValuePair) SetName(name string) {
	p.name = name
}

func (p *NameValuePair) GetValue() string {
	return p.value
}

func (p *NameValuePair) SetValue(value string) {
	p.value = value
}


func LoadFromOsFileSystemOrClasspathAsStream(filePath string) (io.Reader, error) {
	data,err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}


type IniFileReader struct {
	paramTable      Hashtable
	confFilename   string
}

func NewIniFileReader(confFilename string) (*IniFileReader, error) {
	var r = new(IniFileReader)
	r.confFilename = confFilename

	if err := r.loadFromFile(confFilename); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *IniFileReader) GetConfFilename() string {
	return r.confFilename
}

func (r *IniFileReader) GetStrValue(name string) string {
	var val = r.paramTable.Get(name)
	if val != nil {
		if v,ok := val.(string); ok {
			return v
		}
	}

	return ""
}

func (r *IniFileReader) GetIntValue(name string, defaultValue int) int {
	var val = r.GetStrValue(name)
	if val == "" {
		return defaultValue
	}

	num,_ := strconv.Atoi(val)

	return num
}

func (r *IniFileReader) GetBoolValue(name string, defaultValue bool) bool {
	var val = r.GetStrValue(name)
	if val == "" {
		return defaultValue
	}

	val = strings.ToLower(val)

	if val == "true" || val == "yes" || val == "on" || val == "1" {
		return true
	}

	return false
}

func (r *IniFileReader) GetValues(name string) []string {
	var val = r.paramTable.Get(name)
	if val == nil {
		return nil
	}
	var values []string

	if str,ok := val.(string); ok {
		values = make([]string, 1)
		values[0] = str
		return values
	}

	if reflect.TypeOf(val).Kind() == reflect.Slice {
		v := reflect.ValueOf(val)
		values = make([]string, v.Len())
		for i := 0; i < v.Len(); i++ {
			values[i] = fmt.Sprintf("%v", v.Index(i).Interface())
		}
		return values
	}

	return nil
}

func (r *IniFileReader) loadFromFile(confFilePath string) error {
	reader,err := LoadFromOsFileSystemOrClasspathAsStream(confFilePath)
	if err != nil {
		return err
	}

	return r.readToParamTable(reader)
}

func (r *IniFileReader) readToParamTable(in io.Reader) error {
	r.paramTable = NewHashtable()
	if in == nil {
		return nil
	}
	var line string
	var parts []string
	var name string
	var value string
	var inter interface{}
	var valueList []interface{}

	var reader = bufio.NewReader(in)
	for l,_,err := reader.ReadLine(); err == nil; l,_,err = reader.ReadLine() {
		line = strings.TrimSpace(string(l))
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		parts = strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		name = strings.TrimSpace(parts[0])
		value = strings.TrimSpace(parts[1])
		inter = r.paramTable.Get(name)
		if inter == nil {
			r.paramTable.Put(name, value)
		} else if str,ok := inter.(string); ok {
			valueList = nil
			valueList = append(valueList, str)
			valueList = append(valueList, value)
			r.paramTable.Put(name, valueList)
		} else {
			if list,ok := inter.([]interface{}); ok {
				valueList = list
				valueList = append(valueList, value)
			} else {
				return errors.New("unknown type")
			}
		}
	}

	return nil
}

const (
	UTF8    = "UTF-8"
	GB18030 = "GB18030"
)

// return utf8 string
func ConvertByteToString(bytes []byte, charset string) (string, error) {
	var str string
	switch charset {
	case GB18030:
		var decodeBytes,err = simplifiedchinese.GB18030.NewDecoder().Bytes(bytes)
		if err != nil {
			return "", err
		}
		str= string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(bytes)
	}

	return str, nil
}

// return utf8
func ConvertBytesToUTF8(bytes []byte, charset string) ([]byte, error) {
	switch charset {
	case GB18030:
		var decodeBytes,err = simplifiedchinese.GB18030.NewDecoder().Bytes(bytes)
		if err != nil {
			return nil, err
		}
		return decodeBytes, nil
	case UTF8:
		return bytes, nil
	}

	return nil, fmt.Errorf("not support charset %s", charset)
}

// bytes is utf8
func ConvertUTF8ToBytes(bytes []byte, charset string) ([]byte, error) {
	switch charset {
	case GB18030:
		var decodeBytes,err = simplifiedchinese.GB18030.NewEncoder().Bytes(bytes)
		if err != nil {
			return nil, err
		}
		return decodeBytes, nil
	case UTF8:
		return bytes, nil
	}

	return nil, fmt.Errorf("not support charset %s", charset)
}
