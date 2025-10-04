import io from builtin;
import array from builtin; 
import types from builtin;


class Any {
    fn Any() {}
}


class Integer: Any {
    say x: int;
    fn Integer(x: int) {
        this.x = x;
    } 

    fn Value(): int {
        return this.x;
    }
}

fn main(): int32 {
    // say i: Integer = new Integer(20);

    say x: int = 2;
    say y: int = 3;

    // io.printf("type of = %s.    ", types.type(Integer));

    say arr:[][]Integer = array.create(Integer, x, y);
    say cx: [][]Integer = arr;
    say i: int = 0;
    say j: int = 0;

    say n : Integer = arr[i,j];

    arr[0,0] = new Integer(20);

    io.printf("final val = %d  ", arr[0,0].Value());
    // // x = arr[i,j];
    
    // arr[i,j+1] = 1 + j;
    
    // io.printf("value = %d ", cx[i, j+1]);
    // arr[0][0] = new Integer();


    return 0;
}