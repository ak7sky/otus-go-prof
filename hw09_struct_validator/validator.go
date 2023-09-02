package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

const (
	validateTag = "validate"

	strType    = "string"
	strSlcType = "[]string"
	intType    = "int"
	intSlcType = "[]int"
	structType = "struct"

	lenRestr    = "len"
	regexpRestr = "regexp"
	inRestr     = "in"
	minRestr    = "min"
	maxRestr    = "max"
	nested      = "nested"
)

var (
	// app errors.
	errNotStruct         = errors.New("passed arg must be a go struct type")
	errNotNested         = errors.New("supported validate tag for nested struct is 'nested'")
	errEmptyVldTag       = errors.New("validate tag value must not be empty")
	errFieldType         = errors.New("unsupported field type; supported: string, []string, int, []int")
	errRestrType         = errors.New("unknown restriction for field type")
	errNotIntLenRestr    = errors.New("value of 'len' restriction must be int")
	errNotIntMinRestr    = errors.New("value of 'min' restriction must be int")
	errNotIntMaxRestr    = errors.New("value of 'max' restriction must be int")
	errNotIntInRestr4int = errors.New("all values of 'in' restriction must be integers for int type")

	// vld error.
	errValidation = errors.New("violations")
)

type fieldValidator func(name string, value reflect.Value, restrs []string) error

var fieldValidators map[string]fieldValidator

// to avoid initialization loop fieldValidators -> structValidator -> validateField -> fieldValidators.
func init() {
	fieldValidators = map[string]fieldValidator{
		strType:    strValidator,
		strSlcType: strSlcValidator,
		intType:    intValidator,
		intSlcType: intSlcValidator,
		structType: structValidator,
	}
}

type strRestrMatcher func(val string, restrVal string) (bool, error)

var strRestrMatchers = map[string]strRestrMatcher{
	lenRestr:    strLenRestrMatcher,
	regexpRestr: strRegexpRestrMatcher,
	inRestr:     strInRestrMatcher,
}

type intRestrMatcher func(val int, restrVal string) (bool, error)

var intRestrMatchers = map[string]intRestrMatcher{
	minRestr: intMinRestrMatcher,
	maxRestr: intMaxRestrMatcher,
	inRestr:  intInRestrMatcher,
}

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("field: '%s'; %s", v.Field, v.Err)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	builder := strings.Builder{}
	builder.WriteString("validation errors:\n")

	for _, ve := range v {
		builder.WriteString("\t")
		builder.WriteString(ve.Error())
		builder.WriteString("\n")
	}

	return builder.String()
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)

	if val.Kind() != reflect.Struct {
		return errNotStruct
	}

	return structValidator("", val, []string{})
}

func validateField(fieldMeta reflect.StructField, fieldVal reflect.Value) error {
	validateTagVal, validationRequired := fieldMeta.Tag.Lookup(validateTag)
	if !validationRequired {
		return nil
	}
	if validateTagVal == "" {
		return fmt.Errorf("field '%s': %w", fieldMeta.Name, errEmptyVldTag)
	}

	fieldType := underlyingTypeOf(fieldVal)
	fv, fieldTypeSupported := fieldValidators[fieldType]
	if !fieldTypeSupported {
		return fmt.Errorf("field '%s', type %s: %w", fieldMeta.Name, fieldType, errFieldType)
	}

	fieldRestrs := strings.Split(validateTagVal, "|")
	return fv(fieldMeta.Name, fieldVal, fieldRestrs)
}

func underlyingTypeOf(value reflect.Value) string {
	kind := value.Kind()
	if kind == reflect.Slice {
		if value.CanConvert(reflect.TypeOf([]string{})) {
			return strSlcType
		}
		if value.CanConvert(reflect.TypeOf([]int{})) {
			return intSlcType
		}
	}
	return kind.String()
}

func structValidator(name string, value reflect.Value, restrs []string) error {
	if name != "" && (len(restrs) > 1 || restrs[0] != nested) {
		return fmt.Errorf("field '%s': %w", name, errNotNested)
	}

	var verrors ValidationErrors

	for _, fieldMeta := range reflect.VisibleFields(value.Type()) {
		fieldVal := value.FieldByName(fieldMeta.Name)
		err := validateField(fieldMeta, fieldVal)
		if err == nil {
			continue
		}
		vldErr := &ValidationError{}
		if errors.As(err, vldErr) {
			verrors = append(verrors, *vldErr)
			continue
		}
		return err
	}

	if verrors != nil {
		if name != "" {
			return ValidationError{Field: name, Err: verrors}
		}
		return verrors
	}

	return nil
}

func strValidator(name string, value reflect.Value, restrs []string) error {
	val := value.Convert(reflect.TypeOf("")).Interface().(string)
	var violatedRestrs []string
	for _, restr := range restrs {
		restrTypeAndVal := strings.Split(restr, ":")
		restrType, restrVal := restrTypeAndVal[0], restrTypeAndVal[1]
		restrMatcher, restrTypeSupported := strRestrMatchers[restrType]
		if !restrTypeSupported {
			return fmt.Errorf("field '%s', type %s, restr type '%s': %w", name, strType, restrType, errRestrType)
		}
		isMatched, err := restrMatcher(val, restrVal)
		if err != nil {
			return fmt.Errorf("field '%s': %w", name, err)
		}
		if !isMatched {
			violatedRestrs = append(violatedRestrs, restr)
		}
	}
	if violatedRestrs != nil {
		return ValidationError{
			Field: name,
			Err:   fmt.Errorf("%w: val '%s', %v", errValidation, val, violatedRestrs),
		}
	}

	return nil
}

func strSlcValidator(name string, value reflect.Value, restrs []string) error {
	vals := value.Convert(reflect.TypeOf([]string{})).Interface().([]string)
	var elVldErrors []error

	for idx, val := range vals {
		err := strValidator(name, reflect.ValueOf(val), restrs)
		if err == nil {
			continue
		}

		vldErr := &ValidationError{}
		if errors.As(err, vldErr) {
			elVldErrors = append(elVldErrors, fmt.Errorf("idx %d: %w", idx, vldErr.Err))
			continue
		}
		return err
	}

	if elVldErrors != nil {
		return ValidationError{
			Field: name,
			Err:   fmt.Errorf("%w: %v", errValidation, elVldErrors),
		}
	}

	return nil
}

func intValidator(name string, value reflect.Value, restrs []string) error {
	val := value.Convert(reflect.TypeOf(0)).Interface().(int)
	var violatedRestrs []string
	for _, restr := range restrs {
		restrTypeAndVal := strings.Split(restr, ":")
		restrType, restrVal := restrTypeAndVal[0], restrTypeAndVal[1]
		restrMatcher, restrTypeSupported := intRestrMatchers[restrType]
		if !restrTypeSupported {
			return fmt.Errorf("field '%s', type %s, restr type '%s': %w", name, intType, restrType, errRestrType)
		}
		isMatched, err := restrMatcher(val, restrVal)
		if err != nil {
			return fmt.Errorf("field '%s': %w", name, err)
		}
		if !isMatched {
			violatedRestrs = append(violatedRestrs, restr)
		}
	}
	if violatedRestrs != nil {
		return ValidationError{
			Field: name,
			Err:   fmt.Errorf("%w: val '%d', %v", errValidation, val, violatedRestrs),
		}
	}

	return nil
}

func intSlcValidator(name string, value reflect.Value, restrs []string) error {
	vals := value.Convert(reflect.TypeOf([]int{})).Interface().([]int)
	var elVldErrors []error

	for idx, val := range vals {
		err := intValidator(name, reflect.ValueOf(val), restrs)
		if err == nil {
			continue
		}
		vldErr := &ValidationError{}
		if errors.As(err, &ValidationError{}) {
			elVldErrors = append(elVldErrors, fmt.Errorf("idx %d: %w", idx, vldErr.Err))
			continue
		}
		return err
	}

	if elVldErrors != nil {
		return ValidationError{
			Field: name,
			Err:   fmt.Errorf("%w: %v", errValidation, elVldErrors),
		}
	}

	return nil
}

func strLenRestrMatcher(val string, lenRestrVal string) (bool, error) {
	expLen, err := strconv.Atoi(lenRestrVal)
	if err != nil {
		return false, errNotIntLenRestr
	}
	return utf8.RuneCountInString(val) == expLen, nil
}

func strRegexpRestrMatcher(val string, regexpRestrVal string) (bool, error) {
	isMatched, err := regexp.MatchString(regexpRestrVal, val)
	if err != nil {
		return false, fmt.Errorf("err on '%s' regexp restriction matching: %w", regexpRestrVal, err)
	}
	return isMatched, nil
}

func strInRestrMatcher(val string, inRestrVal string) (bool, error) {
	validVals := strings.Split(inRestrVal, ",")
	for _, vval := range validVals {
		if vval == val {
			return true, nil
		}
	}
	return false, nil
}

func intMinRestrMatcher(val int, minRestrVal string) (bool, error) {
	min, err := strconv.Atoi(minRestrVal)
	if err != nil {
		return false, errNotIntMinRestr
	}
	return val >= min, nil
}

func intMaxRestrMatcher(val int, maxRestrVal string) (bool, error) {
	max, err := strconv.Atoi(maxRestrVal)
	if err != nil {
		return false, errNotIntMaxRestr
	}
	return val <= max, nil
}

func intInRestrMatcher(val int, inRestrVal string) (bool, error) {
	validVals := strings.Split(inRestrVal, ",")
	isMatched := false
	for _, vval := range validVals {
		i, err := strconv.Atoi(vval)
		if err != nil {
			return false, errNotIntInRestr4int
		}
		isMatched = isMatched || i == val
	}
	return isMatched, nil
}
