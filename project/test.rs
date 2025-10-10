import io from builtin;
fn main(): int32 {
    say l: int = 0;
    say h: int = 1000000;

    foreach i in l..h {
       io.printf("hi, %d. ", i);
    }
    return 0;
}

// 1447910000 -  time taken by go
// 1181128000 -  time taken by mine