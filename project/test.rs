import io;
import math;

class Calculator {
    fn Calculator() {
        // constructor
    }

    fn add(a: int, b: int): int {
        return a + b;
    }

    fn multiply(a: int, b: int): int {
        return a * b;
    }

    fn divide(a: int, b: int): float64 {
        if (b == 0) {
            io.printf("Division by zero!\n");
            return 0.0;
        } else {
            return a / b;
        }
		return 0;
    }

    fn max(a: int, b: int): int {
        if( a > b ){
            return a;
        } else {
			return b;
		}
		return b;
    }
}


fn main(): int32 {
    say calc: Calculator = new Calculator();
    say adv: AdvancedMath = new AdvancedMath();
	say n: int = 11;
	io.printf("factorial of %d = %d  ", n, adv.factorial(n));
    return 0;
}
