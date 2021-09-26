package validators

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// String validator
type StringLengthValidator struct {
	requiredLength int
}

func (v *StringLengthValidator) Init(validatorValue string) error {
	requiredLength, err := strconv.Atoi(validatorValue)
	if err != nil {
		return fmt.Errorf("unexpected error: %w", err)
	}
	v.requiredLength = requiredLength

	return nil
}

func (v StringLengthValidator) Validate(valueToValidate interface{}) error {
	stringToValidate := fmt.Sprintf("%v", valueToValidate)

	if len(stringToValidate) != v.requiredLength {
		return fmt.Errorf(
			"string length is %d, but length %d is required",
			len(stringToValidate),
			v.requiredLength)
	}
	return nil
}

type RegexpValidator struct {
	regexpPattern *regexp.Regexp
}

func (v *RegexpValidator) Init(validatorValue string) error {
	v.regexpPattern = regexp.MustCompile(validatorValue)
	return nil
}

func (v RegexpValidator) Validate(valueToValidate interface{}) error {
	stringToValidate := fmt.Sprintf("%v", valueToValidate)
	matched := v.regexpPattern.MatchString(stringToValidate)
	if !matched {
		return fmt.Errorf("string %s doesn't fit pattern %s", stringToValidate, v.regexpPattern)
	}
	return nil
}

type StringInValidator struct {
	allowedValues map[string]struct{}
}

func (v *StringInValidator) Init(validatorValue string) error {
	allowedValuesList := strings.Split(validatorValue, ",")
	v.allowedValues = make(map[string]struct{}, len(allowedValuesList))
	for _, value := range allowedValuesList {
		v.allowedValues[value] = struct{}{}
	}
	return nil
}

func (v StringInValidator) Validate(valueToValidate interface{}) error {
	stringToValidate := fmt.Sprintf("%v", valueToValidate)
	_, valueIsAllowed := v.allowedValues[stringToValidate]
	if !valueIsAllowed {
		return fmt.Errorf("string %s doesn't fit allowed set %v", stringToValidate, v.allowedValues)
	}
	return nil
}

// int validator
type MinValidator struct {
	minValue int
}

func (v *MinValidator) Init(validatorValue string) error {
	minValue, err := strconv.Atoi(validatorValue)
	if err != nil {
		return fmt.Errorf("unexpected error: %w", err)
	}
	v.minValue = minValue

	return nil
}

func (v MinValidator) Validate(valueToValidate interface{}) error {
	intToValidate, ok := valueToValidate.(int)
	if !ok {
		return fmt.Errorf("unexpected value %v", valueToValidate)
	}

	if intToValidate < v.minValue {
		return fmt.Errorf(
			"int %d must be greater or equal than %d",
			valueToValidate,
			v.minValue)
	}
	return nil
}

type MaxValidator struct {
	maxValue int
}

func (v *MaxValidator) Init(validatorValue string) error {
	minValue, err := strconv.Atoi(validatorValue)
	if err != nil {
		return fmt.Errorf("unexpected error: %w", err)
	}
	v.maxValue = minValue

	return nil
}

func (v MaxValidator) Validate(valueToValidate interface{}) error {
	intToValidate, ok := valueToValidate.(int)
	if !ok {
		return fmt.Errorf("unexpected value %v", valueToValidate)
	}

	if intToValidate > v.maxValue {
		return fmt.Errorf(
			"int %d must be equal to or lesser than %d",
			valueToValidate,
			v.maxValue)
	}
	return nil
}

type IntInValidator struct {
	allowedValues map[int]struct{}
}

func (v *IntInValidator) Init(validatorValue string) error {
	allowedValuesList := strings.Split(validatorValue, ",")
	v.allowedValues = make(map[int]struct{}, len(allowedValuesList))
	for _, value := range allowedValuesList {
		intValue, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("unexpected error: %w", err)
		}
		v.allowedValues[intValue] = struct{}{}
	}
	return nil
}

func (v IntInValidator) Validate(valueToValidate interface{}) error {
	intToValidate, ok := valueToValidate.(int)
	if !ok {
		return fmt.Errorf("unexpected value %v", valueToValidate)
	}

	_, valueIsAllowed := v.allowedValues[intToValidate]
	if !valueIsAllowed {
		return fmt.Errorf("int %d doesn't fit allowed set %v", intToValidate, v.allowedValues)
	}
	return nil
}

// interface validator
type Validator interface {
	Init(validatorValue string) error
	Validate(valueToValidate interface{}) error
}
