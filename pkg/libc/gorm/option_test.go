package gorm

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
)

func TestSetOption(t *testing.T) {
	_db, _, err := sqlmock.New()
	if err != nil {
		t.Error(err)
		return
	}

	db, err := gorm.Open("postgres", _db)
	if err != nil {
		t.Error(err)
		return
	}

	cases := [][2]interface{}{
		{BlockGlobalUpdateOpt(true), nil},
		{LogModOpt(true), nil},
		{SetMaxIdleConnsOpt(2), nil},
		{SetMaxOpenConnsOpt(3), nil},
		{SetConnMaxLifetimeOpt(300), nil},
	}

	for _, tc := range cases {
		err := tc[0].(Option)(db)
		expect, _ := tc[1].(error)

		if err != expect {
			t.Errorf("err occurred but not eq \nCases:\t%+v\n\t%+v != %+v", tc, expect, err)
		}
	}
}
