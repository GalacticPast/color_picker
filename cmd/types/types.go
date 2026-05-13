package types

type Color struct {
	R uint8
	G uint8
	B uint8
	Hex_code string 
}

type Color_container struct {
	Colors [5]Color
}
