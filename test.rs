class DirectoryReader {
  say y: int = 112;
  // say x: Math = new Math();

  fn sum(x: int): int {
    say n: int = this.y;
    return this.y + x;
  }
}

class Math {
  say pi: int = 100;
}

fn main(): int32 {
    say a: DirectoryReader = new DirectoryReader();
    say z: int = a.sum(10);

    return 0;
}