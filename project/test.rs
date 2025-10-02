import io from builtin;

class Any {
    fn Any() {}
}


class Integer: Any {
    say x: int;
    fn Integer(x: int) {
        this.x = x;
    } 

    fn Value(): Any {
        return null;
    }
}

fn main(): int32 {
    say i: Integer = new Integer(20);

    say x: Any = i.Value();

    // say arr: [][]Integer = arrays.create(Integer, 4, 5);
    // arr[0][0] = new Integer();

    io.printf("print , %s ", Any);

    return 0;
}