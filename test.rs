class DirectoryReader {
  say static x: float32 = 200;
  say y: float64 = 1;

  fn sum(x: int, y: int): int {
    return this.y;
  }

  fn mul(x: int, y: int): int {
    return x * y;
  }

  fn sub(x: int, y: int): int {
    return x - y;
  }
}

class Math {
  say static x: float32 = 200;
  say y: float64 = 1;

  fn sum(x: int, y: int): int {
    return this.y;
  }

  fn mul(x: int, y: int): int {
    return x * y;
  }

  fn sub(x: int, y: int): int {
    return x - y;
  }
}
fn main() {
    say x: DirectoryReader = new DirectoryReader();
    say z: float64 = x.sum(10,10) + x.mul(10,20) ;
}