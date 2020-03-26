package dynamique

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unsafe"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/convox/console/pkg/crypt"
	"github.com/pkg/errors"
)

const (
	SortableTime = "20060102.150405.000000000"
)

func Marshal(target interface{}, key string) (map[string]*dynamodb.AttributeValue, error) {
	item := make(map[string]*dynamodb.AttributeValue)

	v := reflect.ValueOf(target)
	t := v.Type()

	typeElem := t.Elem()
	structElem := v.Elem()
	n := typeElem.NumField()

	for i := 0; i < n; i++ {
		typeField := typeElem.Field(i)
		structField := structElem.Field(i)
		if !structField.CanSet() {
			structField = reflect.NewAt(structField.Type(), unsafe.Pointer(structField.UnsafeAddr())).Elem()
		}
		parts := strings.Split(typeField.Tag.Get("dynamo"), ",")
		tag := parts[0]
		tagopts := map[string]bool{}
		for _, s := range parts[1:] {
			tagopts[s] = true
		}
		if tag != "" {
			switch val := structField.Interface().(type) {
			case bool:
				item[tag] = &dynamodb.AttributeValue{BOOL: aws.Bool(val)}
			case int, int64:
				item[tag] = &dynamodb.AttributeValue{N: aws.String(fmt.Sprintf("%d", val))}
			case string:
				if val != "" {
					if tagopts["encrypted"] {
						enc, err := crypt.Encrypt(key, []byte(val))
						if err != nil {
							return nil, errors.WithStack(err)
						}
						val = enc
					}
					item[tag] = &dynamodb.AttributeValue{S: aws.String(val)}
				}
			case time.Time:
				item[tag] = &dynamodb.AttributeValue{S: aws.String(val.Format(SortableTime))}
			default:
				if val != nil {
					data, err := json.Marshal(val)
					if err != nil {
						return nil, errors.WithStack(fmt.Errorf("could not marshal %s into %q as json: %s", typeField.Name, tag, err.Error()))
					}
					if tagopts["encrypted"] {
						enc, err := crypt.Encrypt(key, data)
						if err != nil {
							return nil, errors.WithStack(err)
						}
						data = []byte(enc)
					}
					item[tag] = &dynamodb.AttributeValue{S: aws.String(string(data))}
				}
			}
		}
	}

	return item, nil
}

func Unmarshal(item map[string]*dynamodb.AttributeValue, target interface{}, key string) error {
	v := reflect.ValueOf(target)
	t := v.Type()

	typeElem := t.Elem()
	structElem := v.Elem()
	n := typeElem.NumField()

	for i := 0; i < n; i++ {
		typeField := typeElem.Field(i)
		structField := structElem.Field(i)
		if !structField.CanSet() {
			structField = reflect.NewAt(structField.Type(), unsafe.Pointer(structField.UnsafeAddr())).Elem()
		}
		parts := strings.Split(typeField.Tag.Get("dynamo"), ",")
		tag := parts[0]
		tagopts := map[string]bool{}
		for _, s := range parts[1:] {
			tagopts[s] = true
		}
		if tag != "" { //&& structField.CanSet() {
			switch structField.Interface().(type) {
			case bool:
				s := item[tag]
				if s != nil && s.BOOL != nil {
					structField.SetBool(*s.BOOL)
				}
			case int, int64:
				s := item[tag]

				if s != nil {
					if s.N == nil {
						return errors.WithStack(fmt.Errorf("%q is not a number", tag))
					}
					n, err := strconv.Atoi(*s.N)
					if err != nil {
						return errors.WithStack(err)
					}
					structField.SetInt(int64(n))
				}
			case string:
				str := coalesceS(item[tag], "")
				if str != "" && tagopts["encrypted"] {
					dec, err := crypt.Decrypt(key, str)
					if err != nil {
						return errors.WithStack(err)
					}
					str = string(dec)
				}
				structField.SetString(str)
			case time.Time:
				//NOTE: this always assumes time.Time is serialized as a string formatted as SortableTime
				if item[tag] != nil {
					t, err := time.Parse(SortableTime, *item[tag].S)
					if err != nil {
						return errors.WithStack(err)
					}
					structField.Set(reflect.ValueOf(t))
				}
			default:
				str := coalesceS(item[tag], "")
				if str != "" {
					if tagopts["encrypted"] {
						dec, err := crypt.Decrypt(key, str)
						if err != nil {
							return errors.WithStack(err)
						}
						str = string(dec)
					}
					t := reflect.New(structField.Type()).Interface()
					err := json.Unmarshal(
						[]byte(str),
						&t,
					)
					if err != nil {
						return errors.WithStack(fmt.Errorf("could not unmarshal %q json to %s: %s", tag, typeField.Name, err.Error()))
					}
					if t != nil {
						structField.Set(reflect.Indirect(reflect.ValueOf(t)))
					}
				}
			}
		}
	}

	return nil
}

func coalesceS(s *dynamodb.AttributeValue, def string) string {
	if s != nil && s.S != nil {
		return *s.S
	} else {
		return def
	}
}
