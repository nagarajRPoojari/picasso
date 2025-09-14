class DirectoryReader {
  say y: double = 112;
  say x: Math = new Math();

  fn sum(x: int): int {
    return this.y + this.x.pi;
  }
}

class Math {
  say pi: double = 100;
}

fn main(): int32 {
    say a: DirectoryReader = new DirectoryReader();
    say z: int = a.sum(10);
    return 0;
}