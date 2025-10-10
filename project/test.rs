import io from builtin;
import array from builtin;

class Any {
  say x: int;
  fn Any() {}
}

class Integer: Any {
  say y: int;
  fn Integer() {
    this.Any();
    
  }
}





fn main(): int32 {
    say size: int = 100;
    say arr: []int = array.create(int, size);
    foreach i in 0..size {
      arr[i] = i * 10;
    }

    foreach i in 0..size {
      io.printf("arr[%d]=%d\n", i, arr[i]);
    }
    return 0;
}

// 1447910000 -  time taken by go
// 1181128000 -  time taken by mine