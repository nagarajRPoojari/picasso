
import math;

class Math {
	say PI: float64 = 3.14;
	fn Math() {

	}

	fn add(a: int, b: int): int {
		return a + b;
	}

	fn multiply(a: int, b: int): int {
		return a * b;
	}
}

fn main(): int32 {
	say m: IO = new IO();
	// io.printf("%f + %f = %d   ", 10, 20, m.add(10,20));

	return 0;
}