import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns

df = pd.read_excel('./data.xlsx')
df = df[df['payload'].isin([60, 500, 1500])]

network_types = ['ethernet', 'wifi', '4G static', '4G mobile']
df['latency'] = df['latency'] / 1_000_000

for network in network_types:
    df_network = df[df['network'] == network]

    df_network['group'] = df_network['interval'].astype(str) + 'ms-' + df_network['payload'].astype(str)

    avg_latency = df_network.groupby(['group', 'protocol'])['latency'].mean().reset_index()

    plt.figure(figsize=(12, 6))
    ax = sns.barplot(x='group', y='latency', hue='protocol', data=avg_latency, errorbar=None, palette='muted')

    for p in ax.patches:
        ax.annotate(format(p.get_height(), '.2f'),
                    (p.get_x() + p.get_width() / 2., p.get_height()),
                    ha = 'center', va = 'center',
                    xytext = (0, 9),
                    textcoords = 'offset points')

    plt.title(f'{network.capitalize()} Network - Average Latency by Interval and Payload', fontsize=14)
    plt.xlabel('Interval and Payload Group', fontsize=12)
    plt.ylabel('Average Latency (ms)', fontsize=12)

    # plt.xticks(rotation=45, ha='right')

    plt.legend(title='Protocol')

    plt.savefig(f'{network}_latency_bar_chart.png')

    plt.show()