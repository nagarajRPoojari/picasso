import io from builtin;
import array from builtin;


class Worker {
  say id: int;
  fn Worker(id: int) {
    this.id = id;
  }

  fn doss() {
    printf("> %d\n", this.id);
    this.recurse();
  }

  fn recurse() {


    this.id = this.id - 1;
    if(this.id == 0){
      return;
    }else {}
    printf("> %d\n", this.id);
    this.recurse();

  }
}




fn start() {
    printf("started... \n");


    say worker: Worker = new Worker(100);

    thread(worker.dos);

}

// 1447910000 -  time taken by go
// 1181128000 -  time taken by mine

