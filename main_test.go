package main

import (
	"testing"
)

func TestNew(t *testing.T) {
	New()
}

func TestCheck(t *testing.T) {
	i := Input{
		UserCipherText: "f715af90e5a3145e5153f9e847f6f12d9965ebbc58e92e810813ea1bcf47e82a#1591320367",
		UserCode:       "6411",
		Action:         "check",
	}
	Check(&i)

}
