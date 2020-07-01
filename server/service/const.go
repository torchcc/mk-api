package service

type PackageTarget = int8
type FeeType = int8
type TransType = int8

const (
	AnyGender       PackageTarget = 0
	MALE                          = 1
	UnMarriedFemale               = 2
	MarriedFemale                 = 3
	Female                        = 2
)

const (
	Income FeeType = 1
	// Outgo          = 2
)

const (
	Earned TransType = 1
)
