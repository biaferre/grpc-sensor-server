import pandas as pd
import matplotlib.pyplot as plt

df_grpc = pd.read_csv("grpc_results.csv")
df_rabbit = pd.read_csv("rabbitmq_results.csv")

plt.plot(df_grpc["InputSize"], df_grpc["LatencyMicroseconds"], marker='o', label="gRPC")
plt.plot(df_rabbit["InputSize"], df_rabbit["LatencyMicroseconds"], marker='x', label="RabbitMQ")

plt.title("Latency vs Input Size")
plt.xlabel("Input Size")
plt.ylabel("Latency (μs)")
plt.grid(True)
plt.legend()

plt.show()
