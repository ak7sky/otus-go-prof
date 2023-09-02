package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type (
	UserRole string

	UserPhysParams struct {
		Age    int `validate:"min:18|max:50"`
		Height int `validate:"min:140|max:240"`
		Weight int `validate:"min:45|max:150"`
	}

	User struct {
		PhysParams UserPhysParams `validate:"nested"`
		ID         string         `json:"id" validate:"len:36"`
		Name       string
		Experience int             `validate:"min:2|max:50"`
		Email      string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role       UserRole        `validate:"in:admin,stuff"`
		Phones     []string        `validate:"len:11"`
		meta       json.RawMessage //nolint:unused
	}
)

func TestValidateReturnsAppErr(t *testing.T) {
	testCases := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		//--------------------- Application error cases ---------------------
		{
			name:        "app-err:not struct",
			in:          map[int]int{},
			expectedErr: errNotStruct,
		},
		{
			name: "app-err:empty validate tag",
			in: struct {
				Pswd string `validate:""`
			}{},
			expectedErr: fmt.Errorf("field 'Pswd': %w", errEmptyVldTag),
		},
		{
			name: "app-err:unsupported type",
			in: struct {
				BooksByAuthor map[string][]string `validate:"len:10000"`
			}{},
			expectedErr: fmt.Errorf("field 'BooksByAuthor', type map: %w", errFieldType),
		},
		{
			name: "app-err:not nested on struct",
			in: struct {
				Nested struct{} `validate:"nest"`
			}{},
			expectedErr: fmt.Errorf("field 'Nested': %w", errNotNested),
		},
		{
			name: "app-err:unknown restr for int",
			in: struct {
				Shoes int `validate:"even:true"`
			}{},
			expectedErr: fmt.Errorf("field 'Shoes', type int, restr type 'even': %w", errRestrType),
		},
		{
			name: "app-err:unknown restr for str",
			in: struct {
				Palindrome string `validate:"palindrome:true"`
			}{},
			expectedErr: fmt.Errorf("field 'Palindrome', type string, restr type 'palindrome': %w", errRestrType),
		},
		{
			name: "app-err:len restr value",
			in: struct {
				Name string `validate:"len:any"`
			}{},
			expectedErr: fmt.Errorf("field 'Name': %w", errNotIntLenRestr),
		},
		{
			name: "app-err:min restr value",
			in: struct {
				Num int `validate:"min:any"`
			}{},
			expectedErr: fmt.Errorf("field 'Num': %w", errNotIntMinRestr),
		},
		{
			name: "app-err:max restr value",
			in: struct {
				Num int `validate:"max:any"`
			}{},
			expectedErr: fmt.Errorf("field 'Num': %w", errNotIntMaxRestr),
		},
		{
			name: "app-err:in restr on int",
			in: struct {
				Num int `validate:"in:10,100,thousand"`
			}{},
			expectedErr: fmt.Errorf("field 'Num': %w", errNotIntInRestr4int),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()
			err := Validate(tc.in)
			fmt.Println(err)
			require.Equal(t, tc.expectedErr, err)
		})
	}
}

func TestValidateReturnsSingleElValidationErrs(t *testing.T) {
	testCases := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		//--------------------- Validation error cases ---------------------
		{
			name: "vld-err:nil",
			in: struct {
				ID string `validate:"len:36"`
			}{
				ID: "d56e42f5-1aa8-4a54-96d9-77be14e725be",
			},
			expectedErr: nil,
		},
		{
			name: "vld-err:len",
			in: struct {
				ID string `validate:"len:36"`
			}{
				ID: "d56e42f5",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   fmt.Errorf("%w: val 'd56e42f5', %v", errValidation, []string{"len:36"}),
				},
			},
		},
		{
			name: "vld-err:min",
			in: struct {
				Experience int `validate:"min:2|max:50"`
			}{
				Experience: 1,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Experience",
					Err:   fmt.Errorf("%w: val '1', %v", errValidation, []string{"min:2"}),
				},
			},
		},
		{
			name: "vld-err:max",
			in: struct {
				Experience int `validate:"min:2|max:50"`
			}{
				Experience: 999,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Experience",
					Err:   fmt.Errorf("%w: val '999', %v", errValidation, []string{"max:50"}),
				},
			},
		},
		{
			name: "vld-err:regexp",
			in: struct {
				Email string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
			}{
				Email: "any@@any.com",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Email",
					Err:   fmt.Errorf("%w: val 'any@@any.com', %v", errValidation, []string{"regexp:^\\w+@\\w+\\.\\w+$"}),
				},
			},
		},
		{
			name: "vld-err:in",
			in: struct {
				Role string `validate:"in:admin,stuff"`
			}{
				Role: "moder",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Role",
					Err:   fmt.Errorf("%w: val 'moder', %v", errValidation, []string{"in:admin,stuff"}),
				},
			},
		},
		{
			name: "vld-err:len on []string",
			in: struct {
				Phones []string `validate:"len:11"`
			}{
				Phones: []string{"89123456789", "898765432109"}, // 2nd num is invalid
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Phones",
					Err:   fmt.Errorf("%w: %v", errValidation, []error{errors.New("idx 1: violations: val '898765432109', [len:11]")}),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()
			err := Validate(tc.in)
			fmt.Println(err)
			require.Equal(t, err, tc.expectedErr)
		})
	}
}

func TestValidateReturnsMultipleValidationErrs(t *testing.T) {
	testCases := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "vld-err:multiple",
			in: User{
				PhysParams: UserPhysParams{Age: 17, Height: 300, Weight: 200},
				ID:         "d56e42f5",
				Name:       "John",
				Experience: 1,
				Email:      "any@@any.com",
				Role:       "moder",
				Phones:     []string{"891237456789", "898765432107"},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "PhysParams",
					Err: ValidationErrors{
						ValidationError{
							Field: "Age",
							Err:   fmt.Errorf("%w: val '17', %v", errValidation, []string{"min:18"}),
						},
						ValidationError{
							Field: "Height",
							Err:   fmt.Errorf("%w: val '300', %v", errValidation, []string{"max:240"}),
						},
						ValidationError{
							Field: "Weight",
							Err:   fmt.Errorf("%w: val '200', %v", errValidation, []string{"max:150"}),
						},
					},
				},
				ValidationError{
					Field: "ID",
					Err:   fmt.Errorf("%w: val 'd56e42f5', %v", errValidation, []string{"len:36"}),
				},
				ValidationError{
					Field: "Experience",
					Err:   fmt.Errorf("%w: val '1', %v", errValidation, []string{"min:2"}),
				},
				ValidationError{
					Field: "Email",
					Err:   fmt.Errorf("%w: val 'any@@any.com', %v", errValidation, []string{"regexp:^\\w+@\\w+\\.\\w+$"}),
				},
				ValidationError{
					Field: "Role",
					Err:   fmt.Errorf("%w: val 'moder', %v", errValidation, []string{"in:admin,stuff"}),
				},
				ValidationError{
					Field: "Phones",
					Err: fmt.Errorf("%w: %v", errValidation, []error{
						errors.New("idx 0: violations: val '891237456789', [len:11]"),
						errors.New("idx 1: violations: val '898765432107', [len:11]"),
					}),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()
			err := Validate(tc.in)
			fmt.Println(err)
			require.Equal(t, err, tc.expectedErr)
		})
	}
}
