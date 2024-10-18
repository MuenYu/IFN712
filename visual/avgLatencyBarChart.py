import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns

df = pd.read_excel('./data.xlsx')
df = df[df['payload'].isin([60, 500, 1500])]

network_types = ['ethernet', 'wifi', '4G static', '4G mobile']
df['latency'] = df['latency'] / 1_000_000

for network in network_types:
    # 过滤数据，针对当前网络环境
    df_network = df[df['network'] == network]

    # 创建一个新的列，结合 interval 和 payload 作为分组的 x 轴标签
    df_network['group'] = df_network['interval'].astype(str) + 'ms-' + df_network['payload'].astype(str)

    # 按照 'group' 和 'protocol' 分组，计算每组的平均延迟
    avg_latency = df_network.groupby(['group', 'protocol'])['latency'].mean().reset_index()

    # 创建柱状图
    plt.figure(figsize=(12, 6))
    ax = sns.barplot(x='group', y='latency', hue='protocol', data=avg_latency, errorbar=None, palette='muted')

    # 在每个柱子顶部添加平均延迟值
    for p in ax.patches:
        ax.annotate(format(p.get_height(), '.2f'),
                    (p.get_x() + p.get_width() / 2., p.get_height()),
                    ha = 'center', va = 'center',
                    xytext = (0, 9),  # 设置标签在柱子上方9个点
                    textcoords = 'offset points')

    # 设置图表标题和标签
    plt.title(f'{network.capitalize()} Network - Average Latency by Interval and Payload', fontsize=14)
    plt.xlabel('Interval and Payload Group', fontsize=12)
    plt.ylabel('Average Latency (ms)', fontsize=12)

    # 旋转 x 轴标签以提高可读性
    plt.xticks(rotation=45, ha='right')

    # 添加图例
    plt.legend(title='Protocol')

    # 保存图像（你也可以选择不保存直接显示）
    plt.savefig(f'{network}_latency_bar_chart.png')

    # 显示图表
    # plt.show()