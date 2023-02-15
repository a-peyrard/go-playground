package main

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"reflect"
	"strings"
)

type NestedStruct struct {
	NestedField1 string
	NestedField2 int
}

type MyStruct struct {
	Field1       string
	Field2       int
	NestedStruct *NestedStruct
}

func BindEnvs(iface interface{}, viperI *viper.Viper, parts ...string) {
	ifv := reflect.ValueOf(iface)
	ift := reflect.TypeOf(iface)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			tv = t.Name
		}
		switch v.Kind() {
		case reflect.Struct:
			BindEnvs(v.Interface(), viperI, append(parts, tv)...)
		case reflect.Pointer:
			BindEnvs(reflect.New(t.Type.Elem()).Elem().Interface(), viperI, append(parts, tv)...)
		default:
			_ = viperI.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
}

type GlobalConfig struct {
	App *AppConfig
	Db  *DatabaseConfig
	Foo string `mapstructure:"MY_FOO"`
}

type AppConfig struct {
	Env  string
	Port int
}

type DatabaseConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

func main() {
	v := viper.New()
	v.AutomaticEnv()
	v.SetEnvPrefix("ES")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	BindEnvs(GlobalConfig{}, v)

	var config GlobalConfig
	if err := v.Unmarshal(&config); err != nil {
		log.Fatalf("%+v", err)
	}

	var dbConfig2 DatabaseConfig
	if err := v.UnmarshalKey("db", &dbConfig2); err != nil {
		log.Fatalf("%+v", err)
	}

	appConfig := config.App
	dbConfig := config.Db
	// err := v.UnmarshalKey("APP", &appConfig)
	//
	//	if err != nil {
	//		log.Fatalf("%+v", err)
	//	}
	//
	// var dbConfig DatabaseConfig
	// err = viper.UnmarshalKey("DB", &dbConfig)
	//
	//	if err != nil {
	//		log.Fatalf("%+v", err)
	//	}
	fmt.Printf("config = %+v\n", config)
	fmt.Printf("appConfig = %+v\n", appConfig)
	fmt.Printf("dbConfig = %+v\n", dbConfig)
	fmt.Printf("dbConfig2 = %+v\n", dbConfig2)
	fmt.Printf("allKeys = %+v\n", v.AllKeys())
}
