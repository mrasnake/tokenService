package internal

import (
	"testing"
)

// errString is used for comparison of error strings
// by returning the string "nil" when given a nil value
func errString(err error) string {
	if err == nil {
		return "nil"
	}
	return err.Error()
}

// WriteTokenTest is a struct to hold all testing parameters
// for each test performed on the WriteToken.
type WriteTokenTest struct {
	name           string
	inputRequest   *WriteTokenRequest
	outputResponse *WriteTokenResponse
	wantErr        string
}

// Test_WriteToken builds a list of custom structs and loops through each of them
// performing the associated unit test on WriteToken() with the specified parameters.
func Test_WriteToken(t *testing.T) {

	// tests contains all the parameters, checks and expected results of each test.
	tests := []WriteTokenTest{
		// WriteToken Test #1 checks that when given a bad 'WriteTokenRequest' value
		// the WriteToken function fails correctly.
		{
			name:           "WriteToken Test #1 - Validation Failure",
			inputRequest:   &WriteTokenRequest{},
			wantErr:        "invalid request: secret: cannot be blank.",
			outputResponse: nil,
		},
		// WriteToken Test #1 checks when passing valid
		// parameters to WriteToken, no errors are returned.
		{
			name: "WriteToken Test #1 - Success",
			inputRequest: &WriteTokenRequest{
				Secret: "foo",
			},
			wantErr:        "nil",
			outputResponse: &WriteTokenResponse{},
		},
	}

	// This loops through each item in the tests list, uses the individual parameters
	// to prepare and perform the unit test and compares the results.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			s, err := NewService()
			if err != nil {
				t.Fatalf("could not initialize service: %v", err)
			}

			resp, err := s.WriteToken(tt.inputRequest)
			if tt.wantErr != errString(err) {
				t.Errorf("unexpected error, got=%v; want=%v", errString(err), tt.wantErr)
			}

			if tt.outputResponse != nil {
				if resp == nil {
					t.Error("expected a response, none given.")
				}
			} else {
				if resp != nil {
					t.Error("no response expected.")
				}
			}

			return

		})
	}
}
