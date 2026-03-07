import java.util.concurrent.CountDownLatch;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class ConcurrencyStress {

    private static final int TASK_COUNT = 30000;
    private final ExecutorService pool;
    private final CountDownLatch latch;

    public ConcurrencyStress() {
        int cores = Runtime.getRuntime().availableProcessors();
        this.pool = Executors.newFixedThreadPool(cores); // optimized
        this.latch = new CountDownLatch(TASK_COUNT);
    }

    private void worker(int count) {
        int c = 0;
        for (int i = 0; i < count; i++) {
            c = c * 100 + c - 100;
            if (c > 10) {
                c -= 29;
            } else {
                c += 90;
            }
        }
        latch.countDown();
    }

    public void testManyCPUBoundTasks() throws InterruptedException {
        for (int i = 0; i < TASK_COUNT; i++) {
            final int work = i * 100;
            pool.execute(() -> worker(work));
        }

        latch.await();
        pool.shutdown();
        System.out.println("\nall tasks done");
    }

    public static void main(String[] args) throws InterruptedException {
        ConcurrencyStress stress = new ConcurrencyStress();
        stress.testManyCPUBoundTasks();
    }
}