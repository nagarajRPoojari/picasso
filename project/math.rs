class AdvancedMath {
    say PI: float64 = 3.14159265;

    fn AdvancedMath() {

    }

    fn circleArea(radius: float64): float64 {
        return this.PI * radius * radius;
    }

    fn circleCircumference(radius: float64): float64 {
        return 2 * this.PI * radius;
    }

    fn power(base: int, exp: int): int {
        // only supports small exponent, no loops, recursive instead
        if (exp == 0 ){
            return 1;
        } else {
            return base * this.power(base, exp-1);
        }

		return 0;
    }

	fn factorial(n :int): int {
		if (n==1){
			return 1;
		}else {
			return n * this.factorial(n-1);
		}

		return 0;
	}
}