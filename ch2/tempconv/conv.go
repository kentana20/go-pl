package tempconv

// CToF 摂氏 → 華氏
func CToF(c Celsius) Fahrenheit {
	return Fahrenheit(c*9/5 + 32)
}

// FToC 華氏 → 摂氏
func FToC(f Fahrenheit) Celsius {
	return Celsius((f - 32) * 5 / 9)
}

// CToK 摂氏 → 絶対温度
func CToK(c Celsius) Kelvin {
	return Kelvin(c - AbsoluteZeroC)
}

// KToC 絶対温度 → 摂氏
func KToC(k Kelvin) Celsius {
	return Celsius(Celsius(k) + AbsoluteZeroC)
}

// FToK 華氏 → 絶対温度
func FToK(f Fahrenheit) Kelvin {
	return Kelvin(FToC(f) - AbsoluteZeroC)
}

// KToF 絶対温度 → 華氏
func KToF(k Kelvin) Fahrenheit {
	return Fahrenheit(CToF(Celsius(k) + AbsoluteZeroC))
}
