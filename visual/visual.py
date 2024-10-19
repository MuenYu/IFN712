import seaborn as sns
import matplotlib.pyplot as plt
import pandas as pd
from matplotlib.pyplot import ylabel

# Load the data
file = './summary.xlsx'
df = pd.read_excel(file)

# Prepare the data for plotting: combine 'payload (bytes)' and 'interval (ms)' into one column
df['payload_interval'] = df['payload (bytes)'].astype(str) + 'B_' + df['interval (ms)'].astype(str) + 'ms'

# List of statistics to visualize
stats = ['min', 'max', 'mean', 'middle', 'standard deviation']

for stat in stats:
    df[stat] = df[stat] / 1_000_000

# Create bar plots for each network
networks = df['network'].unique()

# Create the bar plots for each network again with separation by protocol
for network in networks:
    plt.figure(figsize=(12, 6))

    # Filter data for the specific network
    network_data = df[df['network'] == network]

    # Combine 'payload_interval' with 'protocol' to separate the bars by protocol
    network_data['payload_interval_protocol'] = network_data['payload_interval'] + '_' + network_data['protocol']

    # Reshape the data to a long format for seaborn
    network_data_long = pd.melt(network_data, id_vars=['payload_interval_protocol'], value_vars=stats,
                                var_name='Statistic', value_name='Value')

    # Create the barplot, separating by protocol
    sns.barplot(x='payload_interval_protocol', y='Value',hue='Statistic', data=network_data_long, errorbar=None)

    # Customize the plot
    plt.title(f'{network} - Min, Max, Mean, Middle, and Std Dev for Payload & Interval Combinations (TCP vs KCP)')
    plt.xticks(rotation=45, ha='right')
    plt.ylabel('Latency (ms)')
    plt.tight_layout()

    # Show the plot
    # plt.show()
    plt.savefig(f'{network}_latency.png')