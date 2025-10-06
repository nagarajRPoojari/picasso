import io from builtin;
import array from builtin;
import types from builtin;

// class Integer {
//     say x: int;
//     fn Integer(x: int) {
//         this.x = x;
//     }

//     fn Print() {
//         return 
//     }
// }

class Test {
    fn Test() {}
    fn Update(arr: [][]int): [][]int {
        // io.printf(" ==== %s.  ", types.type(arr));

        io.printf("value = %d. ", arr[0, 1]);

        return arr;
    }   
}
fn main(): int32 {
    say arr: [][]int = array.create(int, 5, 4);
    arr[0, 1] = 1890;


    // io.printf("value = %d. ", arr[0]);
    // say n: []int = arr;

    // io.printf("%d.  ", n[0]);

    say t: Test = new Test();
    say x: [][]int = t.Update(arr);

    io.printf("value = %d. ", x[0, 1]);

    return 0;
}