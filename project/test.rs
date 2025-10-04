import io from builtin;
import array from builtin; 


// class Any {
//     fn Any() {}
// }


// class Integer: Any {
//     say x: int;
//     fn Integer(x: int) {
//         this.x = x;
//     } 

//     fn Value(): Any {
//         return null;
//     }
// }

fn main(): int32 {
    // say i: Integer = new Integer(20);

    say x: int = 2;
    say y: int = 3;

    say arr:[][]int = array.create(x, y);
    say cx: [][]int = arr;
    say i: int = 0;
    say j: int = 0;

    say n : int = arr[i,j];

    io.printf("value = %d ", arr[i,j]);
    // x = arr[i,j];
    
    arr[i,j+1] = 1 + j;
    
    io.printf("value = %d ", cx[i, j+1]);
    // arr[0][0] = new Integer();

    // io.printf("print , %s ", Integer);

    return 0;
}