import io from builtin;
import array from builtin;

import stdlib from c;
import stdio from c;

fn main(): int32 {

    say fp: string = io.fopen("temp.txt", "r+");

    say MAX_LINES: int = 100;
    say BUFFERSIZE: int = 256;
    say BUFFER: string = "";

    foreach i in 0..MAX_LINES {
      say r: int = io.fgets(BUFFER, 256, fp);
      io.printf("%s", BUFFER);

      if(r==0){
        io.fclose(fp);
        return 0;
      }else {}
    }


    io.fclose(fp);
    return 0;
}

// 1447910000 -  time taken by go
// 1181128000 -  time taken by mine




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



fn start(): int32 {
    say size: int = 100;
    say arr: []int = array.create(int, size);
    say str: string = "hello world";
    io.printf("length of %s  \n", str);

    say r: string = "";
    say n: int;

    io.printf("what is your name ?? \n");
    io.scanf("%s %d", r, n);

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
    // return 0;
}

// 1447910000 -  time taken by go
// 1181128000 -  time taken by mine

