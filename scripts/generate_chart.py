
import pandas as pd
import matplotlib.pyplot as plt

try:
    df = pd.read_csv('benchmark_results.csv')
except FileNotFoundError:
    print("Error: benchmark_results.csv not found. Run the bash script first.")
    exit(1)

df = df.sort_values('mean')

plt.figure(figsize=(12, 7))
bars = plt.bar(df['command'], df['mean'], yerr=df['stddev'], 
                capsize=5, color='skyblue', edgecolor='navy')

plt.ylabel('Execution Time (seconds)')
plt.title('Mandelbrot Benchmark (N=1000)', fontsize=14, fontweight='bold')
plt.xticks(rotation=45, ha='right')
plt.grid(axis='y', linestyle='--', alpha=0.7)

for bar in bars:
    yval = bar.get_height()
    plt.text(bar.get_x() + bar.get_width()/2, yval, f'{yval:.3f}s', 
             va='bottom', ha='center', fontsize=10)

plt.tight_layout()

plt.savefig('benchmark_comparison.png')
print("Graph saved as benchmark_comparison.png")
plt.show()