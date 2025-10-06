import io from builtin;
import array from builtin;
import types from builtin;


class Any {
    fn Any() {}
    fn Print() {
        io.printf("unimplemented");
    }
}

class Integer: Any {
    say x: int;
    fn Integer(x: int) {
        this.x = x;
    }

    fn Print() {
        io.printf("x=%d. ", this.x);
    }
}

class Test {
    fn Test() {}
    fn Update(arr: [][]Integer): [][]Integer {
        arr[0,1].Print();
        say n: Integer = arr[0,1];
        n.x = 90;
        return arr;
    }   
}
fn main(): int32 {
    say arr: [][]Integer;
    arr = array.create(Integer, 5, 4);
    arr[0, 1] = new Integer(100);

    say t: Test = new Test();
    say x: [][]Integer = t.Update(arr);


    t.Update(arr);
    return 0;
}