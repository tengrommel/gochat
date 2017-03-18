package Hydraconfigurator

import (
	"bufio"
	"errors"
	"os"
	"reflect"
	"strconv"
	"fmt"
	"strings"
)

type ConfigFields map[string]reflect.Value

func (f ConfigFields)Add(name, value, t string) error {
	switch t {
	case "STRING":
		f[name] = reflect.ValueOf(value)
	case "INTEGER":
		i, err := strconv.Atoi(value)
		if err != nil{
			return err
		}
		f[name] = reflect.ValueOf(i)
	case "FLOAT":
		fl, err := strconv.ParseFloat(value, 64)
		if err != nil{
			return err
		}
		f[name] = reflect.ValueOf(fl)
	case "BOOL":
		b, err := strconv.ParseBool(value)
		if err !=nil{
			return err
		}
		f[name] = reflect.ValueOf(b)
	}
	return nil
}

func MarshalCustomConfig(v reflect.Value, filename string) error {

	defer func() {
		if r:=recover();r!=nil{
			fmt.Println("Panic occured", r)
		}
	}()

	if !v.CanSet(){
		return errors.New("Value passwd not settable")
	}

	file, err := os.Open(filename)
	if err != nil{
		return err
	}

	defer file.Close()
	fields := make(ConfigFields)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("Processing line", line)
		if strings.Count(line, "|") != 1 || strings.Count(line,";") != 1{
			continue
		}
		args := strings.Split(line, "|")
		valuetype := strings.Split(args[1],";")
		name, value, vtype := strings.TrimSpace(args[0]), strings.TrimSpace(valuetype[0]), strings.ToUpper(strings.TrimSpace(valuetype[1]))
		fields.Add(name, value, vtype)
	}

	if err := scanner.Err(); err != nil{
		return err
	}

	vt := v.Type()
	for i:=0;i<v.NumField();i++ {
		fieldType := vt.Field(i)
		fieldValue := v.Field(i)

		name := fieldType.Tag.Get("name")
		if name == ""{
			name =  fieldType.Tag.Get("name")
		}

		if v, ok := fields[name];ok{
			fieldValue.Set(v)
		}
	}
	return nil
}