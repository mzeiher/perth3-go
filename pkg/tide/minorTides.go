package tide

// the algorithm to infer minor tides and compute heights uses another order of
// constituents so we have to map them to the [2][28]height array
//
// tidal array expects the following order:
// [0]=Q1, [1]=O1, [2]=P1, [3]=K1, [4]=N2, [5]=M2, [6]=S2, [7]=K2, [26]=S1, [27]=M4
// we get the heights in the following order
// [0]=Q1, [1]=O1, [2]=P1, [3]=S1], [4]=K1, [5]=N2, [6]=M2, [7]=S2, [8]=K2, [9]=M4
// so we have to re-map them
func MapComputedHeightsToHeightArray(heights [][]float64) [][]float64 {
	newTideHeights := make([][]float64, 2)

	// init new array
	for cossin := 0; cossin < 2; cossin++ {
		newTideHeights[cossin] = make([]float64, 28)
		for tideArg := 0; tideArg < 28; tideArg++ {
			newTideHeights[cossin][tideArg] = 0
		}
	}

	// Q1
	newTideHeights[0][0] = heights[0][0]
	newTideHeights[1][0] = heights[1][0]

	// O1
	newTideHeights[0][1] = heights[0][1]
	newTideHeights[1][1] = heights[1][1]

	// P1
	newTideHeights[0][2] = heights[0][2]
	newTideHeights[1][2] = heights[1][2]

	// K1
	newTideHeights[0][3] = heights[0][4]
	newTideHeights[1][3] = heights[1][4]

	// N2
	newTideHeights[0][4] = heights[0][5]
	newTideHeights[1][4] = heights[1][5]

	// M2
	newTideHeights[0][5] = heights[0][6]
	newTideHeights[1][5] = heights[1][6]

	// S2
	newTideHeights[0][6] = heights[0][7]
	newTideHeights[1][6] = heights[1][7]

	// K2
	newTideHeights[0][7] = heights[0][8]
	newTideHeights[1][7] = heights[1][8]

	// S1
	newTideHeights[0][26] = heights[0][3]
	newTideHeights[1][26] = heights[1][3]

	// m4
	newTideHeights[0][27] = heights[0][9]
	newTideHeights[1][27] = heights[1][9]

	return newTideHeights

}

// infer the minor tides at specific location from major tides (Q1, O1, P1, K1, N2, M2, S2, K2, S1, M4)
func InferMinorTideHeights(heights [][]float64) [][]float64 {

	heights[0][8] = 0.263*heights[0][0] - 0.0252*heights[0][1]    // 2Q1 HCos
	heights[1][8] = 0.263*heights[1][0] - 0.0252*heights[1][1]    // 2Q1 HSin
	heights[0][9] = 0.297*heights[0][0] - 0.0264*heights[0][1]    // sigma1 HCos
	heights[1][9] = 0.297*heights[1][0] - 0.0264*heights[1][1]    // sigma1 HSin
	heights[0][10] = 0.164*heights[0][0] + 0.0048*heights[0][1]   // rho1 HCos
	heights[1][10] = 0.164*heights[1][0] + 0.0048*heights[1][1]   // rho1 HSin
	heights[0][11] = 0.0140*heights[0][1] + 0.0101*heights[0][3]  // M1 HCos
	heights[1][11] = 0.0140*heights[1][1] + 0.0101*heights[1][3]  // M1 HSin
	heights[0][12] = 0.0389*heights[0][1] + 0.0282*heights[0][3]  // M1 HCos
	heights[1][12] = 0.0389*heights[1][1] + 0.0282*heights[1][3]  // M1 HSin
	heights[0][13] = 0.0064*heights[0][1] + 0.0060*heights[0][3]  // chi1 HCos
	heights[1][13] = 0.0064*heights[1][1] + 0.0060*heights[1][3]  // chi1 HSin
	heights[0][14] = 0.0030*heights[0][1] + 0.0171*heights[0][3]  // pi1 HCos
	heights[1][14] = 0.0030*heights[1][1] + 0.0171*heights[1][3]  // pi1 HSin
	heights[0][15] = -0.0015*heights[0][1] + 0.0152*heights[0][3] // phi1 HCos
	heights[1][15] = -0.0015*heights[1][1] + 0.0152*heights[1][3] // phi1 HSin
	heights[0][16] = -0.0065*heights[0][1] + 0.0155*heights[0][3] // theta1 HCos
	heights[1][16] = -0.0065*heights[1][1] + 0.0155*heights[1][3] // theta1 HSin
	heights[0][17] = -0.0389*heights[0][1] + 0.0836*heights[0][3] // J1 HCos
	heights[1][17] = -0.0389*heights[1][1] + 0.0836*heights[1][3] // J1 HSin
	heights[0][18] = -0.0431*heights[0][1] + 0.0613*heights[0][3] // OO1 HCos
	heights[1][18] = -0.0431*heights[1][1] + 0.0613*heights[1][3] // OO1 HSin
	heights[0][19] = 0.264*heights[0][4] - 0.0253*heights[0][5]   // 2N2 HCos
	heights[1][19] = 0.264*heights[1][4] - 0.0253*heights[1][5]   // 2N2 HSin
	heights[0][20] = 0.298*heights[0][4] - 0.0264*heights[0][5]   // mu2 HCos
	heights[1][20] = 0.298*heights[1][4] - 0.0264*heights[1][5]   // mu2 HSin
	heights[0][21] = 0.165*heights[0][4] + 0.00487*heights[0][5]  // nu2 HCos
	heights[1][21] = 0.165*heights[1][4] + 0.00487*heights[1][5]  // nu2 HSin
	heights[0][22] = 0.0040*heights[0][5] + 0.0074*heights[0][6]  // lambda2 HCos
	heights[1][22] = 0.0040*heights[1][5] + 0.0074*heights[1][6]  // lambda2 HSin
	heights[0][23] = 0.0131*heights[0][5] + 0.0326*heights[0][6]  // L2 HCos
	heights[1][23] = 0.0131*heights[1][5] + 0.0326*heights[1][6]  // L2 HSin
	heights[0][24] = 0.0033*heights[0][5] + 0.0082*heights[0][6]  // L2 HCos
	heights[1][24] = 0.0033*heights[1][5] + 0.0082*heights[1][6]  // L2 HSin
	heights[0][25] = 0.0585 * heights[0][6]                       // T2 HCos
	heights[1][25] = 0.0585 * heights[1][6]                       // T2 HSin

	return heights
}
