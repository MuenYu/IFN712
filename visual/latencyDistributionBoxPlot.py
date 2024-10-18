import matplotlib.pyplot as plt
import seaborn as sns
import pandas as pd

# 读取数据
df = pd.read_excel('./data.xlsx')

# 过滤只考虑 payload 为 60, 500 和 1500 的数据
df = df[df['payload'].isin([60, 500, 1500])]
df['latency'] = df['latency'] / 1_000_000

# 创建一个新的列，将 interval 和 payload 组合为分组标识
df['group'] = df['interval'].astype(str) + 'ms-' + df['payload'].astype(str)

# 按网络环境生成 4 张图
network_types = ['ethernet', 'wifi', '4G static', '4G mobile']

# 绘制每个网络环境下的延迟分布情况
for network in network_types:
    # 过滤数据，针对当前网络环境
    df_network = df[df['network'] == network]

    # 创建箱线图
    plt.figure(figsize=(12, 6))
    sns.boxplot(x='group', y='latency', hue='protocol', data=df_network, palette='muted')

    # 设置图表标题和标签
    plt.title(f'{network.capitalize()} Network - Latency Distribution by Interval and Payload', fontsize=14)
    plt.xlabel('Interval and Payload Group', fontsize=12)
    plt.ylabel('Latency (ms)', fontsize=12)

    # 旋转 x 轴标签以提高可读性
    plt.xticks(rotation=45, ha='right')

    # 添加图例
    plt.legend(title='Protocol')

    # 保存图像
    plt.savefig(f'{network}_latency_box_plot.png')

    # 显示图表
    plt.show()
