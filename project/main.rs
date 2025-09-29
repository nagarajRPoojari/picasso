import io from builtin;

interface Any {
    fn hey();
}

class Integer is Any {
    fn Integer(x: int) {}
    fn hey() {
        // do something
    }
}

class Array[T Any] {
    fn Array(size: int) {

    }
    fn append(x: T) {

    }
}

fn main(): int32 {
    say x: Array[int] = new Array[int]();
    

    return 0;
}