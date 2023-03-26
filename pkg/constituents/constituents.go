package constituents

import (
	"errors"
	"strings"
)

var ErrConstituentNotFound = errors.New("constituent not found")

type Constituent int32

const (
	// order and numbering same as NOAA is used
	C_M2   Constituent = 1
	C_S2   Constituent = 2
	C_N2   Constituent = 3
	C_K1   Constituent = 4
	C_M4   Constituent = 5
	C_O1   Constituent = 6
	C_M6   Constituent = 7
	C_MK3  Constituent = 8
	C_S4   Constituent = 9
	C_MN4  Constituent = 10
	C_NU2  Constituent = 11
	C_S6   Constituent = 12
	C_MU2  Constituent = 13
	C_2N2  Constituent = 14
	C_OO1  Constituent = 15
	C_LAM2 Constituent = 16
	C_S1   Constituent = 17
	C_M1   Constituent = 18
	C_J1   Constituent = 19
	C_MM   Constituent = 20
	C_SSA  Constituent = 21
	C_SA   Constituent = 22
	C_MSF  Constituent = 23
	C_MF   Constituent = 24
	C_RHO  Constituent = 25
	C_Q1   Constituent = 26
	C_T2   Constituent = 27
	C_R2   Constituent = 28
	C_2Q1  Constituent = 29
	C_P1   Constituent = 30
	C_2SM2 Constituent = 31
	C_M3   Constituent = 32
	C_L2   Constituent = 33
	C_2MK3 Constituent = 34
	C_K2   Constituent = 35
	C_M8   Constituent = 36
	C_MS4  Constituent = 37
	// no order assigned by NOAA
	C_LAMBDA2 Constituent = 100
	C_SIGMA1  Constituent = 101
	C_CHI1    Constituent = 102
	C_PI1     Constituent = 103
	C_PHI1    Constituent = 104
	C_THETA1  Constituent = 105
	// special cases
	C_UNKNOWN Constituent = 99999
)

func (c Constituent) String() string {
	switch c {
	// NOAA
	case C_M2:
		return "M2"
	case C_S2:
		return "S2"
	case C_N2:
		return "N2"
	case C_K1:
		return "K1"
	case C_M4:
		return "M4"
	case C_O1:
		return "O1"
	case C_M6:
		return "M6"
	case C_MK3:
		return "MK3"
	case C_S4:
		return "S4"
	case C_MN4:
		return "MN4"
	case C_NU2:
		return "NU2"
	case C_S6:
		return "S6"
	case C_MU2:
		return "MU2"
	case C_2N2:
		return "2N2"
	case C_OO1:
		return "OO1"
	case C_LAM2:
		return "LAM2"
	case C_S1:
		return "S1"
	case C_M1:
		return "M1"
	case C_J1:
		return "J1"
	case C_MM:
		return "MM"
	case C_SSA:
		return "SSA"
	case C_SA:
		return "SA"
	case C_MSF:
		return "MSF"
	case C_MF:
		return "MF"
	case C_RHO:
		return "RHO"
	case C_Q1:
		return "Q1"
	case C_T2:
		return "T2"
	case C_R2:
		return "R2"
	case C_2Q1:
		return "2Q1"
	case C_P1:
		return "P1"
	case C_2SM2:
		return "2SM2"
	case C_M3:
		return "M3"
	case C_L2:
		return "L2"
	case C_2MK3:
		return "2MK3"
	case C_K2:
		return "K2"
	case C_M8:
		return "M8"
	case C_MS4:
		return "MS4"
	// no noaa number
	case C_LAMBDA2:
		return "LAMBDA2"
	case C_SIGMA1:
		return "SIGMA1"
	case C_CHI1:
		return "CHI1"
	case C_PI1:
		return "PI1"
	case C_PHI1:
		return "PHI1"
	case C_THETA1:
		return "THETA1"
	}
	return "UNKNOWN"
}

func FromString(constituent string) (Constituent, error) {
	constituentUpperCase := strings.ToUpper(constituent)
	switch constituentUpperCase {
	// noaa number
	case "M2":
		return C_M2, nil
	case "S2":
		return C_S2, nil
	case "N2":
		return C_N2, nil
	case "K1":
		return C_K1, nil
	case "M4":
		return C_M4, nil
	case "O1":
		return C_O1, nil
	case "M6":
		return C_M6, nil
	case "MK3":
		return C_MK3, nil
	case "S4":
		return C_S4, nil
	case "MN4":
		return C_MN4, nil
	case "NU2":
		return C_NU2, nil
	case "S6":
		return C_S6, nil
	case "MU2":
		return C_MU2, nil
	case "2N2":
		return C_2N2, nil
	case "OO1":
		return C_OO1, nil
	case "LAM2":
		return C_LAM2, nil
	case "S1":
		return C_S1, nil
	case "M1":
		return C_M1, nil
	case "J1":
		return C_J1, nil
	case "MM":
		return C_MM, nil
	case "SSA":
		return C_SSA, nil
	case "SA":
		return C_SA, nil
	case "MSF":
		return C_MSF, nil
	case "MF":
		return C_MF, nil
	case "RHO":
		return C_RHO, nil
	case "Q1":
		return C_Q1, nil
	case "T2":
		return C_T2, nil
	case "R2":
		return C_R2, nil
	case "2Q1":
		return C_2Q1, nil
	case "P1":
		return C_P1, nil
	case "2SM2":
		return C_2SM2, nil
	case "M3":
		return C_M3, nil
	case "L2":
		return C_L2, nil
	case "2MK3":
		return C_2MK3, nil
	case "K2":
		return C_K2, nil
	case "M8":
		return C_M8, nil
	case "MS4":
		return C_MS4, nil
	// no noaa number
	case "LAMBDA2":
		return C_LAMBDA2, nil
	case "SIGMA1":
		return C_SIGMA1, nil
	case "CHI1":
		return C_CHI1, nil
	case "PI1":
		return C_PI1, nil
	case "PHI1":
		return C_PHI1, nil
	case "THETA1":
		return C_THETA1, nil
	}
	return C_UNKNOWN, ErrConstituentNotFound
}

func GetAllConstituents() []Constituent {
	return []Constituent{
		// noaa numbers
		C_M2,
		C_S2,
		C_N2,
		C_K1,
		C_M4,
		C_O1,
		C_M6,
		C_MK3,
		C_S4,
		C_MN4,
		C_NU2,
		C_S6,
		C_MU2,
		C_2N2,
		C_OO1,
		C_LAM2,
		C_S1,
		C_M1,
		C_J1,
		C_MM,
		C_SSA,
		C_SA,
		C_MSF,
		C_MF,
		C_RHO,
		C_Q1,
		C_T2,
		C_R2,
		C_2Q1,
		C_P1,
		C_2SM2,
		C_M3,
		C_L2,
		C_2MK3,
		C_K2,
		C_M8,
		C_MS4,
		// no noaa number
		C_LAMBDA2,
		C_SIGMA1,
		C_CHI1,
		C_PI1,
		C_PHI1,
		C_THETA1}
}
