import matplotlib.pyplot as plt
import seaborn as sns
import pandas as pd

df = pd.read_excel('./data.xlsx')

df = df[df['payload'].isin([60, 500, 1500])]
df['latency'] = df['latency'] / 1_000_000

df['group'] = df['interval'].astype(str) + 'ms-' + df['payload'].astype(str)

network_types = ['ethernet', 'wifi', '4G static', '4G mobile']

for network in network_types:
    df_network = df[df['network'] == network]

    plt.figure(figsize=(12, 6))
    sns.boxplot(x='group', y='latency', hue='protocol', data=df_network, palette='muted')

    plt.title(f'{network.capitalize()} Network - Latency Distribution by Interval and Payload', fontsize=14)
    plt.xlabel('Interval and Payload Group', fontsize=12)
    plt.ylabel('Latency (ms)', fontsize=12)

    # plt.xticks(rotation=45, ha='right')

    plt.legend(title='Protocol')

    plt.savefig(f'{network}_latency_box_plot.png')

    plt.show()
