import io from builtin;
import array from builtin;

class Worker {
  say id: int;
  fn Worker(id: int) {
    this.id = id;
  }

  fn do() {
    foreach i in 1..1000 {
      io.printf("--> hello_%d %d \n", this.id, i);
    }
  }
}



fn start(): int32 {
    printf("started... \n");
    

    say worker1: Worker = new Worker(967);
    say worker2: Worker = new Worker(1);

    // worker1.do();

    thread(worker1.do);
    thread(worker1.do);

}

// 1447910000 -  time taken by go
// 1181128000 -  time taken by mine

