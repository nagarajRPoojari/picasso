
say PI: double = 3.14;

class DirectoryReader {
  say y: double = 112;
  say x: Math = new Math();

  fn sum(m: Math): int {
    return m.pi;
  }

  fn math(): Math {
    return this.x;
  }
}

class Math {
  say pi: double = 123;
}

fn main(): int32 {
    say a: DirectoryReader = new DirectoryReader();
    say z: int = a.sum(new Math());



    return 0;
}