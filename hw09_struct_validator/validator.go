package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var builder strings.Builder
	builder.Grow(len(v))
	for _, err := range v {
		builder.WriteString(fmt.Sprintf("Field: %s, error: %s\n", err.Field, err.Err))
	}
	return builder.String()
}

func CreateValidator(validatorName, fieldType string) (Validator, error) {
	if fieldType == "string" {
		switch validatorName {
		case "len":
			return &StringLengthValidator{}, nil
		case "regexp":
			return &RegexpValidator{}, nil
		case "in":
			return &StringInValidator{}, nil
		}
	} else if fieldType == "int" {
		switch validatorName {
		case "min":
			return &MinValidator{}, nil
		case "max":
			return &MaxValidator{}, nil
		case "in":
			return &IntInValidator{}, nil
		}
	}
	return nil, fmt.Errorf("unknown validator %s for fieldType %s", validatorName, fieldType)
}

func getValidators(validateTagString, fieldType string) ([]Validator, error) {
	validatorDefinitions := strings.Split(validateTagString, "|")
	validatorsList := make([]Validator, 0, len(validatorDefinitions))

	for _, validatorDefinition := range validatorDefinitions {
		splittedValidatorDefinition := strings.SplitN(validatorDefinition, ":", 2)
		if len(splittedValidatorDefinition) != 2 {
			return nil, errors.New("wrong validator definition")
		}

		validatorName := splittedValidatorDefinition[0]
		validatorValue := splittedValidatorDefinition[1]
		validator, err := CreateValidator(validatorName, fieldType)
		if err != nil {
			return nil, err
		}

		err = validator.Init(validatorValue)
		if err != nil {
			return nil, fmt.Errorf("validor init error: %w", err)
		}

		validatorsList = append(validatorsList, validator)
	}
	return validatorsList, nil
}

func processValidators(validatorsList []Validator, value interface{}, fieldName string) ValidationErrors {
	validationErrors := make(ValidationErrors, 0)
	for _, validator := range validatorsList {
		err := validator.Validate(value)
		if err != nil {
			validationErrors = append(
				validationErrors,
				ValidationError{
					Field: fieldName,
					Err:   err,
				})
		}
	}
	return validationErrors
}

func Validate(v interface{}) error {
	fmt.Println(reflect.TypeOf(v))
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Struct {
		return errors.New("unsupported value to validate")
	}
	valueType := value.Type()
	allValidationErrorList := make(ValidationErrors, 0)

	for i := 0; i < value.NumField(); i++ {
		fieldOfType := valueType.Field(i)
		field := value.Field(i)
		if !field.CanInterface() {
			continue
		}
		validateTag := fieldOfType.Tag.Get("validate")
		if validateTag == "" {
			continue
		}
		if field.Kind() == reflect.Slice {
			validatorsListForItem, err := getValidators(validateTag, field.Type().Elem().String())
			if err != nil {
				return err
			}
			for i := 0; i < field.Len(); i++ {
				validationErrorsForItem := processValidators(
					validatorsListForItem,
					field.Index(i).Interface(),
					fieldOfType.Name)
				allValidationErrorList = append(allValidationErrorList, validationErrorsForItem...)
			}
		} else {
			validatorsList, err := getValidators(validateTag, field.Kind().String())
			if err != nil {
				return err
			}
			validationErrorsForValue := processValidators(
				validatorsList,
				field.Interface(),
				fieldOfType.Name)
			allValidationErrorList = append(allValidationErrorList, validationErrorsForValue...)
		}
	}

	if len(allValidationErrorList) == 0 {
		return nil
	}
	return allValidationErrorList
}
