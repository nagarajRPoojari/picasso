import io from builtin;
import types from builtin;
// import math;


class Math {
    say x: int = 190;
    fn Math() {}
    fn Print(a : int) {
        io.printf("printing from Math class.    ");
    }
}


class Calculator: Math {
    say y: int = 10;
    fn Calculator() {
        // constructor
    }

    fn add(a: int, b: int): int {
        return a + b;
    }

    fn Print(b : int) {
        say n: int = this.x;
        io.printf("printing from Math Calculator class.   %d ", n);
    }

    // fn multiply(a: int, b: int): int {
    //     return a * b;
    // }

    // fn divide(a: int, b: int): float64 {
    //     if (b == 0) {
    //         io.printf("Division by zero!\n");
    //         return 0.0;
    //     } else {
    //         return a / b;
    //     }
	// 	return 0;
    // }

    // fn max(a: int, b: int): int {
    //     if( a > b ){
    //         return a;
    //     } else {
	// 		return b;
	// 	}
	// 	return b;
    // }
}


fn main(): int32 {
    say c: Calculator = new Calculator();
    c.add(10,20);

    // c.Print(99);


    say a: Na = new Calculator();
    a.Print(9);


    io.printf("typeof a = %s ", types.type(a));

    // say calc: Calculator = new Calculator();
    // say adv: AdvancedMath = new AdvancedMath();
	// say n: int = 11;
	// io.printf("factorial of %d = %d  ", n, adv.factorial(n));
    return 0;
}
