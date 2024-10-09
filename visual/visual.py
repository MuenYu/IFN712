import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns

# Load the Excel file
file_path = './result.xlsx'  # Update the file path to your specific file
df = pd.read_excel(file_path)

# Set up the plotting aesthetics
sns.set(style="whitegrid")

# 1. Boxplot: Protocols (TCP vs KCP) vs Latency, with subplots for different network types
plt.figure(figsize=(12, 6))
sns.boxplot(x='protocol', y='latency', hue='network', data=df)
plt.title('Latency Distribution by Protocol and Network Type')
plt.ylabel('Latency (ns)')
plt.xlabel('Protocol')
plt.legend(title='Network Type')
plt.show()

# 2. Boxplot: Latency distribution for different payload sizes under each protocol
plt.figure(figsize=(12, 6))
sns.boxplot(x='protocol', y='latency', hue='payload', data=df)
plt.title('Latency Distribution by Protocol and Payload Size')
plt.ylabel('Latency (ns)')
plt.xlabel('Protocol')
plt.legend(title='Payload Size')
plt.show()

# 3. Line plot: Latency vs Interval for each protocol with different payloads
plt.figure(figsize=(12, 6))
sns.lineplot(x='interval', y='latency', hue='protocol', style='payload', data=df, markers=True, dashes=False)
plt.title('Latency vs Interval for TCP and KCP with Different Payload Sizes')
plt.ylabel('Latency (ns)')
plt.xlabel('Interval (ms)')
plt.legend(title='Protocol and Payload Size')
plt.show()

# 4. Bar plot: Average latency comparison across protocols and network types
plt.figure(figsize=(12, 6))
sns.barplot(x='protocol', y='latency', hue='network', data=df, errorbar='sd')
plt.title('Average Latency by Protocol and Network Type')
plt.ylabel('Latency (ns)')
plt.xlabel('Protocol')
plt.legend(title='Network Type')
plt.show()
