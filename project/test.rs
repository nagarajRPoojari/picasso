import io from builtin;
import array from builtin;

class Worker {
  say id: int;
  fn Worker(id: int) {
    this.id = id;
  }

  fn do_1() {
    foreach i in 1..1000 {
      io.printf("--> hello_1 %d \n", i);
    }
  }

  fn do_2() {
    foreach i in 1..1000 {
      io.printf("--> hello_2 %d \n", i);
    }
  }
}



fn start(): int32 {
    printf("started... \n");
    

    say worker1: Worker = new Worker(0);
    say worker2: Worker = new Worker(1);

    // worker1.do();

    thread(worker1.do_1);
    thread(worker2.do_2);

}

// 1447910000 -  time taken by go
// 1181128000 -  time taken by mine

