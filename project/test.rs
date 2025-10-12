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


class Student {
  say x: int;
  say y: int;
  fn Class() {}
}



fn main(): int32 {
    say size: int = 100;
    say arr: []int = array.create(int, size);
    say str: string = "hello world";
    io.printf("length of %s  \n", str);

    // say r: string = "";
    // say n: int;

    // io.printf("what is your name ?? \n");
    // io.scanf("%s %d", r, n);

    // io.printf("Hi, %s %d \n", r, n);
    say name: string = "nagaraj";
    say fp: string = io.fopen("temp.txt", "w+");


    fprintf(fp, "Name: %s\n", name);
    io.fflush(fp);
    io.fseek(fp, 0, 0);
    // fprintf(fp, "Name2: %s\n", name);


    say r: string = "";

  
    // io.fscanf(fp, "Name: %s\n", r);

    // io.printf("read from file-> %s\n", r);

    io.fputs("hi hello how are you", fp);


    say x: string = "";
    io.fflush(fp);
    io.fseek(fp, 0, 0);
    
    io.fgets(x, 20, fp);


    io.printf("read-> %s \n", x);
    


    io.fclose(fp);
    return 0;
}

// 1447910000 -  time taken by go
// 1181128000 -  time taken by mine

