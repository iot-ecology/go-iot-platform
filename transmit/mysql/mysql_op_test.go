package mysql

import "testing"

func TestGet(t *testing.T) {
	connection, err := InitMySQLConnection("root", "127.0.0.1", "root123@", "iot", 3306, 1)
	if err != nil {
		t.Error(err)
	}
	sqlConnection, err := InitMySQLConnection("root", "127.0.0.1", "root123@", "iot", 3306, 1)

	if err != nil {
		t.Error(err)
	}
	if connection == sqlConnection {
		t.Logf("success ")
	}

}
