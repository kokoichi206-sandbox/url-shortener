package usecase_test

import (
	"testing"

	"github.com/kokoichi206-sandbox/url-shortener/usecase"
	"github.com/stretchr/testify/assert"
)

func Test_Usecase_GetRoomUsers(t *testing.T) {
	t.Parallel()

	type args struct {
		num int
	}

	testCases := map[string]struct {
		args    args
		wantReg string
		wantErr string
	}{
		"success": {
			args: args{
				num: 3,
			},
			wantReg: "^[a-zA-Z0-9]{3}$",
		},
		"success: length 5": {
			args: args{
				num: 5,
			},
			wantReg: "^[a-zA-Z0-9]{5}$",
		},
	}

	for name, tc := range testCases {
		name := name
		tc := tc

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// Arrange

			// Act
			got, err := usecase.GenerateRandomString(tc.args.num)

			// Assert
			assert.Regexp(t, tc.wantReg, got, "result does not match")
			if tc.wantErr == "" {
				assert.Nil(t, err, "error should be nil")
			} else {
				assert.Equal(t, tc.wantErr, err.Error(), "result does not match")
			}
		})
	}
}
