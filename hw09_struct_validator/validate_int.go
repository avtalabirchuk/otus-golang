package hw09structvalidator

import (
	"fmt"
	"strconv"
	"strings"
)

// int validator.
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
